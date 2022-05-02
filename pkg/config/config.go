package config

import (
	"github.com/goccy/go-yaml"
	"log"
	"os"
	"project/pkg/logger"
)

type Config[H, S any] struct {
	Env     string
	Logger  logger.Config
	Handler H
	Service S
}

func Load[H, S any](path string) *Config[H, S] {
	var cfg Config[H, S]
	b, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	if err = yaml.Unmarshal(b, &cfg); err != nil {
		log.Fatal(err)
	}
	return &cfg
}
