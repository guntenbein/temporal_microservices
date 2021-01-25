//nolint:staticcheck
package propagators

import (
	"context"

	commonpb "go.temporal.io/api/common/v1"
	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/workflow"
)

// stringMapPropagator propagates the list of keys-strings across a workflow,
// interpreting the values as strings.
type stringMapPropagator struct {
	keySet map[string]struct{}
}

// NewSecretPropagator returns a new string context propagator
func NewStringMapPropagator(keys []string) workflow.ContextPropagator {
	keyMap := make(map[string]struct{}, len(keys))
	for _, key := range keys {
		keyMap[key] = struct{}{}
	}
	return &stringMapPropagator{keyMap}
}

// Inject injects values from context into headers for propagation
func (s *stringMapPropagator) Inject(ctx context.Context, writer workflow.HeaderWriter) error {
	for key := range s.keySet {
		if value, ok := ctx.Value(key).(string); ok {
			encodedValue, err := converter.GetDefaultDataConverter().ToPayload(value)
			if err != nil {
				return err
			}
			writer.Set(key, encodedValue)
		}
	}
	return nil
}

// InjectFromWorkflow injects values from context into headers for propagation
func (s *stringMapPropagator) InjectFromWorkflow(ctx workflow.Context, writer workflow.HeaderWriter) error {
	for key := range s.keySet {
		if value, ok := ctx.Value(key).(string); ok {
			encodedValue, err := converter.GetDefaultDataConverter().ToPayload(value)
			if err != nil {
				return err
			}
			writer.Set(key, encodedValue)
		}
	}
	return nil
}

// Extract extracts values from headers and puts them into context
func (s *stringMapPropagator) Extract(ctx context.Context, reader workflow.HeaderReader) (context.Context, error) {
	if err := reader.ForEachKey(func(key string, value *commonpb.Payload) error {
		if _, ok := s.keySet[key]; ok {
			var decodedValue string
			err := converter.GetDefaultDataConverter().FromPayload(value, &decodedValue)
			if err != nil {
				return err
			}
			ctx = context.WithValue(ctx, key, decodedValue)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return ctx, nil
}

// ExtractToWorkflow extracts values from headers and puts them into context
func (s *stringMapPropagator) ExtractToWorkflow(ctx workflow.Context, reader workflow.HeaderReader) (workflow.Context, error) {
	if err := reader.ForEachKey(func(key string, value *commonpb.Payload) error {
		if _, ok := s.keySet[key]; ok {
			var decodedValue string
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
