package nats

import (
	"context"
	"github.com/nats-io/nats.go"
	"nats-viewer/pkg/errors"
	"nats-viewer/pkg/logger"
	"strings"
)

type Client struct {
	conn   *nats.EncodedConn
	logger *logger.Logger
}

func NewClient(encoder nats.Encoder, server string, logger *logger.Logger) (*Client, error) {
	conn, err := nats.Connect(server)
	if err != nil {
		return nil, errors.WrapError(err, "connecting NATS server")
	}

	nats.RegisterEncoder("avro", encoder)

	ec, err := nats.NewEncodedConn(conn, "avro")
	if err != nil {
		return nil, errors.WrapError(err, "creating encoded connection")
	}

	return &Client{conn: ec, logger: logger}, nil
}

func (c *Client) Drain() error {
	return c.conn.Drain()
}

func (c *Client) Publish(ctx context.Context, subject string, message interface{}) error {
	c.logger.Debug().Str("subject", subject).Interface("data", message).Msg("publish message")
	err := c.conn.Publish(subject, message)
	if err != nil {
		return errors.WrapError(err, "publishing message to NATS")
	}
	return nil
}

func (c *Client) SubscribeAll(levels int64, next func(subject string, payload interface{})) error {
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
			l.Debug().Msg("subscription successful")
		}
	}
	return nil
}
