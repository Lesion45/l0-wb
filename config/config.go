package config

import (
	"fmt"
	"github.com/caarlos0/env/v8"
)

type Config struct {
	Env   string `env:"ENV,required"`
	PgDSN string `env:"POSTGRES_DSN,required"`
	Kafka Kafka
}

type Kafka struct {
	Host  string `env:"CONFIG_KAFKA_HOST,required"`
	Topic string `env:"CONFIG_KAFKA_TOPIC,required"`
}

// MustLoad loads configuration from config.yaml
// Throw a panic if the config doesn't exist or if there is an error reading the config.
func MustLoad() *Config {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		panic(fmt.Sprintf("Failed to load config: %s", err))
	}

	return &cfg
}
