package broker

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/canonflow/canonflow-go-ddd/internal/contract"
	"github.com/sirupsen/logrus"
)

type ProducerHandler[T contract.Event] struct {
	Producer sarama.SyncProducer
	Topic    string
	Log      *logrus.Logger
}

func (p *ProducerHandler[T]) GetTopic() *string {
	return &p.Topic
}

func (p *ProducerHandler[T]) Send(event contract.Event) error {
	value, err := json.Marshal(event)
	if err != nil {
		p.Log.WithError(err).Error("Failed to marshal event")
		return err
	}

	message := &sarama.ProducerMessage{
		Topic: p.Topic,
		Key:   sarama.StringEncoder(event.GetId()),
		Value: sarama.ByteEncoder(value),
	}

	partition, offset, err := p.Producer.SendMessage(message)
	if err != nil {
		p.Log.WithError(err).Error("failed to produce message")
		return err
	}

	p.Log.Debugf("Message sent to topic %s, partition %d, offset %d", p.Topic, partition, offset)
	return nil
}
