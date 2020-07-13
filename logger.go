package main

import (
	"github.com/rs/zerolog"
	"os"
	"time"
)

type Logger struct {
	*zerolog.Logger
}

func NewLogger() (*Logger, error) {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger := zerolog.New(output).With().Timestamp().Logger()
	return &Logger{Logger: &logger}, nil
}

func (l *Logger) As(name string) *Logger {
	logger := l.With().Str("as", name).Logger()
	return &Logger{
		Logger: &logger,
	}
}
