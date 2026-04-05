package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/canonflow/canonflow-go-ddd/internal/contract"
	"github.com/rabbitmq/amqp091-go"
)

var (
	QueueClient *Queue
	QueueTTL    = 24 * time.Hour
)

type Queue struct {
	ConsumerConn  *amqp091.Connection
	PublisherConn *amqp091.Connection
	Repository    *QueueRepository
	Handler       map[string]contract.QueueContract
}

func Dispatch(ctx context.Context, handler contract.QueueContract, message QueueMessage) error {
	// QueueClient.Channel.QueueDeclare("d", true, false, false, false, amqp091.Table{
	// 	"x-expires": QueueTTL.Milliseconds(),
	// })

	if QueueClient == nil {
		return errors.New("queue client is not initialized")
	}

	if QueueClient.PublisherConn == nil {
		return errors.New("publisher connection is nil")
	}
	if QueueClient.PublisherConn.IsClosed() {
		return errors.New("publisher connection is closed")
	}

	//* Insert to DB
	tx := QueueClient.Repository.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	//* Create Queue Record
	payloadString, err := json.Marshal(message.Payload)
	if err != nil {
		return err
	}

	queueRecord := &QueueRecord{
		UniqueID: message.UniqueID,
		Queue:    strings.ToLower(handler.Name()),
		Payload:  string(payloadString),
		Status:   StatusPending,
	}

	//* Insert Queue Record to DB
	if err := QueueClient.Repository.Create(tx, queueRecord); err != nil {
		return err
	}

	//*
	queueString, err := json.Marshal(message)
	if err != nil {
		return err
	}

	//* RabbitMQ Channel
	ch, err := QueueClient.PublisherConn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()

	//* Declare queue in RabbitMQ
	q, err := ch.QueueDeclare(
		strings.ToLower(handler.Name()),
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
		return err
	}

	//* Publish to RabbitMQ
	err = ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        queueString,
		},
	)
	if err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func ConsumeQueue(ctx context.Context, queue string) (<-chan amqp091.Delivery, error) {
	//* RabbitMQ Channel
	ch, err := QueueClient.ConsumerConn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// TODO: Declare the queue to ensure it exists
	q, err := ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
		return nil, err
	}

	// SUBSCRIBE TO THE QUEUE
	msgs, err := ch.Consume(
		q.Name, // queue
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Failed to register a RabbitMQ consumer: %s", err)
		return nil, err
	}

	// go func() {
	// 	chClosing := ch.NotifyClose(make(chan *amqp091.Error, 1))
	// 	connClosing := QueueClient.ConsumerConn.NotifyClose(make(chan *amqp091.Error, 1))
	// 	select {
	// 	case <-ctx.Done():
	// 		logrus.Errorf("ctx cancelled: %v", ctx.Err())
	// 		ch.Close()
	// 	case err := <-chClosing:
	// 		logrus.Errorf("channel closed: %v", err)
	// 	case err := <-connClosing:
	// 		logrus.Errorf("CONNECTION closed: %v", err) // ← I bet this fires
	// 	}
	// }()
	go func() {
		<-ctx.Done()
		ch.Close()
	}()

	return msgs, nil
}

func UpdateQueueStatus(uniqueID string, status Status) error {
	if QueueClient == nil {
		return errors.New("queue client is not initialized")
	}

	return QueueClient.Repository.UpdateStatus(uniqueID, status)
}

func RegisterQueueHandler(queues []contract.QueueContract) {
	for _, handler := range queues {
		queueName := strings.ToLower(handler.Name())
		//* Check handler exist in handler map, if not exist then add to handler map
		if _, ok := QueueClient.Handler[queueName]; !ok {
			QueueClient.Handler[queueName] = handler
		}
	}
}

func CloseConnection() {
	if QueueClient != nil {
		if err := QueueClient.PublisherConn.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ connection: %s", err)
		}
		if err := QueueClient.ConsumerConn.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ connection: %s", err)
		}
	}
}
