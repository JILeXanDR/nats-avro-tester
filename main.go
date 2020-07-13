package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"strconv"
)

func main() {
	logger, err := NewLogger()
	if err != nil {
		panic(fmt.Sprintf("creating logger: %+v", err))
	}

	config, err := ReadConfigUsingEnv()
	if err != nil {
		logger.Fatal().Err(err).Msg("reading config")
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

	messages := make(chan *PublishMessageRequest)

	closeSubscription, err := natsClient.SubscribeAll(func(subject string, data interface{}) {
		logger.Debug().Str("subject", subject).Interface("data", data).Msg("receive message")
		//messages <- message
	})
	if err != nil {
		logger.Error().Err(err).Msg("subscribing to all events")
	}

	defer func() {
		closeSubscription()
		natsClient.Drain()
	}()

	handlers := newAPIHandlers(natsClient, codecStorage, logger)

	e := echo.New()

	e.Validator = NewValidator()
	e.HTTPErrorHandler = HTTPErrorHandler(e)

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(middleware.BodyLimit("1M"))
	//e.Use(RequestLogger(logger.As("http")))

	e.Static("/", "./web/dist")

	e.GET("/api/schemas", handlers.GetSchemas)
	e.POST("/api/schemas", handlers.UploadSchema)
	e.POST("/api/message", handlers.CreateMessage)
	e.GET("/api/stream", handlers.ReadStream(messages))

	logger.Debug().Msgf("starting server at http://localhost:%d", config.Port)

	if err := e.Start(":" + strconv.Itoa(int(config.Port))); err != nil {
		logger.Fatal().Err(err).Msg("starting server")
	}
}
