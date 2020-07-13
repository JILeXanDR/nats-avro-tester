package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
)

type avroEncoder struct {
	finder CodecStorage
	logger *Logger
}

func NewAvroEncoder(finder CodecStorage, logger *Logger) (nats.Encoder, error) {
	return &avroEncoder{
		finder: finder,
		logger: logger,
	}, nil
}

func (enc *avroEncoder) Encode(subject string, v interface{}) ([]byte, error) {
	b, err := enc.encode(subject, v)
	if err != nil {
		enc.logger.Err(err).Msg("encoding")
	} else {
		enc.logger.Debug().Bytes("val", b).Int("len", len(b)).Str("raw", fmt.Sprintf("%v", b)).Msg("encoded message")
	}
	return b, err
}

func (enc *avroEncoder) encode(subject string, v interface{}) ([]byte, error) {
	enc.logger.Debug().Str("subject", subject).Msg("encoding message data")
	codec, err := enc.finder.FindByNamespace(subject)
	if err != nil {
		return nil, WrapError(err, "finding codec by subject %s", subject)
	}
	if codec == nil {
		return nil, NewError("codec for subject %s not found", subject)
	}

	native, err := codec.BinaryFromNative(nil, v)
	if err != nil {
		return nil, WrapError(err, "converting plain map to binary Avro data")
	}
	return native, nil
}

func (enc *avroEncoder) Decode(subject string, data []byte, vPtr interface{}) error {
	err := enc.decode(subject, data, vPtr)
	if err != nil {
		enc.logger.Err(err).Msg("decoding")
	} else {
		enc.logger.Debug().Interface("val", vPtr).Msg("decoded message")
	}
	return err
}

func (enc *avroEncoder) decode(subject string, data []byte, vPtr interface{}) error {
	enc.logger.Debug().Str("subject", subject).Bytes("data", data).Msg("decoding message data")
	codec, err := enc.finder.FindByNamespace(subject)
	if err != nil {
		return WrapError(err, `finding codec by subject "%s"`, subject)
	}
	if codec == nil {
		return NewError(`codec for subject "%s" not found`, subject)
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
