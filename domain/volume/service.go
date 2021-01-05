package volume

import (
	"context"
	"temporal_microservices"
	"temporal_microservices/domain"
)

var ParallelepipedVolumeActivityName = domain.GetActivityName(Service{}.CalculateParallelepipedVolume)

type Parallelepiped struct {
	ID     string
	Length float64
	Width  float64
	Height float64
}

type Service struct{}

type CalculateParallelepipedVolumeRequest struct {
	Parallelepipeds []Parallelepiped
}

type CalculateParallelepipedVolumeResponse struct {
	Volumes map[string]float64
}

func (s Service) CalculateParallelepipedVolume(ctx context.Context, req CalculateParallelepipedVolumeRequest) (resp CalculateParallelepipedVolumeResponse, err error) {
	heartbeat := domain.StartHeartbeat(ctx, temporal_microservices.HeartbeatIntervalSec)
	defer heartbeat.Stop()

	resp.Volumes = make(map[string]float64, len(req.Parallelepipeds))
	for _, p := range req.Parallelepipeds {
		volume := p.Width * p.Length * p.Height
		resp.Volumes[p.ID] = volume
	}
	return
}
