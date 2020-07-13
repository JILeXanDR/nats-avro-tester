package main

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ErrorDetails interface{}

type ErrorResponse struct {
	Message string         `json:"message"`
	Code    string         `json:"code"`
	Details []ErrorDetails `json:"details"`
}

func (er *ErrorResponse) Error() string {
	return ""
}

func WrapError(err error, format string, a ...interface{}) error {
	if err == nil {
		return NewError(format, a...)
	}
	text := fmt.Sprintf(format, a...)
	return fmt.Errorf(text+": %w", err)
}

func NewError(format string, a ...interface{}) error {
	return errors.New(fmt.Sprintf(format, a...))
}

func ErrCause(err error) error {
	for true {
		parent := errors.Unwrap(err)
		if parent == nil {
			break
		}
		err = parent
	}
	return err
}

func HTTPErrorHandler(e *echo.Echo) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		//causedErr := ErrCause(err)
		causedErr := err

		if causedErr == nil {
			e.DefaultHTTPErrorHandler(err, c)
			return
		}

		switch causedErr.(type) {
		case *echo.HTTPError:
			e.DefaultHTTPErrorHandler(err, c)
			return
		}

		switch v := causedErr.(type) {
		case ValidationErr:
			var details []ErrorDetails
			for _, meta := range v.metas {
				details = append(details, meta)
			}
			c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
				Message: "Failed to validate request data",
				Code:    "validation_error",
				Details: details,
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: v.Error(),
				Code:    "unexpected_error",
			})
			return
		}
	}
}
