package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Server     Server
	PostgreSQL PostgreSQL
}
type Server struct {
	Port int `env:"SERVER_PORT"`
}

type PostgreSQL struct {
	Host     string `env:"POSTGRES_HOST"`
	Conns    int    `env:"POSTGRES_CONNS"`
	Port     int    `env:"POSTGRES_PORT"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Db       string `env:"POSTGRES_DB"`
}

func MustLoadConfig() (Config, error) {
	var cfg Config
	err := godotenv.Load()
	if err != nil {
		err = fmt.Errorf("Error loading .env file: %w", err)
		return cfg, err
	}

	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		err = fmt.Errorf("Error reading .env file: %w", err)
		return cfg, err
	}

	return cfg, nil
}
