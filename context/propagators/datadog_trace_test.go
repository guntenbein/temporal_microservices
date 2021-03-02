package propagators

import (
	"context"
	"errors"
	workflow_context "temporal_microservices/context/workflow"
	"testing"

	"github.com/stretchr/testify/suite"
	temporal_testsuite "go.temporal.io/sdk/testsuite"
	temporal_workflow "go.temporal.io/sdk/workflow"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type TracePropagatorTestSuite struct {
	suite.Suite
	temporal_testsuite.WorkflowTestSuite
	propagator temporal_workflow.ContextPropagator
}

func TestTracePropagator(t *testing.T) {
	s := TracePropagatorTestSuite{propagator: NewDataDogTraceContextPropagator()}
	suite.Run(t, &s)
}

func (s *TracePropagatorTestSuite) TestOrdinaryContextInjectExtractLoop() {
	headersRWStub := makeStubHeaderReaderWriter()
	ctx := context.Background()
	_, ctx = tracer.StartSpanFromContext(ctx, "first operation")
	err := s.propagator.Inject(ctx, headersRWStub)
	s.NoError(err)
	ctxOut, err := s.propagator.Extract(context.Background(), headersRWStub)
	s.NoError(err)
	_, ok := tracer.SpanFromContext(ctxOut)
	s.True(ok, "cannot find output span")
}

func (s *TracePropagatorTestSuite) TestWorkflowContextInjectExtractLoop() {
	env := s.NewTestWorkflowEnvironment()

	env.ExecuteWorkflow(s.CheckTraceWorkflow)

	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())
}

func (s *TracePropagatorTestSuite) CheckTraceWorkflow(ctx temporal_workflow.Context) error {
	headersRWStub := makeStubHeaderReaderWriter()
	ctxOrdinary := context.Background()
	span, _ := tracer.StartSpanFromContext(ctxOrdinary, "first operation")
	ctx = temporal_workflow.WithValue(ctx, workflow_context.DataDogSpanKey, span.Context())
	err := s.propagator.InjectFromWorkflow(ctx, headersRWStub)
	if err != nil {
		return err
	}
	ctx = temporal_workflow.WithValue(ctx, workflow_context.DataDogSpanKey, nil)
	ctx, err = s.propagator.ExtractToWorkflow(ctx, headersRWStub)
	if err != nil {
		return err
	}
	spanContext := ctx.Value(workflow_context.DataDogSpanKey)
	if spanContext == nil {
		return errors.New("cannot find output span context")
	}
	return nil
}
