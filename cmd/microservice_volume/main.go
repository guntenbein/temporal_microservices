package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"temporal_microservices"
	"temporal_microservices/context/propagators"
	"temporal_microservices/domain/volume"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func main() {
	log.Print("starting VOLUME microservice")

	temporalClient := initTemporalClient()

	worker := initActivityWorker(temporalClient)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-signals

	worker.Stop()

	log.Print("closing VOLUME microservice")
}

func initTemporalClient() client.Client {
	temporalClientOptions := client.Options{HostPort: net.JoinHostPort("localhost", "7233"),
		ContextPropagators: []workflow.ContextPropagator{
			propagators.NewStringMapPropagator([]string{temporal_microservices.ProcessIDContextField}),
			propagators.NewSecretPropagator(propagators.SecretPropagatorConfig{
				Keys:   []string{temporal_microservices.JWTContextField},
				Crypto: propagators.Base64Crypto{},
			}),
		},
	}
	temporalClient, err := client.NewClient(temporalClientOptions)
	if err != nil {
		log.Fatal("cannot start temporal client: " + err.Error())
	}
	return temporalClient
}

func initActivityWorker(temporalClient client.Client) worker.Worker {
	workerOptions := worker.Options{
		MaxConcurrentActivityExecutionSize: temporal_microservices.MaxConcurrentVolumeActivitySize,
	}
	worker := worker.New(temporalClient, temporal_microservices.VolumeActivityQueue, workerOptions)
	worker.RegisterActivity(volume.Service{}.CalculateParallelepipedVolume)

	err := worker.Start()
	if err != nil {
		log.Fatal("cannot start temporal worker: " + err.Error())
	}
	return worker
}
