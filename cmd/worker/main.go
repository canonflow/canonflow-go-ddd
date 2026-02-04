package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/canonflow/canonflow-go-ddd/internal/config"
	"github.com/canonflow/canonflow-go-ddd/internal/domain/user/delivery/messaging"
	"github.com/canonflow/canonflow-go-ddd/pkg/broker"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	viper := config.NewViper()
	logger := config.NewLogrus(viper)
	logger.Info("Starting worker service")

	ctx, cancel := context.WithCancel(context.Background())

	// Run Consumer
	go RunUserConsumer(logger, viper, ctx)

	// Graceful shutdown
	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	<-terminateSignals
	logger.Info("Got one of stop signals, shutting down worker gracefully")
	cancel()

	time.Sleep(5 * time.Second)
}

func RunUserConsumer(logger *logrus.Logger, viperConfig *viper.Viper, ctx context.Context) {
	logger.Info("Running user consumer")
	userConsumerGroup := config.NewKafkaConsumerGroup(viperConfig, logger)
	userHandler := messaging.NewUserConsumer(logger)

	broker.ConsumeTopic(ctx, userConsumerGroup, broker.USER_TOPIC, logger, userHandler)
}
