package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"temporal_microservices"
	"temporal_microservices/cmd"
	"temporal_microservices/context"
	"temporal_microservices/domain/square"
	"temporal_microservices/tracing"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	dtatdog_tracing "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	log.Print("starting SQUARE microservice")

	temporalClient := cmd.InitTemporalClient()

	dtatdog_tracing.Start(
		dtatdog_tracing.WithEnv("local"),
		dtatdog_tracing.WithServiceName("square-service"),
	)
	defer dtatdog_tracing.Stop()

	tracer := tracing.DataDogTracer{}

	worker := initActivityWorker(temporalClient, tracer)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-signals

	worker.Stop()

	log.Print("closing SQUARE microservice")
}

func initActivityWorker(temporalClient client.Client, tracer tracing.Tracer) worker.Worker {
	workerOptions := worker.Options{
		MaxConcurrentActivityExecutionSize: temporal_microservices.MaxConcurrentSquareActivitySize,
	}
	worker := worker.New(temporalClient, temporal_microservices.SquareActivityQueue, workerOptions)
	handler, err := square.MakeService(context.SimpleRegistrar{}, tracer)
	if err != nil {
		log.Fatal("cannot init the service: " + err.Error())
	}
	worker.RegisterActivity(handler.CalculateRectangleSquare)

	err = worker.Start()
	if err != nil {
		log.Fatal("cannot start temporal worker: " + err.Error())
	}
	return worker
}
