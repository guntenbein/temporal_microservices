package tracing

import (
	"context"
)

var _ Tracer = NoopTracer{}
var _ Span = NoopSpan{}

// NoopTracer is a datadog Tracer implementation.
type NoopTracer struct {
}

// NoopSpan encapsulates datadog's Tracer.Span.
type NoopSpan struct {
}

func MakeNoopSpan() NoopSpan {
	return NoopSpan{}
}

// StartSpan start new span using provided context.
func (t NoopTracer) StartSpan(ctx context.Context, _ string) (Span, context.Context) {
	return NoopSpan{}, ctx
}

// Finish closes the span.
func (s NoopSpan) Finish(_ error) {
}
