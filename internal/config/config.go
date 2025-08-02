package config

import (
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Port  string `env:"PORT" envDefault:"8080"`
	Host  string `env:"HOST" envDefault:"localhost"`
	Token int    `env:"TELEGRAM_TOKEN"`
}

type serverConfigBuilder struct {
	Config *Config
}

func (c *serverConfigBuilder) fromEnv() *serverConfigBuilder {
	if err := env.Parse(c); err != nil {
		log.Fatal(err)
	}
	return c
}

func GetConfig() (*Config, error) {
	var c serverConfigBuilder
	c.fromEnv()
	//TODO: add validation here
	return c.Config, nil
}
