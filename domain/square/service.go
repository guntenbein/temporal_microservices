package square

import (
	"context"
	"temporal_microservices/domain"
)

var RectangleSquareActivityName = domain.GetActivityName(Service{}.CalculateRectangleSquare)

type Rectangle struct {
	ID     string
	Length float64
	Width  float64
}

type Service struct{}

type CalculateRectangleSquareRequest struct {
	Rectangles []Rectangle
}

type CalculateRectangleSquareResponse struct {
	Squares map[string]float64
}

func (s Service) CalculateRectangleSquare(_ context.Context, req CalculateRectangleSquareRequest) (resp CalculateRectangleSquareResponse, err error) {
	resp.Squares = make(map[string]float64, len(req.Rectangles))
	for _, r := range req.Rectangles {
		resp.Squares[r.ID] = r.Width * r.Length
	}
	return
}
