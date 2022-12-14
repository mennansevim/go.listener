package kafka

import (
	"time"

	"github.com/Shopify/sarama"
)

func Config(version string) *sarama.Config {
	v, err := sarama.ParseKafkaVersion(version)
	if err != nil {
		panic(err)
	}

	config := sarama.NewConfig()
	config.Version = v

	config.Metadata.Retry.Max = 3
	config.Metadata.Retry.Backoff = 10 * time.Second

	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Group.Session.Timeout = 30 * time.Second
	config.Consumer.Group.Heartbeat.Interval = 6 * time.Second
	config.Consumer.MaxProcessingTime = 3 * time.Second
	config.Consumer.Fetch.Default = 2048 * 1024
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Group.Rebalance.Timeout = 3 * time.Minute

	config.Producer.Retry.Max = 3
	config.Producer.Retry.Backoff = 10 * time.Second
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Timeout = 3 * time.Minute
	config.Producer.MaxMessageBytes = 2000000

	config.Net.ReadTimeout = 3 * time.Minute
	config.Net.DialTimeout = 3 * time.Minute
	config.Net.WriteTimeout = 3 * time.Minute

	return config
}
