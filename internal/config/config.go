package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Port           string `env:"PORT" envDefault:"8080"`
	Host           string `env:"HOST" envDefault:""`
	Token          string `env:"TELEGRAM_TOKEN"`
	PokingInterval int    `env:"POKING_INTERVAL" envDefault:"1000"`
	PokingURL      string `env:"POKING_URL" envDefault:"https://ya.ru"`
}

func GetConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil { //nolint:typecheck
		return nil, err
	}
	// TODO: add validation here
	return cfg, nil
}
