package tracing

import (
	"context"
	"temporal_microservices"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var _ Tracer = DataDogTracer{}
var _ Span = DataDogSpan{}

// DataDogTracer is a datadog Tracer implementation.
type DataDogTracer struct {
}

// DataDogSpan encapsulates datadog's Tracer.Span.
type DataDogSpan struct {
	span ddtracer.Span
}

func MakeDatDogSpan(span ddtracer.Span) DataDogSpan {
	return DataDogSpan{span}
}

// StartSpan start new span using provided context.
func (t DataDogTracer) StartSpan(ctx context.Context, spanName string) (Span, context.Context) {
	dataDogSpan, updatedCtx := ddtracer.StartSpanFromContext(ctx, spanName)
	return t.wrap(updatedCtx, dataDogSpan)
}

func (t DataDogTracer) wrap(updatedCtx context.Context, span ddtracer.Span) (DataDogSpan, context.Context) {
	s := DataDogSpan{span: span}
	s.span.SetTag(temporal_microservices.ProcessIDContextField, updatedCtx.Value(temporal_microservices.ProcessIDContextField))
	return s, updatedCtx
}

// Finish closes the span.
func (s DataDogSpan) Finish(err error) {
	s.span.Finish(func(cfg *ddtrace.FinishConfig) {
		cfg.Error = err
	})
}
