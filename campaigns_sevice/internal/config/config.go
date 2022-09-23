package config

import "github.com/caarlos0/env/v6"

type Config struct {
	Port      string `env:"PORT" envDefault:"8080"`
	DBConn    string `env:"DBConn" envDefault:"host=localhost user=postgres password=postgres dbname=backend port=5432 sslmode=disable"`
	RedisAddr string `env:"DB_CONN_STRING,notEmpty,unset" envDefault:"localhost:6379"`
}

func New() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
