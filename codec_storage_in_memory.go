package main

import (
	"bytes"
	"encoding/json"
	"github.com/karrick/goavro"
)

type imMemoryCodecStorage struct {
	// key is namespace
	codecs           map[string]*CodecWrapper
	exampleGenerator *ExampleGenerator
}

func NewInMemoryCodecStorage(exampleGenerator *ExampleGenerator) (CodecStorage, error) {
	lcf := &imMemoryCodecStorage{
		codecs:           make(map[string]*CodecWrapper),
		exampleGenerator: exampleGenerator,
	}
	return lcf, nil
}

func (l *imMemoryCodecStorage) GetAll() ([]*CodecWrapper, error) {
	res := make([]*CodecWrapper, 0, len(l.codecs))
	for _, codec := range l.codecs {
		res = append(res, codec)
	}
	return res, nil
}

func (l *imMemoryCodecStorage) FindByNamespace(namespace string) (*CodecWrapper, error) {
	for k, codec := range l.codecs {
		if k == namespace {
			return codec, nil
		}
	}
	return nil, nil
}

func (l *imMemoryCodecStorage) SyncSchemas(schemas ...string) error {
	l.codecs = make(map[string]*CodecWrapper)

	for _, schema := range schemas {
		if err := l.addSchema(schema); err != nil {
			return err
		}
	}

	return nil
}

func (l *imMemoryCodecStorage) addSchema(content string) error {
	var schema map[string]interface{}

	if err := json.NewDecoder(bytes.NewBuffer([]byte(content))).Decode(&schema); err != nil {
		return WrapError(err, "decoding text avro schema into a map")
	}

	name, _ := schema["name"].(string)
	namespace, _ := schema["namespace"].(string)

	codec, err := goavro.NewCodec(content)
	if err != nil {
		return WrapError(err, "creating new codec from avro schema")
	}

	example, err := l.exampleGenerator.Generate(schema)
	if err != nil {
		return WrapError(err, "generating schema example, name: %s, namespace: %s", name, namespace)
	}

	l.codecs[namespace] = &CodecWrapper{
		Codec:     codec,
		name:      name,
		namespace: namespace,
		schema:    schema,
		example:   example,
	}

	return nil
}
