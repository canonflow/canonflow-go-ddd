package messaging

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/canonflow/canonflow-go-ddd/internal/contract"
	"github.com/canonflow/canonflow-go-ddd/internal/domain/user/model"
	"github.com/sirupsen/logrus"
)

type UserConsumer struct {
	Log *logrus.Logger
}

func NewUserConsumer(log *logrus.Logger) contract.ConsumerContract {
	return &UserConsumer{
		Log: log,
	}
}

func (c *UserConsumer) Consume(message *sarama.ConsumerMessage) error {
	userEvent := new(model.UserEvent)

	if err := json.Unmarshal(message.Value, userEvent); err != nil {
		c.Log.WithError(err).Error("[User Consumer] Error unmarshalling user event")
		return err
	}

	// TODO: process event
	c.Log.Infof("[User Consumer] Received topic users with event %v from partition %d", userEvent, message.Partition)
	return nil
}
