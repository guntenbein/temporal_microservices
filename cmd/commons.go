package cmd

import (
	"log"
	"net"
	"temporal_microservices"
	"temporal_microservices/context/propagators"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

func InitTemporalClient() client.Client {
	temporalClientOptions := client.Options{HostPort: net.JoinHostPort("localhost", "7233"),
		ContextPropagators: []workflow.ContextPropagator{
			propagators.NewStringMapPropagator([]string{temporal_microservices.ProcessIDContextField}),
			propagators.NewSecretPropagator(propagators.SecretPropagatorConfig{
				Keys:   []string{temporal_microservices.JWTContextField},
				Crypto: propagators.Base64Crypto{},
			}),
			propagators.NewDataDogTraceContextPropagator(),
		},
	}
	temporalClient, err := client.NewClient(temporalClientOptions)
	if err != nil {
		log.Fatal("cannot start temporal client: " + err.Error())
	}
	return temporalClient
}
