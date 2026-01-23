package messaging

import (
	"github.com/IBM/sarama"
	"github.com/canonflow/canonflow-go-ddd/internal/contract"
	"github.com/canonflow/canonflow-go-ddd/internal/domain/user/model"
	"github.com/canonflow/canonflow-go-ddd/pkg/broker"
	"github.com/sirupsen/logrus"
)

type UserProducer struct {
	ProducerHandler broker.ProducerHandler[*model.UserEvent]
}

func NewUserProducer(producer sarama.SyncProducer, log *logrus.Logger) contract.ProducerContract {
	return &UserProducer{
		ProducerHandler: broker.ProducerHandler[*model.UserEvent]{
			Producer: producer,
			Topic:    "users",
			Log:      log,
		},
	}
}

func (u *UserProducer) GetTopic() *string {
	return u.ProducerHandler.GetTopic()
}

func (u *UserProducer) Send(event contract.Event) error {
	// userEvent := event.(*model.UserEvent)
	return u.ProducerHandler.Send(event)
}
