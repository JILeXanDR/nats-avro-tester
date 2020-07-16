package api

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"nats-viewer/pkg/avro"
	"nats-viewer/pkg/errors"
	"nats-viewer/pkg/logger"
	"nats-viewer/pkg/nats"
	"nats-viewer/pkg/sse"

	"github.com/labstack/echo/v4"
)

type apiHandlers struct {
	nats        *nats.Client
	codecFinder avro.CodecStorage
	logger      *logger.Logger
	hub         sse.Hub
}

func NewAPIHandlers(nats *nats.Client, codecFinder avro.CodecStorage, logger *logger.Logger, hub sse.Hub) *apiHandlers {
	return &apiHandlers{
		nats:        nats,
		codecFinder: codecFinder,
		logger:      logger,
		hub:         hub,
	}
}

func (h *apiHandlers) GetSchemas(c echo.Context) error {
	codecs, err := h.codecFinder.GetAll()
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
		list = append(list, &responseItem{
			Name:      codec.Name(),
			Namespace: codec.Namespace(),
			Schema:    codec.Schema(),
			Example:   codec.Example(),
		})
	}

	return c.JSON(http.StatusOK, list)
}

func (h *apiHandlers) UploadSchema(c echo.Context) error {
	formFile, err := c.FormFile("file")
	if err != nil {
		return errors.Wrap(err, "getting form file")
	}

	f, err := formFile.Open()
	if err != nil {
		return errors.Wrap(err, "opening form file")
	}
	defer func() {
		if err := f.Close(); err != nil {
			h.logger.Err(err).Msg("failed to close form file descriptor")
		}
	}()

	schemas, err := h.parseSchemas(f)
	if err != nil {
		return err
	}

	if len(schemas) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "no schemas found")
	}

	if err := h.codecFinder.SyncSchemas(schemas...); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{})
}

// TODO: move to some service
func (h *apiHandlers) parseSchemas(f multipart.File) ([]string, error) {
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, errors.Wrap(err, "reading file content")
	}

	zipReader, err := zip.NewReader(f, int64(len(b)))
	if err != nil {
		return nil, errors.Wrap(err, "creating zip reader from form file content")
	}

	var schemas = make([]string, 0)
	for _, file := range zipReader.File {
		info := file.FileInfo()
		if info.IsDir() {
			continue
		}
		lgr := h.logger.With().Str("filename", info.Name()).Logger()
		f, err := file.Open()
		if err != nil {
			lgr.Err(err).Msg("failed to open file from ZIP archive")
			continue
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			lgr.Err(err).Msg("reading file content")
			continue
		}
		schemas = append(schemas, string(b))
	}
	return schemas, nil
}

func (h *apiHandlers) CreateMessage(c echo.Context) error {
	ctx := c.Request().Context()

	var req = struct {
		Subject string      `json:"subject" validate:"required"`
		Payload interface{} `json:"payload" validate:"required"`
	}{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	if err := h.nats.Publish(ctx, req.Subject, req.Payload); err != nil {
		return errors.Wrap(err, "publishing message")
	}

	return c.JSON(http.StatusOK, echo.Map{})
}

func (h *apiHandlers) MessagesStream(c echo.Context) error {
	client := sse.NewClient(uuid.New().String())
	h.hub.Register(client)
	go func() {
		<-c.Request().Context().Done()
		h.hub.Unregister(client)
	}()

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

	return client.Wait(func(msg interface{}) {
		b, err := json.Marshal(msg)
		if err != nil {
			h.logger.Err(err).Interface("data", msg).Msgf("marshaling data to JSON")
		} else {
			fmt.Fprintf(rw, "data: %s\n\n", string(b))
			flusher.Flush()
		}
	})
}
