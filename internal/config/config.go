package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel     string        `env:"LOG_LEVEL"`
	Addr         string        `env:"ADDR"`
	WriteTimeout time.Duration `env:"WRITETIMEOUT"`
	ApiBaseURL   string        `env:"API_BASE_URL"`
	Database     ConfigDatabase
}

type ConfigDatabase struct {
	Port     string `env:"PORT" env-default:"5432"`
	Host     string `env:"HOST" env-default:"localhost"`
	Name     string `env:"NAME" env-default:"postgres"`
	User     string `env:"USER" env-default:"user"`
	Password string `env:"PASSWORD"`
	DBName   string `env:"DBNAME"`
}

func Load(path string) (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
