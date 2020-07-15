package main

import (
	"context"
	"github.com/nats-io/nats.go"
	"strings"
)

type natsClient struct {
	conn   *nats.EncodedConn
	logger *Logger
}

func NewNATSClient(encoder nats.Encoder, server string, logger *Logger) (*natsClient, error) {
	conn, err := nats.Connect(server)
	if err != nil {
		return nil, WrapError(err, "connecting NATS server")
	}

	nats.RegisterEncoder("avro", encoder)

	ec, err := nats.NewEncodedConn(conn, "avro")
	if err != nil {
		return nil, WrapError(err, "creating encoded connection")
	}

	return &natsClient{conn: ec, logger: logger}, nil
}

func (c *natsClient) Drain() error {
	return c.conn.Drain()
}

func (c *natsClient) Publish(ctx context.Context, subject string, message interface{}) error {
	c.logger.Debug().Str("subject", subject).Interface("data", message).Msg(`publish data`)
	err := c.conn.Publish(subject, message)
	if err != nil {
		return WrapError(err, "publishing message to NATS connection")
	}
	return nil
}

func (c *natsClient) SubscribeAll(levels int64, next func(string, interface{})) error {
	max := make([]string, 0, levels)
	for i := 0; i < int(levels); i++ {
		max = append(max, "*")
		wildcardSubject := strings.Join(max, ".")
		l := c.logger.With().Str("subject", wildcardSubject).Logger()
		_, err := c.conn.Subscribe(wildcardSubject, func(subject string, vPrt interface{}) {
			c.logger.Debug().Str("subject", subject).Interface("data", vPrt).Msg("got decoded from subscriber")
			next(subject, vPrt)
		})
		if err != nil {
			l.Err(err).Msg("failed to subscribe")
		} else {
			l.Info().Msg("subscription successful")
		}
	}
	return nil
}
