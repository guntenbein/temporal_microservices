package service

import "context"

type Figure struct {
	Length int
	Width  int
	Height int
	Square int
	Volume int
}

type FigureService struct{}

type FigureRequest struct {
	figures []Figure
}

type FigureResponse struct {
	figures []Figure
}

func (s FigureService) Square(_ context.Context, req FigureRequest) (resp FigureResponse, err error) {
	if req.figures == nil {
		return
	}
	resp.figures = make([]Figure, 0, len(req.figures))
	for _, f := range req.figures {
		f.Square = f.Width * f.Length
		resp.figures = append(resp.figures, f)
	}
	return
}

func (s FigureService) Volume(_ context.Context, req FigureRequest) (resp FigureResponse, err error) {
	if req.figures == nil {
		return
	}
	resp.figures = make([]Figure, 0, len(req.figures))
	for _, f := range req.figures {
		f.Volume = f.Square * f.Height
		resp.figures = append(resp.figures, f)
	}
	return
}
