package configs

import "github.com/caarlos0/env/v11"

type Config struct {
	DdHost     string `env:"DB_HOST"`
	DbPort     string `env:"DB_PORT"`
	DbName     string `env:"DB_NAME"`
	DbUsername string `env:"DB_USERNAME"`
	DbPassword string `env:"DB_PASSWORD"`
}

func NewConfig() (*Config, error) {
	cfg, err := env.ParseAs[Config]()

	return &cfg, err
}
