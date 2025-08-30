package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Port            string   `env:"PORT" envDefault:"8080"`
	Host            string   `env:"HOST" envDefault:""`
	Token           string   `env:"TELEGRAM_TOKEN"`
	PokingInterval  int      `env:"POKING_INTERVAL" envDefault:"1000"`
	PokingURL       string   `env:"POKING_URL" envDefault:"https://ya.ru"`
	DefaultLanguage Language `env:"DEFAULT_LANGUAGE" envDefault:"rus"`
}

func GetConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil { //nolint:typecheck
		return nil, err
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) validate() error {
	if _, ok := SupportedLanguages[c.DefaultLanguage]; !ok {
		return fmt.Errorf("language '%s' is not supported", c.DefaultLanguage)
	}
	return nil
}
