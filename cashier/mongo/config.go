package db

import (
	"fmt"

	env "github.com/Netflix/go-env"
	"github.com/joho/godotenv"
)

type config struct {
	DBHost     string `env:"DB_HOST"`
	DBPort     string `env:"DB_PORT"`
	DBName     string `env:"DB_NAME"`
	DBUser     string `env:"DB_USER"`
	DBPwd      string `env:"DB_PWD"`
	BindAddr   string `env:"BIND_ADDR"`
	ConnString string
}

type Config interface {
	BindAddress() string
	DatabaseConnString() string
	DatabaseName() string
}

func NewConfig(path string) (Config, error) {
	var cfg config

	if err := godotenv.Load(path); err != nil {
		return nil, fmt.Errorf("failed to load .env: %v", err)
	}

	es, err := env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse env: %v", err)
	}

	if err := env.Unmarshal(es, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config from env: %v", err)
	}

	cfg.ConnString = fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", cfg.DBUser, cfg.DBPwd, cfg.DBHost, cfg.DBPort, cfg.DBName)

	return cfg, nil
}

func (c config) BindAddress() string {
	return c.BindAddr
}

func (c config) DatabaseConnString() string {
	return c.ConnString
}

func (c config) DatabaseName() string {
	return c.DBName
}
