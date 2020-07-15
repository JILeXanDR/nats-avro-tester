package main

import (
	"github.com/rs/zerolog"
	"os"
	"time"
)

type Logger struct {
	*zerolog.Logger
}

func NewLogger(level string) (*Logger, error) {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger := zerolog.New(output).With().Timestamp().Logger()
	if l, err := zerolog.ParseLevel(level); err != nil {
		return nil, WrapError(err, "parsing log level %s", level)
	} else {
		logger = logger.Level(l)
	}
	return &Logger{Logger: &logger}, nil
}

func (l *Logger) As(name string) *Logger {
	logger := l.With().Str("as", name).Logger()
	return &Logger{
		Logger: &logger,
	}
}
