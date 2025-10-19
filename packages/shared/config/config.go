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
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	Host           string `yaml:"host"`
	Database       string `yaml:"database"`
	MaxConnections int    `yaml:"maxConnections"`
}

func (c *Config) LoadConfig(path string) error {
	f, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config %v", err)
	}

	if err := yaml.Unmarshal(f, c); err != nil {
		return fmt.Errorf("failed to parse config %v", err)
	}

	return nil
}
