package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		Log  `yaml:"logger"`
	}

	App struct {
		Name    string `env:"APP_NAME"    env-required:"true" yaml:"name"`
		Version string `env:"APP_VERSION" env-required:"true" yaml:"version"`
	}

	HTTP struct {
		Port string `env:"HTTP_PORT" env-required:"true" yaml:"port"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL" env-required:"true" yaml:"log_level"`
	}
)

func NewConfig() (*Config, error) {
	//nolint:exhaustruct // Fields are initialized by cleanenv.ReadConfig
	cfg := &Config{}

	configPath := os.Getenv("CONFIG_PATH")

	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	return cfg, nil
}
