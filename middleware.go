package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"strings"
)

func RequestLogger(logger *Logger) echo.MiddlewareFunc {
	return middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{
		Skipper: func(c echo.Context) bool {
			return !strings.HasPrefix(c.Request().URL.String(), "/api")
		},
		Handler: func(c echo.Context, reqBody, resBody []byte) {
			lgr := logger.With().
				Str("method", c.Request().Method).
				Str("path", c.Request().URL.String()).
				Logger()

			lgr.Debug().
				Str("request-id", c.Request().Header.Get(echo.HeaderXRequestID)).
				Bytes("body", reqBody).
				Msg("request")

			lgr.Debug().
				Str("request-id", c.Response().Header().Get(echo.HeaderXRequestID)).
				Bytes("body", resBody).
				Msg("response")
		},
	})
}
