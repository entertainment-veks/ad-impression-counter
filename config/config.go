package config

import (
	"log"
	"time"

	"github.com/caarlos0/env"
)

type Config struct {
	ServerPort string        `env:"SERVER_PORT" envDefault:"8080"` // ServerPort is the port the HTTP server will listen on.
	TTL        time.Duration `env:"TTL" envDefault:"3600s"`        // TTL is the duplicate impressions detection window duration in seconds.
}

func LoadConfig() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Printf("failed to load config: %v", err)
		return Config{}, err
	}

	return cfg, nil
}
