package configs

import (
	"github.com/spf13/viper"
	c "[REPLACEME]/gocommons/configs"
	"os"
	"strconv"
)

type appConfig struct {
	v              *viper.Viper
	ConfluentKafka KafkaConfig      `required:"true" split_words:"true" yaml:"confluentKafka"`
	NewRelic       c.NewRelicConfig `required:"false" split_words:"true" yaml:"newRelic"`
	OrderApi       c.ClientConfig   `required:"true" split_words:"true" yaml:"orderApi"`
}

func (a *appConfig) readWithViper(shouldPanic bool) error {
	if a.v == nil {
		v := viper.New()
		v.AddConfigPath("./config")
		v.SetConfigName("configs")
		v.SetConfigType("yaml")
		a.v = v
	}

	// manual overrides from environment
	a.NewRelic.LicenseKey = os.Getenv("NEW_RELIC_LICENSE_KEY")
	a.NewRelic.AppName = os.Getenv("NEW_RELIC_APP_NAME")
	a.NewRelic.AgentEnabled, _ = strconv.ParseBool(os.Getenv("NEW_RELIC_AGENT_ENABLED"))

	err := a.v.ReadInConfig()
	if err != nil {
		if shouldPanic {
			panic(err)
		}
		return err
	}

	err = a.v.Unmarshal(&AppConfig)
	if err != nil {
		if shouldPanic {
			panic(err)
		}
		return err
	}

	return nil
}
