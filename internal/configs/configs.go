package configs

import (
	"github.com/caarlos0/env/v6"
	"log"
)

type Config struct {
	RunAddr        string `env:"RUN_ADDRESS" envDefault:"127.0.0.1:8080"`
	DatabaseURI    string `env:"DATABASE_URI" envDefault:"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"`
	AccrualSysAddr string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

//envDefault:"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

func NewConfigs() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal(err)
	}
	log.Printf("configs: %v", *cfg)
	return cfg
}
