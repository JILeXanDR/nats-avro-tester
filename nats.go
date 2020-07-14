package main

import (
	"context"
	"github.com/nats-io/nats.go"
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

func (c *natsClient) SubscribeAll(next func(string, interface{})) error {
	_, err := c.conn.QueueSubscribe("*", "*", func(subject string, vPrt interface{}) {
		c.logger.Debug().Str("subject", subject).Interface("data", vPrt).Msg("got decoded from subscriber")
		next(subject, vPrt)
	})
	if err != nil {
		return WrapError(err, "subscribing to subject %s", "*")
	}
	return nil
}
