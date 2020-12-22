package service

import (
	"context"
	"reflect"
	"runtime"
	"strings"
)

var SquareActivityName = GetActivityName(FigureService{}.Square)
var VolumeActivityName = GetActivityName(FigureService{}.Volume)

type Figure struct {
	Length int
	Width  int
	Height int
	Square int
	Volume int
}

type FigureService struct{}

type FigureRequest struct {
	Figures []Figure
}

type FigureResponse struct {
	Figures []Figure
}

func (s FigureService) Square(_ context.Context, req FigureRequest) (resp FigureResponse, err error) {
	resp.Figures = make([]Figure, 0, len(req.Figures))
	for _, f := range req.Figures {
		f.Square = f.Width * f.Length
		resp.Figures = append(resp.Figures, f)
	}
	return
}

func (s FigureService) Volume(_ context.Context, req FigureRequest) (resp FigureResponse, err error) {
	resp.Figures = make([]Figure, 0, len(req.Figures))
	for _, f := range req.Figures {
		f.Volume = f.Square * f.Height
		resp.Figures = append(resp.Figures, f)
	}
	return
}

func GetActivityName(i interface{}) string {
	fullName := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	parts := strings.Split(fullName, ".")
	lastPart := parts[len(parts)-1]
	parts2 := strings.Split(lastPart, "-")
	firstPart := parts2[0]
	return firstPart
}
