package broker

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/canonflow/canonflow-go-ddd/internal/contract"
	"github.com/sirupsen/logrus"
)

type ConsumerGroupHandler struct {
	Handler contract.ConsumerContract
	Log     *logrus.Logger
}

func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			err := h.Handler.Consume(message)
			if err != nil {
				h.Log.WithError(err).Error("Failed to process message")
			} else {
				session.MarkMessage(message, "")
			}

		case <-session.Context().Done():
			return nil
		}
	}
}

func ConsumeTopic(ctx context.Context, consumerGroup sarama.ConsumerGroup, topic string, log *logrus.Logger, handler contract.ConsumerContract) {
	consumerHandler := &ConsumerGroupHandler{
		Handler: handler,
		Log:     log,
	}

	go func() {
		for {
			if err := consumerGroup.Consume(ctx, []string{topic}, consumerHandler); err != nil {
				log.WithError(err).Error("Error from consumer")
			}

			if ctx.Err() != nil {
				log.Info("Context cancelled, stopping consumer")
				return
			}
		}
	}()

	go func() {
		for err := range consumerGroup.Errors() {
			log.WithError(err).Error("Error from consumer")
		}
	}()
	<-ctx.Done()

	log.Infof("Closing consumer group for topic: %s", topic)
	if err := consumerGroup.Close(); err != nil {
		log.WithError(err).Error("Error closing consumer group")
	}
}
