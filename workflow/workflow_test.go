package workflow

import (
	"temporal_microservices/service"
	"testing"

	"github.com/stretchr/testify/assert"
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
	activities := service.FigureService{}

	figures := makeFigures(276)

	env.OnActivity(activities.Square, mock.Anything, mock.Anything).
		Return([]service.Figure{}, nil).Return(service.FigureResponse{}, nil).Times(28)
	env.OnActivity(activities.Volume, mock.Anything, mock.Anything).
		Return([]service.Figure{}, nil).Return(service.FigureResponse{}, nil).Times(28)

	env.ExecuteWorkflow(FigureWorkflow, FigureWorkflowRequest{figures})

	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())

	env.AssertExpectations(s.T())
}

func (s *IndexationWorkflowTestSuite) Test_WorkflowFail() {
	env := s.NewTestWorkflowEnvironment()

	figures := makeFigures(0)

	env.ExecuteWorkflow(FigureWorkflow, FigureWorkflowRequest{figures})

	s.True(env.IsWorkflowCompleted())
	s.Error(env.GetWorkflowError())

	env.AssertExpectations(s.T())
}

func makeFigures(count int) []service.Figure {
	out := make([]service.Figure, 0, count)
	for i := 0; i < count; i++ {
		out = append(out, service.Figure{Length: 10, Height: 10, Width: 10})
	}
	return out
}

func TestBatchCount(t *testing.T) {
	assert.EqualValues(t, 1, BatchCount(3, 10))
	assert.EqualValues(t, 10, BatchCount(99, 10))
	assert.EqualValues(t, 0, BatchCount(0, 10))
	assert.EqualValues(t, 0, BatchCount(99, 0))
}

func TestBatch(t *testing.T) {
	figures := makeFigures(6)
	outFigures := Batch(figures, 2, 2)
	assert.EqualValues(t, 2, len(outFigures))

	outFigures = Batch(figures, 2, 10)
	assert.EqualValues(t, 0, len(outFigures))

	outFigures = Batch(figures, -2, 3)
	assert.EqualValues(t, 0, len(outFigures))

	outFigures = Batch(figures, 2, 0)
	assert.EqualValues(t, 0, len(outFigures))

	outFigures = Batch(figures, 1, 5)
	assert.EqualValues(t, 1, len(outFigures))
}
