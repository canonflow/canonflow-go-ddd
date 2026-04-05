package main

// @title           Canonflow API
// @version         1.0
// @description     Canonflow API documentation
// @host            localhost:8000
// @BasePath        /

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"sync"
	"syscall"
	"time"

	_ "github.com/canonflow/canonflow-go-ddd/cmd/api/docs"
	"github.com/canonflow/canonflow-go-ddd/internal/config"
	"github.com/canonflow/canonflow-go-ddd/internal/contract"
	userQueue "github.com/canonflow/canonflow-go-ddd/internal/domain/user/queue"
	"github.com/canonflow/canonflow-go-ddd/pkg/queue"
	queuePkg "github.com/canonflow/canonflow-go-ddd/pkg/queue"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogrus(viperConfig)
	validate := config.NewValidator()

	// Set Mode for Gin
	if viperConfig.GetString("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
		log.Info("ON Production Environment\n")
	}

	app := config.NewGin(viperConfig, log)
	db := config.NewDatabase(viperConfig, log)
	redis := config.NewRedis(viperConfig, log)
	producer := config.NewKafkaProducer(viperConfig, log)

	// TODO: Init Queue
	config.NewQueue(viperConfig, log)
	queue.RegisterQueueHandler([]contract.QueueContract{
		userQueue.NewUserQueue(userQueue.QUEUE_NAME),
	})

	defer queuePkg.CloseConnection()

	// TODO: Get the web port
	webHost := viperConfig.GetString("DB_HOST")
	webPort := viperConfig.GetInt("WEB_PORT")

	url := ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", webPort))
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	// TODO: Bootstrap all configs
	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Redis:    redis,
		Log:      log,
		Validate: validate,
		Config:   viperConfig,
		Producer: producer,
	})

	server := &http.Server{
		Addr:    webHost + ":" + strconv.Itoa(webPort),
		Handler: app.Handler(),
	}

	go func() {
		// Run
		logrus.Infof("Running on: %s:%s", webHost, strconv.Itoa(webPort))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Listen: %s\n", err)
		}
	}()

	//* Run Queue
	ctxQueue, cancelQueue := context.WithCancel(context.Background())
	go runQueue(ctxQueue, log, db)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)

	// kill (no params) by default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Infoln("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		cancelQueue()
		logrus.Infoln("Server Shutdown:", err)
	}
	logrus.Infoln("Server exiting")
}

func runQueue(ctx context.Context, logrus *logrus.Logger, db *gorm.DB) {
	logrus.Infof("Registered handlers: %d", len(queuePkg.QueueClient.Handler))
	var wg sync.WaitGroup

	for queueName, handler := range queuePkg.QueueClient.Handler {
		logrus.Infof("Consuming queue: %s", queueName)
		wg.Add(1)
		go func(ctx context.Context, queueName string, handler contract.QueueContract, db *gorm.DB) {
			defer wg.Done()

			logrus.Infof("Calling ConsumeQueue for: %s", queueName)
			msgs, err := queuePkg.ConsumeQueue(ctx, queueName)
			if err != nil {
				logrus.Errorf("Failed to consume queue %s: %s", queueName, err)
			}

			logrus.Infof("Waiting for messages on: %s", queueName)

			go func() {
				<-ctx.Done()
				logrus.Errorf("!!! ctx was cancelled: %v", ctx.Err())
				debug.PrintStack()
			}()

			for msg := range msgs {
				var queueMessage queuePkg.QueueMessage

				err := json.Unmarshal(msg.Body, &queueMessage)
				if err != nil {
					logrus.Printf("Error reading coffee order (please check the JSON format): %s", err)
					msg.Nack(false, false)
					continue
				}

				//* Handle the message
				if err := handler.Handle(queueMessage.Payload); err != nil {
					logrus.Printf("Failed to handle message for queue %s: %s", queueName, err)
					msg.Nack(false, false)
					queuePkg.UpdateQueueStatus(queueMessage.UniqueID, queuePkg.StatusFailed)
					continue
				}

				//* Acknowledge the message
				msg.Ack(false)

				//* Update the queue record status to success
				queuePkg.UpdateQueueStatus(queueMessage.UniqueID, queuePkg.StatusDone)
			}
			logrus.Infof("Consumer loop exited for: %s", queueName)
		}(ctx, queueName, handler, db)
	}

	wg.Wait()
	logrus.Errorf("!!! runQueue exited — all consumers dead")
}
