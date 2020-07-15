package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"strconv"
	"time"
)

func main() {
	config, err := ReadConfigUsingEnv()
	if err != nil {
		panic(fmt.Sprintf("reading config: %+v", err))
	}

	logger, err := NewLogger(config.LogLevel)
	if err != nil {
		panic(fmt.Sprintf("creating logger: %+v", err))
	}

	logger.Debug().Interface("config", config).Msg("app initialization...")

	codecStorage, err := NewInMemoryCodecStorage()
	if err != nil {
		logger.Fatal().Err(err).Msg("creating local codec finder")
	}

	encoder, err := NewAvroEncoder(codecStorage, logger.As("AVRO-ENCODER"))
	if err != nil {
		logger.Fatal().Err(err).Msg("creating NATS avro encoder")
	}

	natsClient, err := NewNATSClient(encoder, config.NATSServer, logger.As("NATS-CLIENT"))
	if err != nil {
		logger.Fatal().Err(err).Msg("creating NATS client")
	}

	sseHub := NewSSEHub(logger.As("SSE-HUB"))
	go sseHub.Run()

	if err := natsClient.SubscribeAll(config.MaxHierarchyLevel, func(subject string, data interface{}) {
		go sseHub.Notify(map[string]interface{}{
			"subject": subject,
			"payload": data,
			"when":    time.Now().UTC().Format(time.RFC3339),
		})
	}); err != nil {
		logger.Error().Err(err).Msg("failed to subscribe to all NATS subjects")
	}

	defer func() {
		logger.Info().Msg("draining NATS connection...")
		if err := natsClient.Drain(); err != nil {
			logger.Err(err).Msg("failed to drain NATS connection")
		}
	}()

	handlers := newAPIHandlers(natsClient, codecStorage, logger, sseHub)

	e := echo.New()

	e.HideBanner, e.HidePort = true, true
	e.Validator = NewValidator()
	e.HTTPErrorHandler = HTTPErrorHandler(e)

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(middleware.BodyLimit("1M"))
	e.Use(middleware.Logger())
	//e.Use(RequestLogger(logger.As("http")))

	e.Static("/", "./web/dist")

	e.GET("/api/schemas", handlers.GetSchemas)
	e.POST("/api/schemas", handlers.UploadSchema)
	e.POST("/api/message", handlers.CreateMessage)
	e.GET("/api/stream", handlers.WriteMessagesStream)

	logger.Debug().Msgf("starting server at http://localhost:%d", config.Port)

	if err := e.Start(":" + strconv.Itoa(int(config.Port))); err != nil {
		logger.Fatal().Err(err).Msg("starting server")
	}
}
