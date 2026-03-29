package main

// @title           Canonflow API
// @version         1.0
// @description     Canonflow API documentation
// @host            localhost:8000
// @BasePath        /

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "github.com/canonflow/canonflow-go-ddd/cmd/api/docs"
	"github.com/canonflow/canonflow-go-ddd/internal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
		logrus.Infoln("Server Shutdown:", err)
	}
	logrus.Infoln("Server exiting")
}
