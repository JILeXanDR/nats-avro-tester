package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/karrick/goavro"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type localCodecStorage struct {
	codecs map[string]*CodecWrapper
}

func NewLocalCodecStorage(dir string) (CodecStorage, error) {
	lcf := &localCodecStorage{}

	schemas, err := lcf.readAllAvroSchemas(dir)
	if err != nil {
		return nil, fmt.Errorf("reading all avro schemas: %w", err)
	}

	codecs := make(map[string]*CodecWrapper)

	for _, content := range schemas {
		var schema map[string]interface{}
		if err := json.NewDecoder(bytes.NewBuffer(content)).Decode(&schema); err != nil {
			err = WrapError(err, "decoding from JSON into a map")
			return nil, err
		}

		name, ok := schema["name"].(string)
		if !ok || name == "" {
			return nil, NewError("schema name not found or empty")
		}
		namespace, ok := schema["namespace"].(string)
		if !ok || namespace == "" {
			return nil, NewError("schema namespace not found or empty")
		}
		codec, err := goavro.NewCodec(string(content))
		if err != nil {
			return nil, fmt.Errorf("creating new codec from avro schema: %w", err)
		}
		codecs[namespace] = &CodecWrapper{
			Codec:     codec,
			schema:    schema,
			name:      name,
			namespace: namespace,
		}
	}

	lcf.codecs = codecs

	return lcf, nil
}

// GetAll returns all codes were found in the local directory.
func (l *localCodecStorage) GetAll() ([]*CodecWrapper, error) {
	res := make([]*CodecWrapper, 0, len(l.codecs))
	for _, codec := range l.codecs {
		res = append(res, codec)
	}
	return res, nil
}

// FindByName space returns codec found with specified "namespace" from the schema definition.
func (l *localCodecStorage) FindByNamespace(namespace string) (*CodecWrapper, error) {
	for k, codec := range l.codecs {
		if k == namespace {
			return codec, nil
		}
	}
	return nil, nil
}

func (l *localCodecStorage) SyncSchemas(schemas ...string) error {
	return nil
}

func (l *localCodecStorage) readAllAvroSchemas(dir string) ([][]byte, error) {
	list := make([][]byte, 0)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, "avsc") {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return WrapError(err, "opening path %s", path)
		}

		content, err := ioutil.ReadAll(f)
		if err != nil {
			return WrapError(err, "reading content of file")
		}

		list = append(list, content)
		return nil
	})
	if err != nil {
		return nil, WrapError(err, "walking inside directory %s", dir)
	}
	return list, nil
}
