package main

import (
	"fmt"
	kafka_wrapper "github.com/Trendyol/kafka-wrapper"
	"github.com/Trendyol/kafka-wrapper/execution_behaviour"
	"github.com/Trendyol/kafka-wrapper/execution_behaviour/behavioral"
	"github.com/go-resty/resty/v2"
	_ "gitlab.trendyol.com/tglobal/gocommons/logger"
	"gitlab.trendyol.com/tglobal/gocommons/newrelic"
	t "gitlab.trendyol.com/tglobal/gocommons/tracing"
	"gitlab.trendyol.com/tglobal/order-management/gp-product-consumer/configs"
	"gitlab.trendyol.com/tglobal/order-management/gp-product-consumer/pkg/common"
	"gitlab.trendyol.com/tglobal/order-management/gp-product-consumer/pkg/infra/errors"
	"gitlab.trendyol.com/tglobal/order-management/gp-product-consumer/pkg/infra/kafka"
	"go.uber.org/zap"
	"time"
)

const (
	ErrorPostFix  = "error"
	RetryPostFix  = "retry"
	ConsumerGroup = ""
)

func init() {
	configs.InitConfigs()
}

var quit = make(chan struct{})

func main() {
	zap.L().Info("server is starting")
	newrelic.CreateAgent(configs.AppConfig.NewRelic, true)

	zap.L().Info(fmt.Sprintf("[x] gp-product-consumer is running"))
	<-quit
}

func createKafkaConfig(topicName string) kafka_wrapper.ConnectionParameters {
	return kafka_wrapper.ConnectionParameters{
		ConsumerGroupID: ConsumerGroup,
		Brokers:         configs.AppConfig.ConfluentKafka.Brokers,
		Conf:            kafka.Config(configs.AppConfig.ConfluentKafka.Version),
		Topics:          []string{topicName},
		RetryTopic:      fmt.Sprintf("%s.%s.%s", topicName, ConsumerGroup, RetryPostFix),
		ErrorTopic:      fmt.Sprintf("%s.%s.%s", topicName, ConsumerGroup, ErrorPostFix),
	}
}

func newListener(config kafka_wrapper.ConnectionParameters, service behavioral.LogicOperator) error {
	producer, err := kafka_wrapper.NewProducer(config)
	if err != nil {
		return err
	}

	errorOperator := errors.NewErrorOperator()
	behaviourSelector := execution_behaviour.NewRetryBehaviourSelector(service, errorOperator, producer, 5,
		config.RetryTopic, config.ErrorTopic)
	eventHandler := kafka.NewKafkaListener(behaviourSelector)
	kafkaConsumer, err := kafka_wrapper.NewConsumer(config)
	if err != nil {
		return err
	}
	kafkaConsumer.Subscribe(eventHandler)

	return nil
}

func createRestyClient(timeoutValue int, retryCount int) *resty.Client {
	httpClient := resty.New().
		SetRetryCount(retryCount).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			return r.StatusCode() >= 500 || r.StatusCode() == 409 || r.StatusCode() == 0
		}).
		SetHeader("Content-Type", "application/json").
		SetHeader(t.AgentName, common.AgentName).
		SetTimeout(time.Duration(timeoutValue) * time.Millisecond)
	return httpClient
}
