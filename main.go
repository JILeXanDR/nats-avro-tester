package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

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

	defer func() {
		lgr.Info().Msg("draining NATS connection...")
		if err := natsClient.Drain(); err != nil {
			lgr.Err(err).Msg("failed to drain NATS connection")
		}
	}()

	if err := natsClient.SubscribeAll(cfg.MaxHierarchyLevel, func(subject string, data interface{}) {
		sseHub.NotifyAll(map[string]interface{}{
			"subject": subject,
			"payload": data,
			"when":    time.Now().UTC().Format(time.RFC3339),
		})
	}); err != nil {
		lgr.Error().Err(err).Msg("failed to subscribe to all NATS subjects")
	}

	handlers := api.NewAPIHandlers(natsClient, codecStorage, lgr, sseHub)

	e := echo.New()

	e.HideBanner, e.HidePort, e.Debug = true, true, true
	e.Validator = validator.New()

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		switch val := err.(type) {
		case *echo.HTTPError:
			c.JSON(val.Code, echo.Map{
				"message": val.Message,
			})
		default:
			c.JSON(http.StatusInternalServerError, echo.Map{
				"message": err.Error(),
			})
		}
	}

	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.BodyLimit("1M"))

	e.Static("/", "./web/dist")

	e.GET("/api/schemas", handlers.GetSchemas)
	e.POST("/api/schemas", handlers.UploadSchema)
	e.POST("/api/message", handlers.CreateMessage)
	e.GET("/api/stream", handlers.MessagesStream)
	e.GET("/api/check_version", handlers.CheckVersion)

	lgr.Info().Msgf("starting server at http://localhost:%d", cfg.Port)

	if err := e.Start(":" + strconv.Itoa(int(cfg.Port))); err != nil {
		lgr.Fatal().Err(err).Msg("starting server")
	}
}
