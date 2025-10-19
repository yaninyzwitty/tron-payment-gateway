package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Debug          bool           `yaml:"debug"`
	AppPort        int            `yaml:"appPort"`
	DatabaseConfig DatabaseConfig `yaml:"database"`
}

type DatabaseConfig struct {
	User string `yaml:"user"`
	// Password       string `yaml:"password"` TODO-use environmental variable for password
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	Database       string `yaml:"database"`
	MaxConnections int    `yaml:"maxConnections"`
}

func (c *Config) LoadConfig(path string) error {
	f, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config %w", err)
	}

	if err := yaml.Unmarshal(f, c); err != nil {
		return fmt.Errorf("failed to parse config %w", err)
	}

	return nil
}
