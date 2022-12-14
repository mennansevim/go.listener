package errors

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/Trendyol/kafka-wrapper/execution_behaviour/behavioral"
)

type errorOperator struct {
}

func NewErrorOperator() behavioral.LogicOperator {
	return &errorOperator{}
}

func (p *errorOperator) Operate(ctx context.Context, message *sarama.ConsumerMessage) error {
	return nil
}
