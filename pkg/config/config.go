package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"project/pkg/logger"
)

type Config[H, S any] struct {
	App struct {
		Env    string
		Logger logger.Config
	}
	Handler H
	Service S
}

func Load[H, S any]() *Config[H, S] {
	var cfg Config[H, S]
	b, err := os.ReadFile("conf.yaml")
	if err != nil {
		log.Fatal(err)
	}
	if err = yaml.Unmarshal(b, &cfg); err != nil {
		log.Fatal(err)
	}
	return &cfg
}
