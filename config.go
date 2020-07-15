package main

import (
	"os"
	"strconv"
)

type Config struct {
	Port              uint
	NATSServer        string
	MaxHierarchyLevel int64
	LogLevel          string
}

func ReadConfigUsingEnv() (*Config, error) {
	cfg := Config{
		Port:              8080,
		NATSServer:        "http://localhost:4222",
		MaxHierarchyLevel: 1,
		LogLevel:          "trace",
	}

	port := os.Getenv("PORT")
	if port != "" {
		val, err := strconv.Atoi(port)
		if err != nil {
			return nil, WrapError(err, "getting PORT env var")
		} else {
			cfg.Port = uint(val)
		}
	}

	natsServer := os.Getenv("NATS_SERVER")
	if natsServer != "" {
		cfg.NATSServer = natsServer
	}

	maxHierarchyLevel := os.Getenv("MAX_HIERARCHY_LEVEL")
	if maxHierarchyLevel != "" {
		val, err := strconv.Atoi(maxHierarchyLevel)
		if err != nil {
			return nil, WrapError(err, "getting MAX_HIERARCHY_LEVEL env var")
		} else if val < 1 {
			return nil, NewError("MAX_HIERARCHY_LEVEL must be >= 1")
		} else {
			cfg.MaxHierarchyLevel = int64(val)
		}
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel != "" {
		cfg.LogLevel = logLevel
	}

	return &cfg, nil
}
