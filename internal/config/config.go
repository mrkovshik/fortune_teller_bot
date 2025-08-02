package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Port  string `env:"PORT" envDefault:"8080"`
	Host  string `env:"HOST" envDefault:"localhost"`
	Token string `env:"TELEGRAM_TOKEN"`
}

func GetConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	// TODO: add validation here
	return cfg, nil
}
