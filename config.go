package main

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	Port       uint
	NATSServer string
}

func ReadConfigUsingFlags() (*Config, error) {
	cfg := Config{
		Port:       *flag.Uint("port", 8999, "Server port"),
		NATSServer: *flag.String("nats_server", "http://localhost:4222", "NATS server"),
	}
	flag.Parse()
	return &cfg, nil
}

func ParseConfigUsingEnv() (*Config, error) {
	cfg := Config{
		Port:       8080,
		NATSServer: "http://localhost:4222",
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return nil, WrapError(err, "getting PORT env var")
	} else {
		cfg.Port = uint(port)
	}

	natsServer := os.Getenv("NATS_SERVER")
	if natsServer == "" {
		return nil, NewError("env var NATS_SERVER is empty")
	} else {
		cfg.NATSServer = natsServer
	}

	return &cfg, nil
}
