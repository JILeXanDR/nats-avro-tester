package main

import (
	"context"
	"github.com/nats-io/nats.go"
)

type natsClient struct {
	conn *nats.EncodedConn
}

func (c *natsClient) Drain() error {
	return c.conn.Drain()
}

func NewNATSClient(encoder nats.Encoder, server string) (*natsClient, error) {
	conn, err := nats.Connect(server)
	if err != nil {
		return nil, WrapError(err, "connecting NATS server")
	}

	nats.RegisterEncoder("avro", encoder)

	ec, err := nats.NewEncodedConn(conn, "avro")
	if err != nil {
		return nil, WrapError(err, "creating encoded connection")
	}

	return &natsClient{conn: ec}, nil
}

func (c *natsClient) Publish(ctx context.Context, subject string, message interface{}) error {
	err := c.conn.Publish(subject, message)
	if err != nil {
		return WrapError(err, "publishing message to NATS connection")
	}
	return nil
}

func (c *natsClient) SubscribeAll(next func(string, interface{})) (func(), error) {
	subscription, err := c.conn.Subscribe("*", func(data map[string]interface{}) {
		println("xxx")
		//next(m.Subject, m.Data)
	})
	if err != nil {
		return nil, WrapError(err, "queue subsribing")
	}
	closeFn := func() {
		subscription.Unsubscribe()
	}
	return closeFn, nil
}
