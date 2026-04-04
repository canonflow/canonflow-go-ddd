package queue

import (
	"fmt"

	"github.com/canonflow/canonflow-go-ddd/pkg/queue"
)

var QUEUE_NAME = "user_queue"

type UserQueue struct {
	QueueName  string
	Payload    map[string]interface{}
	Repository *queue.QueueRepository
}

func NewUserQueue(name string) *UserQueue {
	return &UserQueue{
		QueueName: name,
	}
}

func (q *UserQueue) Name() string {
	return q.QueueName
}

func (q *UserQueue) InjectPayload(payload map[string]interface{}) error {
	q.Payload = payload
	return nil
}

func (q *UserQueue) Handle() error {
	fmt.Printf(">>> [User Queue] Handling queue: %s\n", q.QueueName)

	for key, value := range q.Payload {
		fmt.Printf(">>> [User Queue] Key: %s, Value: %v\n", key, value)
	}

	return nil
}
