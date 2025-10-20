package config

import "fmt"

type Config struct {
	Address  string `yaml:"address"`
	Port     string `yaml:"port"`
	LogLevel string `yaml:"log_level"`
}

func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Address, c.Port)
}

func (c *Config) GetLogLevel() string {
	return c.LogLevel
}
