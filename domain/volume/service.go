package volume

import (
	"context"
	"temporal_microservices"
	"temporal_microservices/domain"
	"temporal_microservices/tracing"
)

var ParallelepipedVolumeActivityName = domain.GetActivityName(application{}.CalculateParallelepipedVolume)

type Parallelepiped struct {
	ID     string
	Length float64
	Width  float64
	Height float64
}

type Service interface {
	CalculateParallelepipedVolume(ctx context.Context, req CalculateParallelepipedVolumeRequest) (resp CalculateParallelepipedVolumeResponse, err error)
}

type application struct {
	tracer tracing.Tracer
}

func MakeService(tracer tracing.Tracer) application {
	return application{tracer: tracer}
}

type CalculateParallelepipedVolumeRequest struct {
	Parallelepipeds []Parallelepiped
}

type CalculateParallelepipedVolumeResponse struct {
	Volumes map[string]float64
}

func (s application) CalculateParallelepipedVolume(ctx context.Context, req CalculateParallelepipedVolumeRequest) (resp CalculateParallelepipedVolumeResponse, err error) {
	span, ctx := s.tracer.StartSpan(ctx, "Volume")
	defer func() {
		span.Finish(err)
	}()

	heartbeat := domain.StartHeartbeat(ctx, temporal_microservices.HeartbeatIntervalSec)
	defer heartbeat.Stop()

	resp.Volumes = make(map[string]float64, len(req.Parallelepipeds))
	for _, p := range req.Parallelepipeds {
		volume := p.Width * p.Length * p.Height
		resp.Volumes[p.ID] = volume
	}
	// time.Sleep(5*time.Second)
	return
}
