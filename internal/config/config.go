package config

import (
	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"
)

type Config struct {
	SECRET string `env:"SECRET" envDefault:"yourSecretKey"`
}

func NewConfig() (*Config, error) {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	logrus.Println(cfg)

	return &cfg, nil
}
