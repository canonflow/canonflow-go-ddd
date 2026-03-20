package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/canonflow/canonflow-go-ddd/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

	// TODO: Get the web port
	webHost := viperConfig.GetString("DB_HOST")
	webPort := viperConfig.GetInt("WEB_PORT")

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
