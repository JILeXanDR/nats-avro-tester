package logger

import (
	"github.com/rs/zerolog"
	"nats-viewer/pkg/errors"
	"os"
	"time"
)

type Logger struct {
	*zerolog.Logger
}

func New(level string) (*Logger, error) {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger := zerolog.New(output).With().Timestamp().Logger()
	if l, err := zerolog.ParseLevel(level); err != nil {
		return nil, errors.Wrap(err, "parsing log level %s", level)
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
