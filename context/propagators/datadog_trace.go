package propagators

import (
	"context"

	workflow_context "temporal_microservices/context/workflow"

	commonpb "go.temporal.io/api/common/v1"
	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/workflow"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// WorkflowHeaderWriteAdapter writes to the workflow headers.
type WorkflowHeaderWriteAdapter struct {
	writer workflow.HeaderWriter
}

var _ tracer.TextMapWriter = WorkflowHeaderWriteAdapter{}

// Set sets the key/values to the workflow headers.
func (whwa WorkflowHeaderWriteAdapter) Set(key, val string) {
	encodedValue, err := converter.GetDefaultDataConverter().ToPayload(val)
	if err != nil {
		panic(err)
	}
	whwa.writer.Set(key, encodedValue)
}

type WorkflowHeaderReadAdapter struct {
	reader workflow.HeaderReader
}

var _ tracer.TextMapReader = WorkflowHeaderReadAdapter{}

// ForeachKey implements TextMapReader.
func (whra WorkflowHeaderReadAdapter) ForeachKey(handler func(key, val string) error) (err error) {
	err = whra.reader.ForEachKey(func(key string, value *commonpb.Payload) (err4eachKey error) {
		if string(value.GetMetadata()[converter.MetadataEncoding]) == converter.MetadataEncodingJSON {
			var decodedValue string
			err4eachKey = converter.GetDefaultDataConverter().FromPayload(value, &decodedValue)
			if err4eachKey != nil {
				return
			}
			return handler(key, decodedValue)
		}
		return
	})
	return
}

type dataDogTraceContextPropagator struct {
}

// NewDataDogTraceContextPropagator returns a new DataDog context propagator
func NewDataDogTraceContextPropagator() workflow.ContextPropagator {
	return &dataDogTraceContextPropagator{}
}

// Inject injects values from context into headers for propagation.
func (s *dataDogTraceContextPropagator) Inject(ctx context.Context, writer workflow.HeaderWriter) (err error) {
	span, ok := tracer.SpanFromContext(ctx)
	if ok {
		err = tracer.Inject(span.Context(), WorkflowHeaderWriteAdapter{writer: writer})
		if err != nil {
			return
		}
	}
	return nil
}

// InjectFromWorkflow injects the span context from the workflow context context into headers for propagation.
func (s *dataDogTraceContextPropagator) InjectFromWorkflow(workflowCtx workflow.Context, writer workflow.HeaderWriter) (err error) {
	value := workflowCtx.Value(workflow_context.DataDogSpanKey)
	if value != nil {
		spanContext, ok := value.(ddtrace.SpanContext)
		if ok {
			err := tracer.Inject(spanContext, WorkflowHeaderWriteAdapter{writer: writer})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Extract extracts DataDog span context from headers and the span into context.
func (s *dataDogTraceContextPropagator) Extract(ctx context.Context, reader workflow.HeaderReader) (ctxOut context.Context, err error) {
	spanContext, err := tracer.Extract(WorkflowHeaderReadAdapter{reader})
	if err != nil && err != tracer.ErrSpanContextNotFound {
		return
	}
	if spanContext == nil || err == tracer.ErrSpanContextNotFound {
		return ctx, nil
	}
	span, ctxOut := tracer.StartSpanFromContext(ctx, "propagators.Extract", tracer.ChildOf(spanContext))
	defer span.Finish()
	return
}

// ExtractToWorkflow extracts DataDog span context from headers and puts it into workflow context.
func (s *dataDogTraceContextPropagator) ExtractToWorkflow(ctx workflow.Context,
	reader workflow.HeaderReader) (ctxOut workflow.Context, err error) {
	spanContext, err := tracer.Extract(WorkflowHeaderReadAdapter{reader})
	if err != nil && err != tracer.ErrSpanContextNotFound {
		return
	}
	if spanContext == nil || err == tracer.ErrSpanContextNotFound {
		return ctx, nil //errors.Wrap(err, "no DatDog trace exists in the workflow headers")
	}
	ctxOut = workflow.WithValue(ctx, workflow_context.DataDogSpanKey, spanContext)
	return
}
