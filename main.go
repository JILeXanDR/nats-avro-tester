package main

import (
	"fmt"
	"strconv"

	"nats-viewer/api"
	"nats-viewer/config"
	"nats-viewer/pkg/avro"
	"nats-viewer/pkg/logger"
	"nats-viewer/pkg/nats"
	"nats-viewer/pkg/sse"
	"nats-viewer/pkg/validator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg, err := config.NewFromEnv("")
	if err != nil {
		panic(fmt.Sprintf("reading config: %+v", err))
	}

	lgr, err := logger.New(cfg.LogLevel)
	if err != nil {
		panic(fmt.Sprintf("creating logger: %+v", err))
	}

	lgr.Info().Interface("config", cfg).Msg("app initialization...")

	exampleGenerator := avro.NewExampleGenerator()

	codecStorage, err := avro.NewInMemoryCodecStorage(exampleGenerator)
	if err != nil {
		lgr.Fatal().Err(err).Msg("creating local codec finder")
	}

	encoder, err := avro.NewAvroEncoder(codecStorage, lgr.As("AVRO-ENCODER"))
	if err != nil {
		lgr.Fatal().Err(err).Msg("creating NATS avro encoder")
	}

	natsClient, err := nats.NewClient(encoder, cfg.NATSServer, lgr.As("NATS-CLIENT"))
	if err != nil {
		lgr.Fatal().Err(err).Msg("creating NATS client")
	}

	sseHub := sse.NewHub(lgr.As("SSE-HUB"))
	go sseHub.Run()

	// TODO: use struct?
	ListenNATSMessages(natsClient, sseHub, lgr, cfg.MaxHierarchyLevel)

	defer func() {
		lgr.Info().Msg("draining NATS connection...")
		if err := natsClient.Drain(); err != nil {
			lgr.Err(err).Msg("failed to drain NATS connection")
		}
	}()

	handlers := api.NewAPIHandlers(natsClient, codecStorage, lgr, sseHub)

	e := echo.New()

	e.HideBanner, e.HidePort = true, true
	e.Validator = validator.New()
	//e.HTTPErrorHandler = HTTPErrorHandler(e)

	e.Use(middleware.Recover())
	//e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(middleware.BodyLimit("1M"))
	//e.Use(middleware.Logger())
	//e.Use(RequestLogger(logger.As("http")))

	e.Static("/", "./web/dist")

	e.GET("/api/schemas", handlers.GetSchemas)
	e.POST("/api/schemas", handlers.UploadSchema)
	e.POST("/api/message", handlers.CreateMessage)
	e.GET("/api/stream", handlers.MessagesStream)

	lgr.Info().Msgf("starting server at http://localhost:%d", cfg.Port)

	if err := e.Start(":" + strconv.Itoa(int(cfg.Port))); err != nil {
		lgr.Fatal().Err(err).Msg("starting server")
	}
}
