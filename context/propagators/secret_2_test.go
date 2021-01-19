//nolint:staticcheck
package propagators

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	commonpb "go.temporal.io/api/common/v1"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/converter"
	temporal_testsuite "go.temporal.io/sdk/testsuite"
	temporal_workflow "go.temporal.io/sdk/workflow"
)

const jwtValue = "sample_jwt"

var encryptedJwt, _ = Base64Crypto{}.Encrypt([]byte(jwtValue))

type SecretPropagatorTestSuite2 struct {
	suite.Suite
	temporal_testsuite.WorkflowTestSuite
}

func TestSecretPropagatorSuite2(t *testing.T) {
	s := SecretPropagatorTestSuite2{}

	payload, _ := converter.GetDefaultDataConverter().ToPayload(encryptedJwt)
	s.SetHeader(&commonpb.Header{
		Fields: map[string]*commonpb.Payload{
			jwt: payload,
		},
	})

	suite.Run(t, &s)
}

func (s *SecretPropagatorTestSuite2) Test_SecretPropagation() {
	handler := &Handler{checkCtx: func(workflowCtx temporal_workflow.Context) error {
		workflowJWTVal := workflowCtx.Value(jwt)
		if !reflect.DeepEqual(encryptedJwt, workflowJWTVal) {
			return errors.New("jwt is not encrypted at the workflow context")
		}
		return nil
	}}

	env := s.NewTestWorkflowEnvironment()
	env.SetContextPropagators([]temporal_workflow.ContextPropagator{NewSecretPropagator(SecretPropagatorConfig{
		Keys:   []string{jwt},
		Crypto: Base64Crypto{},
	})})
	env.RegisterActivity(handler.DummyActivity)

	var activityJWTValue string
	env.SetOnActivityStartedListener(func(activityInfo *activity.Info, ctx context.Context, args converter.EncodedValues) {
		activityJWTValue = ctx.Value(jwt).(string)
	})

	env.OnActivity(handler.DummyActivity, mock.Anything).Return(nil).Once()

	env.ExecuteWorkflow(handler.CtxCHeckWorkflow)
	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())

	s.Equal(jwtValue, activityJWTValue)

	env.AssertExpectations(s.T())
}

type Handler struct {
	checkCtx func(workflowCtx temporal_workflow.Context) error
}

func (h *Handler) DummyActivity(ctx context.Context) error {
	return nil
}

func (h *Handler) CtxCHeckWorkflow(ctx temporal_workflow.Context) error {
	ao := temporal_workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
	}
	activityCtx := temporal_workflow.WithActivityOptions(ctx, ao)
	if err := temporal_workflow.ExecuteActivity(activityCtx, h.DummyActivity).Get(activityCtx, nil); err != nil {
		return err
	}
	return h.checkCtx(ctx)
}
