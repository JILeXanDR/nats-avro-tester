package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"strconv"
)

func main() {
	config, err := ParseConfigUsingEnv()
	if err != nil {
		log.Fatalf("reading config: %+v", err)
	}
	log.Printf("loded config: %+v", config)

	codecStorage, err := NewInMemoryCodecStorage()
	if err != nil {
		log.Fatalf("creating local codec finder: %+v", err)
	}

	encoder, err := NewAvroEncoder(codecStorage)
	if err != nil {
		log.Fatalf("creating NATS avro encoder: %+v", err)
	}

	natsClient, err := NewNATSClient(encoder, config.NATSServer)
	if err != nil {
		log.Fatalf("creating NATS client: %+v", err)
	}

	messages := make(chan *PublishMessageRequest)

	closeSubscription, err := natsClient.SubscribeAll(func(subject string, data interface{}) {
		log.Printf("NATS message: %+v", data)
		//messages <- message
	})
	if err != nil {
		log.Printf("subscribing to all events: %+v", err)
	}

	defer func() {
		closeSubscription()
		natsClient.Drain()
	}()

	handlers := newAPIHandlers(natsClient, codecStorage)

	e := echo.New()

	e.Debug = true
	e.Validator = NewValidator()
	e.HTTPErrorHandler = HTTPErrorHandler(e)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.BodyLimit("1M"))

	e.Static("/", "./web/dist")

	e.GET("/api/schemas", handlers.GetSchemas)
	e.POST("/api/schemas", handlers.UploadSchema)
	e.POST("/api/message", handlers.CreateMessage)
	e.GET("/api/stream", handlers.ReadStream(messages))

	log.Printf("starting server at http://localhost:%d", config.Port)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(int(config.Port))))
}
