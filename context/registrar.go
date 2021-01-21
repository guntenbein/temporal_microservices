package context

import (
	"context"
	"log"
	"temporal_microservices"
)

type SimpleRegistrar struct{}

func (s SimpleRegistrar) Register(ctx context.Context) {
	processID := ctx.Value(temporal_microservices.ProcessIDContextField)
	jwt := ctx.Value(temporal_microservices.JWTContextField)
	log.Printf("Registered context values:\n* process id : %s\n* jwt        : %s\n", processID, jwt)
}
