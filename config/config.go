package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	Port              uint   `envconfig:"port" default:"8080"`
	NATSServer        string `envconfig:"nats_server" default:"http://localhost:4222"`
	MaxHierarchyLevel int64  `envconfig:"max_hierarchy_level" default:"1"`
	LogLevel          string `envconfig:"log_level" default:"trace"`
}

func NewFromEnv(prefix string) (*Config, error) {
	var cfg Config
	if err := envconfig.Process(prefix, &cfg); err != nil {
		return nil, errors.Wrap(err, "building config based on environment variables")
	}
	return &cfg, nil
}
