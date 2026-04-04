package config

import (
	"github.com/canonflow/canonflow-go-ddd/pkg/queue"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewQueue(viperConfig *viper.Viper, logrus *logrus.Logger) {
	//* Connect to RabbitMQ
	conn, err := amqp091.Dial(viperConfig.GetString("RABBITMQ_CONNECTION_URL"))
	if err != nil {
		logrus.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	//* Open a RabbitMQ Channel
	ch, err := conn.Channel()
	if err != nil {
		logrus.Fatalf("Failed to open a RabbitMQ Channel: %s", err)
	}

	//* Create DB
	db := NewDatabase(viperConfig, logrus)
	queueRepo := queue.NewQueueRepository(db)

	queue.QueueClient = &queue.Queue{
		Conn:       conn,
		Channel:    ch,
		Repository: queueRepo,
	}
}
