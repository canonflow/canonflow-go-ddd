package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/canonflow/canonflow-go-ddd/internal/config"
	"github.com/canonflow/canonflow-go-ddd/internal/contract"
	queuePkg "github.com/canonflow/canonflow-go-ddd/pkg/queue"
	"gorm.io/gorm"
)

func main() {
	viperConfig := config.NewViper()
	logrus := config.NewLogrus(viperConfig)
	db := config.NewDatabase(viperConfig, logrus)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// TODO: Consume all messages in RabbitMQ Queue
	for queueName, handler := range queuePkg.QueueClient.Handler {
		logrus.Infof("Consuming queue: %s", queueName)
		go func(ctx context.Context, queueName string, handler contract.QueueContract, db *gorm.DB) {
			msgs, err := queuePkg.ConsumeQueue(ctx, queueName)
			if err != nil {
				logrus.Errorf("Failed to consume queue %s: %s", queueName, err)
			}

			for msg := range msgs {
				var queueMessage queuePkg.QueueMessage

				err := json.Unmarshal(msg.Body, &queueMessage)
				if err != nil {
					logrus.Printf("Error reading coffee order (please check the JSON format): %s", err)
					continue
				}

				//* Inject Payload to Handler
				if err := handler.InjectPayload(queueMessage.Payload); err != nil {
					logrus.Printf("Failed to inject payload to handler for queue %s: %s", queueName, err)
					queuePkg.UpdateQueueStatus(queueMessage.UniqueID, queuePkg.StatusFailed)
					continue
				}

				//* Handle the message
				if err := handler.Handle(); err != nil {
					logrus.Printf("Failed to handle message for queue %s: %s", queueName, err)
					queuePkg.UpdateQueueStatus(queueMessage.UniqueID, queuePkg.StatusFailed)
					continue
				}

				//* Acknowledge the message
				msg.Ack(false)

				//* Update the queue record status to success
				queuePkg.UpdateQueueStatus(queueMessage.UniqueID, queuePkg.StatusDone)
			}
		}(ctx, queueName, handler, db)

	}

	// Graceful shutdown
	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	<-terminateSignals
	logrus.Info("Got one of stop signals, shutting down queue gracefully")
	cancel()

	time.Sleep(5 * time.Second)
}
