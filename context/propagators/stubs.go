package propagators

import (
	temporal_commonpb "go.temporal.io/api/common/v1"
	"go.temporal.io/sdk/converter"
)

type stubHeaderReaderWriter struct {
	storage map[string]*temporal_commonpb.Payload
}

func makeStubHeaderReaderWriter() stubHeaderReaderWriter {
	return stubHeaderReaderWriter{map[string]*temporal_commonpb.Payload{}}
}

func (s stubHeaderReaderWriter) Set(key string, value *temporal_commonpb.Payload) {
	s.storage[key] = value
}

func (s stubHeaderReaderWriter) Get(key string, decodedValuePtr interface{}) error {
	return converter.GetDefaultDataConverter().FromPayload(s.storage[key], decodedValuePtr)
}

func (s stubHeaderReaderWriter) ForEachKey(handler func(string, *temporal_commonpb.Payload) error) error {
	for key, value := range s.storage {
		if err := handler(key, value); err != nil {
			return err
		}
	}
	return nil
}
