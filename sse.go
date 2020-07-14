package main

import "time"

type SSEHub interface {
	Run()
	Register(*sseClient)
	Unregister(*sseClient)
	Notify(interface{})
}

type sseHub struct {
	clients    map[int64]*sseClient
	connect    chan *sseClient
	disconnect chan *sseClient

	logger *Logger
}

func NewSSEHub(logger *Logger) SSEHub {
	return &sseHub{
		clients:    make(map[int64]*sseClient),
		connect:    make(chan *sseClient),
		disconnect: make(chan *sseClient),
		logger:     logger,
	}
}

type sseClient struct {
	ID     int64
	data   chan interface{}
	closed chan struct{}
}

func (c *sseClient) Wait(next func(interface{})) error {
	for {
		select {
		case <-c.closed:
			return nil
		case v, ok := <-c.data:
			if ok {
				next(v)
			}
		}
	}
}

func NewSSEClient() *sseClient {
	return &sseClient{
		ID:     time.Now().UnixNano(),
		data:   make(chan interface{}),
		closed: make(chan struct{}, 1),
	}
}

func (hub *sseHub) Register(client *sseClient) {
	hub.logger.Info().Int64("id", client.ID).Msg("register client")
	go func() {
		hub.connect <- client
	}()
}

func (hub *sseHub) Unregister(client *sseClient) {
	hub.logger.Info().Int64("id", client.ID).Msg("unregister client")
	go func() {
		hub.disconnect <- client
	}()
}

func (hub *sseHub) Notify(v interface{}) {
	hub.logger.Info().Int("count", len(hub.clients)).Msg("notify all clients")
	for _, client := range hub.clients {
		client.data <- v
	}
}

func (hub *sseHub) Run() {
	for {
		select {
		case client := <-hub.connect:
			hub.clients[client.ID] = client
			hub.logger.Info().Int64("id", client.ID).Int("count", len(hub.clients)).Msg("registered")
		case client := <-hub.disconnect:
			delete(hub.clients, client.ID)
			close(client.data)
			client.closed <- struct{}{}
			hub.logger.Info().Int64("id", client.ID).Int("count", len(hub.clients)).Msg("unregistered")
		default:
			continue
		}
	}
}
