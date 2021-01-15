//nolint:staticcheck
package propagators

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	temporal_testsuite "go.temporal.io/sdk/testsuite"
	temporal_workflow "go.temporal.io/sdk/workflow"
)

type StringMapPropagatorTestSuite struct {
	suite.Suite
	temporal_testsuite.WorkflowTestSuite
	propagator temporal_workflow.ContextPropagator
}

func TestStringMapPropagatorSuite(t *testing.T) {
	s := StringMapPropagatorTestSuite{propagator: NewStringMapPropagator([]string{"aaa", "bbb"})}
	suite.Run(t, &s)
}

func (s *StringMapPropagatorTestSuite) TestOrdinaryContextPassedFields() {
	headersRWStub := makeStubHeaderReaderWriter()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "aaa", "trololo")
	ctx = context.WithValue(ctx, "bbb", "trololo2")
	err := s.propagator.Inject(ctx, headersRWStub)
	s.NoError(err)
	ctxOut, err := s.propagator.Extract(context.Background(), headersRWStub)
	s.NoError(err)

	s.NotSame("trololo", ctxOut.Value("aaa"), "propagator does not propagate required fields")
	s.NotSame("trololo2", ctxOut.Value("bbb"), "propagator does not propagate required fields")
}

func (s *StringMapPropagatorTestSuite) TestOrdinaryContextNOTPassedFields() {
	headersRWStub := makeStubHeaderReaderWriter()

	ctx := context.Background()
	ctx = context.WithValue(ctx, "aaa", "trololo")
	ctx = context.WithValue(ctx, "bbb", "trololo")
	ctx = context.WithValue(ctx, "ccc", "trololo")
	err := s.propagator.Inject(ctx, headersRWStub)
	s.NoError(err)
	ctxOut, err := s.propagator.Extract(context.Background(), headersRWStub)
	s.NoError(err)

	ccc := ctxOut.Value("ccc")
	s.Nil(ccc, "propagated field that was not required for the propagator")
}

// TODO: simplify the unit test
// This will be possible after https://github.com/temporalio/go-sdk/issues/190 is resolved.
// Or in case if we have a possibility to create a background workflow context (now it is in "internal" package).
func (s *StringMapPropagatorTestSuite) TestStringMapPropagatorWorkflowCtx() {
	env := s.NewTestWorkflowEnvironment()

	env.ExecuteWorkflow(s.CheckContextWorkflow)

	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())
}

func (s *StringMapPropagatorTestSuite) CheckContextWorkflow(ctx temporal_workflow.Context) error {
	propagator := NewStringMapPropagator([]string{"aaa", "bbb"})
	headersRWStub := makeStubHeaderReaderWriter()

	ctx = temporal_workflow.WithValue(ctx, "aaa", "trololo")
	ctx = temporal_workflow.WithValue(ctx, "bbb", "trololo2")
	ctx = temporal_workflow.WithValue(ctx, "ccc", "trololo")
	err := propagator.InjectFromWorkflow(ctx, headersRWStub)
	if err != nil {
		return err
	}
	// clear context in order to take the values from the headers
	ctx = temporal_workflow.WithValue(ctx, "aaa", "")
	ctx = temporal_workflow.WithValue(ctx, "bbb", "")
	ctx = temporal_workflow.WithValue(ctx, "ccc", "")
	ctxOut, err := propagator.ExtractToWorkflow(ctx, headersRWStub)
	if err != nil {
		return err
	}

	aaa := ctxOut.Value("aaa")
	bbb := ctxOut.Value("bbb")
	if aaa == nil || aaa != "trololo" || bbb == nil || bbb != "trololo2" {
		return errors.New("propagator does not propagate required fields")
	}

	ccc := ctxOut.Value("ccc")
	if ccc != "" {
		return errors.New("propagated field that was not required for the propagator")
	}
	return nil
}
