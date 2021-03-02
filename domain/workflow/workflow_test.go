package workflow

import (
	"temporal_microservices/context"
	"temporal_microservices/domain/square"
	"temporal_microservices/domain/volume"
	"temporal_microservices/tracing"
	"testing"

	"github.com/leonelquinteros/gorand"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
)

type IndexationWorkflowTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
}

func TestIndexationWorkflowTestSuite(t *testing.T) {
	suite.Run(t, new(IndexationWorkflowTestSuite))
}

func (s *IndexationWorkflowTestSuite) Test_WorkflowSuccess() {
	env := s.NewTestWorkflowEnvironment()
	squareService, err := square.MakeService(context.SimpleRegistrar{}, tracing.NoopTracer{})
	s.NoError(err)
	env.RegisterActivity(squareService.CalculateRectangleSquare)
	volumeService := volume.MakeService(tracing.NoopTracer{})
	env.RegisterActivity(volumeService.CalculateParallelepipedVolume)

	pp, err := makeParallelepipeds(276, true)
	s.NoError(err)

	env.OnActivity(squareService.CalculateRectangleSquare, mock.Anything, mock.Anything).
		Return([]square.Rectangle{}, nil).Return(square.CalculateRectangleSquareResponse{}, nil).Times(28)
	env.OnActivity(volumeService.CalculateParallelepipedVolume, mock.Anything, mock.Anything).
		Return([]square.Rectangle{}, nil).Return(volume.CalculateParallelepipedVolumeResponse{}, nil).Times(28)

	env.ExecuteWorkflow(CalculateParallelepipedWorkflow, CalculateParallelepipedWorkflowRequest{BatchSize: 10, Parallelepipeds: pp})

	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())

	env.AssertExpectations(s.T())
}

func (s *IndexationWorkflowTestSuite) Test_WorkflowFailNoInput() {
	env := s.NewTestWorkflowEnvironment()
	squareService, err := square.MakeService(context.SimpleRegistrar{}, tracing.NoopTracer{})
	s.NoError(err)
	env.RegisterActivity(squareService.CalculateRectangleSquare)
	volumeService := volume.MakeService(tracing.NoopTracer{})
	env.RegisterActivity(volumeService.CalculateParallelepipedVolume)

	pp, err := makeParallelepipeds(0, true)
	s.NoError(err)

	env.ExecuteWorkflow(CalculateParallelepipedWorkflow, CalculateParallelepipedWorkflowRequest{BatchSize: 10, Parallelepipeds: pp})

	s.True(env.IsWorkflowCompleted())
	s.Error(env.GetWorkflowError())

	env.AssertExpectations(s.T())
}

func (s *IndexationWorkflowTestSuite) Test_WorkflowFailNoIDs() {
	env := s.NewTestWorkflowEnvironment()
	squareService, err := square.MakeService(context.SimpleRegistrar{}, tracing.NoopTracer{})
	s.NoError(err)
	env.RegisterActivity(squareService.CalculateRectangleSquare)
	volumeService := volume.MakeService(tracing.NoopTracer{})
	env.RegisterActivity(volumeService.CalculateParallelepipedVolume)

	pp, err := makeParallelepipeds(276, false)
	s.NoError(err)

	env.ExecuteWorkflow(CalculateParallelepipedWorkflow, CalculateParallelepipedWorkflowRequest{BatchSize: 10, Parallelepipeds: pp})

	s.True(env.IsWorkflowCompleted())
	s.Error(env.GetWorkflowError())

	env.AssertExpectations(s.T())
}

func makeParallelepipeds(count int, withIDs bool) (out []Parallelepiped, err error) {
	out = make([]Parallelepiped, 0, count)
	for i := 0; i < count; i++ {
		p := Parallelepiped{Length: 10, Height: 10, Width: 10}
		if withIDs {
			p.ID, err = generateID()
			if err != nil {
				return
			}
		}
		out = append(out, p)
	}
	return out, nil
}
func generateID() (string, error) {
	uuid, err := gorand.UUIDv4()
	if err != nil {
		return "", err
	}
	uuidStr, err := gorand.MarshalUUID(uuid)
	if err != nil {
		return "", err
	}
	return uuidStr, nil
}
