package main

import (
	"github.com/karrick/goavro"
)

type CodecWrapper struct {
	*goavro.Codec
	name      string
	namespace string
	schema    map[string]interface{}
}

func (cw *CodecWrapper) Name() string {
	return cw.name
}

func (cw *CodecWrapper) Namespace() string {
	return cw.namespace
}

func (cw *CodecWrapper) Schema() map[string]interface{} {
	return cw.schema
}

type CodecStorage interface {
	GetAll() ([]*CodecWrapper, error)
	FindByNamespace(subject string) (*CodecWrapper, error)
	SyncSchemas(schemas ...string) error
}
