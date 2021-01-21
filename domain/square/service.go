package square

import (
	"context"
	"errors"
	"temporal_microservices"
	"temporal_microservices/domain"
)

var RectangleSquareActivityName = domain.GetActivityName(application{}.CalculateRectangleSquare)

type Rectangle struct {
	ID     string
	Length float64
	Width  float64
}

type Service interface {
	CalculateRectangleSquare(ctx context.Context, req CalculateRectangleSquareRequest) (resp CalculateRectangleSquareResponse, err error)
}

func MakeService(ctxReg ContextRegistrar) (application, error) {
	if ctxReg == nil {
		return application{}, errors.New("context propagator should not be nil")
	}
	return application{ctxReg: ctxReg}, nil
}

type application struct {
	ctxReg ContextRegistrar
}

type ContextRegistrar interface {
	Register(ctx context.Context)
}

type CalculateRectangleSquareRequest struct {
	Rectangles []Rectangle
}

type CalculateRectangleSquareResponse struct {
	Squares map[string]float64
}

func (s application) CalculateRectangleSquare(ctx context.Context, req CalculateRectangleSquareRequest) (resp CalculateRectangleSquareResponse, err error) {
	s.ctxReg.Register(ctx)
	heartbeat := domain.StartHeartbeat(ctx, temporal_microservices.HeartbeatIntervalSec)
	defer heartbeat.Stop()

	resp.Squares = make(map[string]float64, len(req.Rectangles))
	for _, r := range req.Rectangles {
		resp.Squares[r.ID] = r.Width * r.Length
	}
	return
}
