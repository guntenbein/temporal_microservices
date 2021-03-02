package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"temporal_microservices"
	"temporal_microservices/cmd"
	"temporal_microservices/controller"
	"temporal_microservices/domain/workflow"
	"temporal_microservices/tracing"
	"time"

	"github.com/gorilla/mux"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	datatdog_tracing "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const svcName = "figure-workflow-service"

func main() {
	log.Print("starting FIGURE WORKFLOW microservice")

	temporalClient := cmd.InitTemporalClient()

	datatdog_tracing.Start(
		datatdog_tracing.WithEnv("local"),
		datatdog_tracing.WithServiceName(svcName),
	)
	defer datatdog_tracing.Stop()

	tracer := tracing.DataDogTracer{}

	worker := initWorkflowWorker(temporalClient)

	httpServer := initHTTPServer(controller.MakeFiguresHandleFunc(temporalClient, tracer))

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-signals

	err := httpServer.Shutdown(context.Background())
	if err != nil {
		log.Fatal("cannot gracefully stop HTTP server: " + err.Error())
	}

	worker.Stop()

	log.Print("closing FIGURE WORKFLOW microservice")
}

func initWorkflowWorker(temporalClient client.Client) worker.Worker {
	workerOptions := worker.Options{
		MaxConcurrentWorkflowTaskExecutionSize: temporal_microservices.MaxConcurrentFigureWorkflowSize,
	}
	worker := worker.New(temporalClient, temporal_microservices.FigureWorkflowQueue, workerOptions)
	worker.RegisterWorkflow(workflow.CalculateParallelepipedWorkflow)

	err := worker.Start()
	if err != nil {
		log.Fatal("cannot start temporal worker: " + err.Error())
	}

	return worker
}

func initHTTPServer(singleHandler func(http.ResponseWriter, *http.Request)) *http.Server {
	router := mux.NewRouter()
	router.Methods(http.MethodPost).Path("/").HandlerFunc(singleHandler)
	server := &http.Server{
		Addr:         net.JoinHostPort("localhost", "8080"),
		Handler:      router,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("cannot start HTTP server: " + err.Error())
		}
	}()
	return server
}
