package config

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/HeyReyHR/keyvaluedb/internal/common"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Network NetworkConfig
	Logging LoggingConfig
}

type NetworkConfig struct {
	Address        string
	MaxConnections int
	MaxMessageSize int
	IdleTimeout    time.Duration
}
type LoggingConfig struct {
	Level  slog.Level
	Output string
}
type rawConfig struct {
	Network struct {
		Address        string `yaml:"address"`
		MaxConnections int    `yaml:"max_connections"`
		MaxMessageSize string `yaml:"max_message_size"`
		IdleTimeout    string `yaml:"idle_timeout"`
	} `yaml:"network"`

	Logging struct {
		Level string `yaml:"level"`
	} `yaml:"logging"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %w", err)
	}

	var raw rawConfig
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("unable to parse yaml: %w", err)
	}

	return normalize(raw)
}

func normalize(raw rawConfig) (*Config, error) {
	var cfg Config

	cfg.Network.Address = raw.Network.Address
	cfg.Network.MaxConnections = raw.Network.MaxConnections

	size, err := common.ParseSize(raw.Network.MaxMessageSize)
	if err != nil {
		return nil, fmt.Errorf("network.max_message_size: %w", err)
	}
	cfg.Network.MaxMessageSize = size

	timeout, err := time.ParseDuration(raw.Network.IdleTimeout)
	if err != nil {
		return nil, fmt.Errorf("network.idle_timeout: %w", err)
	}
	cfg.Network.IdleTimeout = timeout

	var level slog.Level
	if err := level.UnmarshalText([]byte(raw.Logging.Level)); err != nil {
		return nil, fmt.Errorf("logging.level: %w", err)
	}
	cfg.Logging.Level = level

	return &cfg, nil
}
