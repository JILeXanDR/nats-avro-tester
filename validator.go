package main

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Meta struct {
	Field  string
	Reason string
	Tag    string
	Path   string
	Value  interface{}
}

type ValidationErr struct {
	metas []Meta
}

func (ValidationErr) Error() string {
	return ""
}

type customValidator struct {
	validator *validator.Validate
	trans     ut.Translator
}

func (cv *customValidator) Validate(i interface{}) error {
	metas := cv.extractMetas(cv.validator.Struct(i))
	return metas
}

func (cv *customValidator) extractMetas(err error) error {
	var validationErr ValidationErr

	verrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}

	for _, err := range verrs {
		validationErr.metas = append(validationErr.metas, Meta{
			Field:  err.StructField(),
			Reason: err.Translate(cv.trans),
			Tag:    err.Tag(),
			Path:   err.Namespace(),
			Value:  err.Value(),
		})
	}

	return validationErr
}

func NewValidator() echo.Validator {
	trans, _ := ut.New(en.New()).GetTranslator("en")

	return &customValidator{
		validator: validator.New(),
		trans:     trans,
	}
}
