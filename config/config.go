package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"path/filepath"
)

const DefaultConfigPath = "config/config.yaml"

type (
	Config struct {
		Postgresql `yaml:"postgresql"`
	}

	Postgresql struct {
		Address  string `yaml:"address"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	}
)

func NewConfig(configPath *string) (*Config, error) {
	if configPath == nil {
		configPath = new(string)
		*configPath = DefaultConfigPath
	}
	cfg := Config{}

	err := cleanenv.ReadConfig(*configPath, &cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	err = cleanenv.UpdateEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("error updating env: %w", err)
	}

	return &cfg, nil
}

func NewConfigWithDiscover(configPath *string) (*Config, error) {
	var currentPath string
	if configPath == nil {
		configPath = new(string)
		currentPath = DefaultConfigPath
	} else {
		currentPath = *configPath
	}
	for tries := 10; tries > 0; tries-- {
		if _, err := os.Stat(currentPath); err == nil {
			config, err := NewConfig(&currentPath)
			return config, err
		}
		currentPath = filepath.Join("..", currentPath)
	}
	return nil, fmt.Errorf("could not discover config")
}
