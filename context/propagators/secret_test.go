//nolint:staticcheck
package propagators

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
	temporal_testsuite "go.temporal.io/sdk/testsuite"
	temporal_workflow "go.temporal.io/sdk/workflow"
)

type SecretPropagatorTestSuite struct {
	suite.Suite
	temporal_testsuite.WorkflowTestSuite
}

const (
	jwt = "jwt"
)

func TestSecretPropagatorSuite(t *testing.T) {
	s := SecretPropagatorTestSuite{}
	suite.Run(t, &s)
}

func (s *SecretPropagatorTestSuite) TestInject() {
	ctx := context.Background()
	jwtInput := "jwt-token"
	crypto := Base64Crypto{}
	expectedEncryptedJwt, err := crypto.Encrypt([]byte(jwtInput))
	s.NoError(err)

	propagator := NewSecretPropagator(SecretPropagatorConfig{
		Keys:   []string{jwt},
		Crypto: crypto,
	})
	headersRWStub := makeStubHeaderReaderWriter()

	ctx = context.WithValue(ctx, interface{}(jwt), jwtInput)
	err = propagator.Inject(ctx, headersRWStub)
	s.NoError(err)

	var headersJwtValue []byte
	err = headersRWStub.GetValue(jwt, &headersJwtValue)
	s.NoError(err)
	s.EqualValuesf(expectedEncryptedJwt, headersJwtValue, "the token was not properly encrypted")
}

func (s *SecretPropagatorTestSuite) TestExtract() {
	ctx := context.Background()
	jwtInput := "jwt-token"

	propagator := NewSecretPropagator(SecretPropagatorConfig{
		Keys:   []string{jwt},
		Crypto: Base64Crypto{},
	})
	headersRWStub := makeStubHeaderReaderWriter()

	ctx = context.WithValue(ctx, jwt, jwtInput)
	err := propagator.Inject(ctx, headersRWStub)
	s.NoError(err)
	ctxOut, err := propagator.Extract(ctx, headersRWStub)
	s.NoError(err)

	jwtOutput := ctxOut.Value(jwt).(string)
	s.Equal(jwtInput, jwtOutput, "the token was not passed properly")
}

// TODO: simplify the unit test
// This will be possible after https://github.com/temporalio/go-sdk/issues/190 is resolved.
// Or in case if we have a possibility to create a background workflow context (now it is in "internal" package).
func (s *SecretPropagatorTestSuite) TestSecretPassedToWorkflowCtx() {
	env := s.NewTestWorkflowEnvironment()

	env.ExecuteWorkflow(s.PassJWT2WorkflowCtx)

	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())
}

// Workflow context should be just a carrier for the secret fields - they should be taken
// from the headers as-it-is then passed to the workflow context and then injected to the
// headers as-it-is.
func (s *SecretPropagatorTestSuite) PassJWT2WorkflowCtx(ctx temporal_workflow.Context) error {
	jwtInput := "jwt-token"
	crypto := Base64Crypto{}
	encryptedJwt, err := crypto.Encrypt([]byte(jwtInput))
	if err != nil {
		return err
	}

	headersRWStub := makeStubHeaderReaderWriter()
	propagator := NewSecretPropagator(SecretPropagatorConfig{
		Keys:   []string{jwt},
		Crypto: crypto,
	})

	ctx = temporal_workflow.WithValue(ctx, jwt, encryptedJwt)
	err = propagator.InjectFromWorkflow(ctx, headersRWStub)
	if err != nil {
		return err
	}
	var headersJwtValue []byte
	err = headersRWStub.GetValue(jwt, &headersJwtValue)
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(encryptedJwt, headersJwtValue) {
		return errors.New("jwt was not passed properly (as it is) to the workflow headers")
	}

	// clear jwt field from the workflow context and read from the headers
	ctx = temporal_workflow.WithValue(ctx, jwt, nil)
	ctx, err = propagator.ExtractToWorkflow(ctx, headersRWStub)
	if err != nil {
		return err
	}
	outputJwtField := ctx.Value(jwt)
	if !reflect.DeepEqual(encryptedJwt, outputJwtField) {
		return errors.New("jwt was not extracted properly (as it is) to the workflow context")
	}
	return nil
}
