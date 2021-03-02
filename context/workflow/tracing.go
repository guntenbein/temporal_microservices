package workflow

import (
	"context"
	"temporal_microservices/tracing"

	"go.temporal.io/sdk/workflow"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type dataDogSpanKey struct{}

var DataDogSpanKey dataDogSpanKey

func StartSpan(workflowCtx workflow.Context, operation string) (workflowCtxOut workflow.Context, span tracing.Span) {
	value := workflowCtx.Value(DataDogSpanKey)
	if value != nil {
		spanContext, ok := value.(ddtrace.SpanContext)
		if ok {
			span, _ := ddtracer.StartSpanFromContext(context.Background(), operation, ddtracer.ChildOf(spanContext))
			newSpanContext := span.Context()
			newWorkflowCtx := workflow.WithValue(workflowCtx, DataDogSpanKey, newSpanContext)
			newAppSpan := tracing.MakeDatDogSpan(span)
			return newWorkflowCtx, newAppSpan
		}
	}
	return workflowCtx, tracing.MakeNoopSpan()
}
