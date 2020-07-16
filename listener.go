package main

import (
	"nats-viewer/pkg/logger"
	"nats-viewer/pkg/nats"
	"nats-viewer/pkg/sse"
	"time"
)

func ListenNATSMessages(natsClient *nats.Client, sseHub sse.Hub, lgr *logger.Logger, maxHierarchyLevel int64) {
	if err := natsClient.SubscribeAll(maxHierarchyLevel, func(subject string, data interface{}) {
		sseHub.NotifyAll(map[string]interface{}{
			"subject": subject,
			"payload": data,
			"when":    time.Now().UTC().Format(time.RFC3339),
		})
	}); err != nil {
		lgr.Error().Err(err).Msg("failed to subscribe to all NATS subjects")
	}
}
