package queue

import (
	"fmt"
)

var QUEUE_NAME = "user_queue"

type UserQueue struct {
	QueueName string
	// Payload    map[string]interface{}
}

func NewUserQueue(name string) *UserQueue {
	return &UserQueue{
		QueueName: name,
	}
}

func (q *UserQueue) Name() string {
	return q.QueueName
}

func (q *UserQueue) Handle(payload map[string]interface{}) error {
	fmt.Printf(">>> [User Queue] Handling queue: %s\n", q.QueueName)

	for key, value := range payload {
		fmt.Printf(">>> [User Queue] Key: %s, Value: %v\n", key, value)
	}

	return nil
}
