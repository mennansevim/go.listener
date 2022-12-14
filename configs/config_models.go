package configs

type KafkaConfig struct {
	Brokers string `required:"true" split_words:"true" json:"brokers"`
	Version string `required:"true" split_words:"true" json:"version"`
}
