package temporal_microservices

const (
	SquareActivityQueue = "SquareActivityQueue"
	VolumeActivityQueue = "VolumeActivityQueue"
	FigureWorkflowQueue = "FigureWorkflowQueue"
	BatchSize           = 10

	MaxConcurrentSquareActivitySize = 10
	MaxConcurrentVolumeActivitySize = 10
	MaxConcurrentFigureWorkflowSize = 3
)
