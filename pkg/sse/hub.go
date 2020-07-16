package sse

import "nats-viewer/pkg/logger"

type hub struct {
	clients    map[string]*client
	connect    chan *client
	disconnect chan *client

	logger *logger.Logger
}

func NewHub(logger *logger.Logger) Hub {
	return &hub{
		clients:    make(map[string]*client),
		connect:    make(chan *client),
		disconnect: make(chan *client),
		logger:     logger,
	}
}

func (hub *hub) Register(client *client) {
	hub.logger.Debug().Str("id", client.ID()).Msg("register client")
	hub.connect <- client
}

func (hub *hub) Unregister(client *client) {
	hub.logger.Debug().Str("id", client.ID()).Msg("unregister client")
	hub.disconnect <- client
}

func (hub *hub) NotifyAll(v interface{}) {
	hub.logger.Debug().Int("count", len(hub.clients)).Msg("notify all clients")
	for _, client := range hub.clients {
		client.data <- v
	}
}

func (hub *hub) Run() {
	for {
		select {
		case client := <-hub.connect:
			hub.clients[client.ID()] = client
			hub.logger.Debug().Str("id", client.ID()).Int("count", len(hub.clients)).Msg("client registered")
		case client := <-hub.disconnect:
			delete(hub.clients, client.ID())
			close(client.data)
			client.closed <- struct{}{}
			hub.logger.Debug().Str("id", client.ID()).Int("count", len(hub.clients)).Msg("client unregistered")
		default:
			continue
		}
	}
}
