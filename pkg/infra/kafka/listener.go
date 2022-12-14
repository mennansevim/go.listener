package kafka

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
	kafka_wrapper "github.com/Trendyol/kafka-wrapper"
	"github.com/Trendyol/kafka-wrapper/execution_behaviour"
	"go.uber.org/zap"
)

type eventHandler struct {
	behavioralSelector execution_behaviour.BehavioralSelector
}

func NewKafkaListener(selector execution_behaviour.BehavioralSelector) kafka_wrapper.EventHandler {
	return &eventHandler{
		behavioralSelector: selector,
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (e *eventHandler) Setup(session sarama.ConsumerGroupSession) error {
	zap.L().Info("kafka listener is starting")
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (e *eventHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (e *eventHandler) ConsumeClaim(session sarama.ConsumerGroupSession, c sarama.ConsumerGroupClaim) error {
	processor := e.behavioralSelector.GetBehavioral(c)
	for {
		select {
		case message := <-c.Messages():
			ctx := context.Background()
			err := processor.Process(ctx, message)
			if err != nil {
				zap.L().Error(fmt.Sprintf("Error executing message: %+v , err: %+v", message.Value, err))
			}
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
