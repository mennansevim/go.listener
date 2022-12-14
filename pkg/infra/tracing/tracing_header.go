package tracing

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	t "gitlab.trendyol.com/tglobal/gocommons/tracing"
)

const (
	CorrelationId = "x-correlationid"
	AgentName     = "x-agentname"
	ExecutorUser  = "x-executor-user"
)

func SetHeaderToContext(ctx context.Context, headers []*sarama.RecordHeader) context.Context {
	for _, v := range headers {
		if string(v.Key) == CorrelationId {
			ctx = context.WithValue(ctx, CorrelationId, string(v.Value))
		}

		if string(v.Key) == AgentName {
			ctx = context.WithValue(ctx, AgentName, string(v.Value))
		}

		if string(v.Key) == ExecutorUser {
			ctx = context.WithValue(ctx, ExecutorUser, string(v.Value))
		}
	}

	if ctx.Value(t.CorrelationId) == nil {
		ctx = context.WithValue(ctx, t.CorrelationId, uuid.New().String())
	}

	return ctx
}
