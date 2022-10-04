package configs

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
	"os"
)

type Config struct {
	RunAddr        string `env:"RUN_ADDRESS" envDefault:"127.0.0.1:8080"`
	DatabaseURI    string `env:"DATABASE_URI" envDefault:"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"`
	AccrualSysAddr string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	SecretKey      string `env:"SECRET_KEY"`
}

//envDefault:"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

func NewConfigs() *Config {
	cfg := &Config{}
	err := os.Setenv("SECRET_KEY", "be55d1079e6c6167118ac91318fe")
	if err != nil {
		log.Fatal(err)
	}
	if err = env.Parse(cfg); err != nil {
		log.Fatal(err)
	}
	log.Printf("configs: %#v", *cfg)
	return cfg
}

func (c *Config) ParseFlags() {
	flag.StringVar(&c.RunAddr, "a", c.RunAddr, "Run server address")
	flag.StringVar(&c.DatabaseURI, "d", c.DatabaseURI, "Database URI")
	flag.StringVar(&c.AccrualSysAddr, "r", c.AccrualSysAddr, "Accrual system address")
	flag.Parse()
}
