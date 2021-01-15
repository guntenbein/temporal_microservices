//nolint:staticcheck
package propagators

import (
	"context"
	"fmt"

	commonpb "go.temporal.io/api/common/v1"
	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/workflow"
)

type Crypto interface {
	Encrypt(data []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
}

type SecretPropagatorConfig struct {
	Keys   []string
	Crypto Crypto
}

type secretPropagator struct {
	keySet map[string]struct{}
	crypto Crypto
}

func NewSecretPropagator(config SecretPropagatorConfig) workflow.ContextPropagator {
	if config.Crypto == nil {
		panic(fmt.Sprintf("SecretPropagator cannot be initialised with the parameters: %+v", config))
	}
	keyMap := make(map[string]struct{}, len(config.Keys))
	for _, key := range config.Keys {
		keyMap[key] = struct{}{}
	}
	return &secretPropagator{
		keySet: keyMap,
		crypto: config.Crypto,
	}
}

// Inject injects values from context into headers for propagation
func (s *secretPropagator) Inject(ctx context.Context, writer workflow.HeaderWriter) error {
	for key := range s.keySet {
		if value, ok := ctx.Value(key).(string); ok {
			encryptedValue, err := s.crypto.Encrypt([]byte(value))
			if err != nil {
				return err
			}
			encodedValue, err := converter.GetDefaultDataConverter().ToPayload(encryptedValue)
			if err != nil {
				return err
			}
			writer.Set(key, encodedValue)
		}
	}
	return nil
}

// InjectFromWorkflow injects values from context into headers for propagation
func (s *secretPropagator) InjectFromWorkflow(ctx workflow.Context, writer workflow.HeaderWriter) error {
	for key := range s.keySet {
		value := ctx.Value(key)
		encodedValue, err := converter.GetDefaultDataConverter().ToPayload(value)
		if err != nil {
			return err
		}
		writer.Set(key, encodedValue)
	}
	return nil
}

// Extract extracts values from headers and puts them into context
func (s *secretPropagator) Extract(ctx context.Context, reader workflow.HeaderReader) (context.Context, error) {
	if err := reader.ForEachKey(func(key string, value *commonpb.Payload) error {
		if _, ok := s.keySet[key]; ok {
			var decodedValue []byte
			err := converter.GetDefaultDataConverter().FromPayload(value, &decodedValue)
			if err != nil {
				return err
			}
			decryptedValue, err := s.crypto.Decrypt(decodedValue)
			if err != nil {
				return err
			}
			ctx = context.WithValue(ctx, interface{}(key), decryptedValue)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return ctx, nil
}

// ExtractToWorkflow extracts values from headers and puts them into context
func (s *secretPropagator) ExtractToWorkflow(ctx workflow.Context, reader workflow.HeaderReader) (workflow.Context, error) {
	if err := reader.ForEachKey(func(key string, value *commonpb.Payload) error {
		if _, ok := s.keySet[key]; ok {
			var decodedValue []byte
			err := converter.GetDefaultDataConverter().FromPayload(value, &decodedValue)
			if err != nil {
				return err
			}
			ctx = workflow.WithValue(ctx, key, decodedValue)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return ctx, nil
}
