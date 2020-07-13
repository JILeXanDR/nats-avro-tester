package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"log"
	"net/http"
)

type apiHandlers struct {
	nats        *natsClient
	codecFinder CodecStorage
}

func newAPIHandlers(nats *natsClient, codecFinder CodecStorage) *apiHandlers {
	return &apiHandlers{
		nats:        nats,
		codecFinder: codecFinder,
	}
}

func (a *apiHandlers) GetSchemas(c echo.Context) error {
	codecs, err := a.codecFinder.GetAll()
	if err != nil {
		return err
	}

	type responseItem struct {
		Name      string                 `json:"name"`
		Namespace string                 `json:"namespace"`
		Schema    interface{}            `json:"schema"`
		Example   map[string]interface{} `json:"example"`
	}

	list := make([]*responseItem, 0, len(codecs))

	for _, codec := range codecs {
		schema := codec.Schema()
		list = append(list, &responseItem{
			Name:      codec.Name(),
			Namespace: codec.Namespace(),
			Schema:    schema,
			Example:   GenerateAvroJSONExample(schema),
		})
	}

	return c.JSON(http.StatusOK, list)
}

func (a *apiHandlers) UploadSchema(c echo.Context) error {
	formFile, err := c.FormFile("file")
	if err != nil {
		return err
	}

	f, err := formFile.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return WrapError(err, "creating zip reader")
	}

	var schemas = make([]string, 0)
	for _, file := range zipReader.File {
		if file.FileInfo().IsDir() {
			continue
		}
		f, err := file.Open()
		if err != nil {
			continue
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			err = WrapError(err, "reading file %s", file.Name)
			log.Printf("err: %+v", err)
			continue
		}
		schemas = append(schemas, string(b))
	}

	if len(schemas) == 0 {
		return NewError("file is empty")
	}

	if err := a.codecFinder.SyncSchemas(schemas...); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": fmt.Sprintf("total count of schemas were uploaded %d", len(schemas)),
	})
}

func (a *apiHandlers) CreateMessage(c echo.Context) error {
	ctx := c.Request().Context()
	var req PublishMessageRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}
	if err := a.nats.Publish(ctx, req.Subject, req.Payload); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "ok"})
}

func (a *apiHandlers) ReadStream(messages chan *PublishMessageRequest) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
		c.Response().Header().Set("Content-Type", "text/event-stream")
		c.Response().Header().Set("Cache-Control", "no-cache")
		c.Response().Header().Set("Connection", "keep-alive")
		c.Response().WriteHeader(http.StatusOK)

		rw := c.Response().Writer

		flusher, ok := rw.(http.Flusher)
		if !ok {
			return echo.NewHTTPError(http.StatusBadRequest, "Streaming unsupported!")
		}

		for {
			msg := <-messages
			//msg.ID = time.Now().UTC().String()
			b, err := json.Marshal(msg)
			if err != nil {
				log.Printf("marshaling msg struct into JSON: %+v", err)
				continue
			}
			fmt.Fprintf(rw, "data: %s\n\n", string(b))
			flusher.Flush()
		}
	}
}
