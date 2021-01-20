package temporal_microservices

const (
	SquareActivityQueue = "SquareActivityQueue"
	VolumeActivityQueue = "VolumeActivityQueue"
	FigureWorkflowQueue = "FigureWorkflowQueue"

	MaxConcurrentSquareActivitySize = 10
	MaxConcurrentVolumeActivitySize = 10
	MaxConcurrentFigureWorkflowSize = 3

	HeartbeatIntervalSec = 1

	ProcessIDContextField = "process-id"
	JWTContextField       = "jwt"

	ProcessIDHTTPHeader     = "process-id"
	AuthorizationHTTPHeader = "authorization"
)
