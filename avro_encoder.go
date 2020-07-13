package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
)

type avroEncoder struct {
	finder CodecStorage
}

func NewAvroEncoder(finder CodecStorage) (nats.Encoder, error) {
	return &avroEncoder{finder: finder}, nil
}

func (enc *avroEncoder) Encode(subject string, v interface{}) ([]byte, error) {
	codec, err := enc.finder.FindByNamespace(subject)
	if err != nil {
		return nil, WrapError(err, "finding codec by subject %s", subject)
	}
	if codec == nil {
		return nil, NewError("codec for subject %s not found", subject)
	}

	native, err := codec.BinaryFromNative(nil, v)
	if err != nil {
		return nil, fmt.Errorf("converting plain map into binary Avro format: %w", err)
	}
	return native, nil
}

func (enc *avroEncoder) Decode(subject string, data []byte, vPtr interface{}) error {
	codec, err := enc.finder.FindByNamespace(subject)
	if err != nil {
		return WrapError(err, "finding codec by subject %s", subject)
	}
	if codec == nil {
		return NewError("codec for subject %s not found", subject)
	}

	native, _, err := codec.NativeFromBinary(data)
	if err != nil {
		// not avro message
		vPtr = data
		return WrapError(err, "not avro message")
	}

	// avro message
	vPtr = native.(map[string]interface{})

	return nil
}
