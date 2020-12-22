package workflow

import (
	"math"
	"temporal_microservices"
	"temporal_microservices/service"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type FigureWorkflowRequest struct {
	Figures []service.Figure
}

type FigureWorkflowResponse struct {
	Figures []service.Figure
}

func FigureWorkflow(ctx workflow.Context, req FigureWorkflowRequest) (resp FigureWorkflowResponse, err error) {
	if len(req.Figures) == 0 {
		err = BusinessError{"there are no figures to process"}
		return
	}

	if temporal_microservices.BatchSize <= 0 {
		err = BusinessError{"batch size cannot be less or equal to zero"}
		return
	}

	batchCount := BatchCount(len(req.Figures), temporal_microservices.BatchSize)

	selector := workflow.NewNamedSelector(ctx, "select-figures-batches")
	var errOneBatch error
	cancelCtx, cancelHandler := workflow.WithCancel(ctx)

	var count int
	var outputFigures = make([]service.Figure, 0, len(req.Figures))
	for i := 0; i < batchCount; i++ {
		batch := Batch(req.Figures, i, temporal_microservices.BatchSize)

		future := processFigures(cancelCtx, batch)

		selector.AddFuture(future, func(f workflow.Future) {
			respVolume := service.FigureResponse{}
			if err := f.Get(cancelCtx, &respVolume); err != nil {
				cancelHandler()
				errOneBatch = err
			} else {
				outputFigures = append(outputFigures, resp.Figures...)
			}
		})

		count++
	}

	for i := 0; i < count; i++ {
		selector.Select(ctx)
		if errOneBatch != nil {
			return FigureWorkflowResponse{}, errOneBatch
		}
	}
	return FigureWorkflowResponse{Figures: outputFigures}, nil
}

func processFigures(cancelCtx workflow.Context, batch []service.Figure) workflow.Future {
	future, settable := workflow.NewFuture(cancelCtx)
	workflow.Go(cancelCtx, func(ctx workflow.Context) {
		ctx = WithActivityOptions(ctx, temporal_microservices.SquareActivityQueue)
		respSquare := service.FigureResponse{}
		err := workflow.ExecuteActivity(ctx, service.SquareActivityName, service.FigureRequest{Figures: batch}).
			Get(ctx, &respSquare)
		if err != nil {
			settable.Set(nil, err)
			return
		}
		respVolume := service.FigureResponse{}
		err = workflow.ExecuteActivity(ctx, service.VolumeActivityName, service.FigureRequest{Figures: respSquare.Figures}).
			Get(ctx, &respVolume)
		settable.Set(respVolume, err)
	})
	return future
}

func WithActivityOptions(ctx workflow.Context, queue string) workflow.Context {
	ao := workflow.ActivityOptions{
		TaskQueue:              queue,
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    24 * time.Hour,
		HeartbeatTimeout:       time.Second * 5,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:        time.Second,
			BackoffCoefficient:     2.0,
			MaximumInterval:        time.Minute * 5,
			NonRetryableErrorTypes: []string{"BusinessError"},
		},
	}
	ctxOut := workflow.WithActivityOptions(ctx, ao)
	return ctxOut
}

func BatchCount(wholeSize, batchSize int) int {
	if batchSize <= 0 {
		return 0
	}
	return int(math.Ceil(float64(wholeSize) / float64(batchSize)))
}

func Batch(source []service.Figure, batchNumber, batchSize int) []service.Figure {
	start := batchNumber * batchSize
	end := start + batchSize
	atAll := len(source)
	if batchSize <= 0 || batchNumber < 0 || atAll < start {
		return []service.Figure{}
	}
	if end > atAll {
		end = atAll
	}
	return source[start:end]
}
