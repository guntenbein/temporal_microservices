package tracing

import (
	"context"
)

// Tracer is a simple contract to trace services. Uses context as a source of the data.
type Tracer interface {
	StartSpan(ctx context.Context, operationName string) (Span, context.Context)
}

// Span is a span abstraction.
type Span interface {
	Finish(err error)
}
