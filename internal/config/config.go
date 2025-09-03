package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
)

type (
	Config struct {
		Port            string   `env:"PORT" envDefault:"8080"`
		Host            string   `env:"HOST" envDefault:""`
		Token           string   `env:"TELEGRAM_TOKEN"`
		PokingInterval  int      `env:"POKING_INTERVAL" envDefault:"10000"`
		PokingURL       string   `env:"POKING_URL" envDefault:"https://ya.ru"`
		DefaultLanguage Language `env:"DEFAULT_LANGUAGE" envDefault:"rus"`
		DatabaseURI     string   `env:"DATABASE_URI"`
		RDB             *RDB
	}
	RDB struct {
		Addr         string        `env:"RDB_ADDRESS"`
		Password     string        `env:"RDB_PASSWORD"`
		DB           int           `env:"RDB_DB"`
		DialTimeout  time.Duration `env:"RDB_DIAL_TIMEOUT" envDefault:"5s"`
		ReadTimeout  time.Duration `env:"RDB_READ_TIMEOUT" envDefault:"5s"`
		WriteTimeout time.Duration `env:"RDB_WRITE_TIMEOUT" envDefault:"5s"`
	}
)

func GetConfig() (*Config, error) {
	var rdb RDB
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil { //nolint:typecheck
		return nil, err
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	if err := env.Parse(&rdb); err != nil {
		return nil, err
	}
	if rdb.Addr != "" {
		cfg.RDB = &rdb
	}
	return cfg, nil
}

func (c *Config) validate() error {
	if _, ok := SupportedLanguages[c.DefaultLanguage]; !ok {
		return fmt.Errorf("language '%s' is not supported", c.DefaultLanguage)
	}
	return nil
}
