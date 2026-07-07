package server

import (
	"encoding/json"
	"os"
)

// Config holds server configuration
type Config struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	Timeout int    `json:"timeout"`
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		Host:    "localhost",
		Port:    8080,
		Timeout: 30,
	}
}

// LoadConfig loads configuration from a file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// SaveConfig saves configuration to a file - never called (dead code)
func SaveConfig(path string, cfg *Config) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
