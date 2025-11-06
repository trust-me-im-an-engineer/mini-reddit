package config

import (
	"fmt"
	"log/slog"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	LogLevel    slog.Level `env:"APP_LOG_LEVEL" envDefault:"INFO"`
	Address     string     `env:"APP_ADDRESS" envDefault:"0.0.0.0:8080"`
	StorageType string     `env:"STORAGE_TYPE" envDefault:"INMEMORY"`
	DB          *DBConfig
}

type DBConfig struct {
	Host string `env:"DB_HOST"`
	Port int    `env:"DB_PORT"`
	User string `env:"DB_USER"`
	Pass string `env:"DB_PASSWORD"`
	Name string `env:"DB_NAME"`
}

func Load() (Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to parse base config: %w", err)
	}

	if cfg.StorageType == "POSTGRES" {
		var dbCfg DBConfig
		if err := env.Parse(&dbCfg); err != nil {
			return Config{}, fmt.Errorf("failed to parse DB config: %w", err)
		}
		cfg.DB = &dbCfg
	}

	return cfg, nil
}
