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
	"temporal_microservices/context/propagators"
	"temporal_microservices/controller"
	"temporal_microservices/domain/workflow"
	"time"

	"github.com/gorilla/mux"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	temporal_workflow "go.temporal.io/sdk/workflow"
)

func main() {
	log.Print("starting FIGURE WORKFLOW microservice")

	temporalClient := initTemporalClient()

	worker := initWorkflowWorker(temporalClient)

	httpServer := initHTTPServer(controller.MakeFiguresHandleFunc(temporalClient))

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

func initTemporalClient() client.Client {
	temporalClientOptions := client.Options{HostPort: net.JoinHostPort("localhost", "7233"),
		ContextPropagators: []temporal_workflow.ContextPropagator{
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
