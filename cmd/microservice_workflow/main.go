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
	"temporal_microservices/controller"
	"temporal_microservices/workflow"
	"time"

	"github.com/gorilla/mux"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	log.Print("starting FIGURE WORKFLOW microservice")

	temporalClient := initTemporalClient()

	worker := initWorkflowWorker(temporalClient)

	httpServer := initHTTPServer(controller.MakeFiguresHandleFunc(temporalClient))

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	<-signals

	err := httpServer.Shutdown(context.Background())
	if err != nil {
		log.Fatal("cannot gracefully stop HTTP server: " + err.Error())
	}

	worker.Stop()

	log.Print("closing FIGURE WORKFLOW microservice")
}

func initTemporalClient() client.Client {
	temporalClientOptions := client.Options{HostPort: net.JoinHostPort("localhost", "7233")}
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
	worker.RegisterWorkflow(workflow.FigureWorkflow)

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
