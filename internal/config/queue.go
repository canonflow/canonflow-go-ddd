package config

import (
	"github.com/canonflow/canonflow-go-ddd/internal/contract"
	"github.com/canonflow/canonflow-go-ddd/pkg/queue"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewQueue(viperConfig *viper.Viper, logrus *logrus.Logger) {
	//* Connect to RabbitMQ
	publisherConn, err := amqp091.Dial(viperConfig.GetString("RABBITMQ_CONNECTION_URL"))
	if err != nil {
		logrus.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	consumerConn, err := amqp091.Dial(viperConfig.GetString("RABBITMQ_CONNECTION_URL"))
	if err != nil {
		logrus.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	//* Create DB
	db := NewDatabase(viperConfig, logrus)
	queueRepo := queue.NewQueueRepository(db)

	queue.QueueClient = &queue.Queue{
		ConsumerConn:  consumerConn,
		PublisherConn: publisherConn,
		Repository:    queueRepo,
		Handler:       map[string]contract.QueueContract{},
	}
}
