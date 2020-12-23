package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"temporal_microservices"
	"temporal_microservices/service"
	"temporal_microservices/workflow"

	"go.temporal.io/sdk/client"
)

func MakeFiguresHandleFunc(temporalClient client.Client) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx := context.Background()
		input, err := getInputFigures(req)
		if err != nil {
			writeError(rw, err)
			return
		}
		output, err := executeWorkflow(ctx, temporalClient, input)
		if err != nil {
			writeError(rw, err)
			return
		}
		err = writeOutputFigures(rw, output)
		if err != nil {
			writeError(rw, err)
			return
		}
	}
}

func writeError(rw http.ResponseWriter, err error) {
	log.Print(err.Error())
	rw.WriteHeader(http.StatusInternalServerError)
	rw.Write([]byte(err.Error()))
}

func getInputFigures(req *http.Request) (figures []service.Figure, err error) {
	defer func() {
		if closeErr := req.Body.Close(); closeErr != nil {
			log.Print("error closing HTTP body: " + closeErr.Error())
		}
	}()
	err = json.NewDecoder(req.Body).Decode(&figures)
	return
}

func executeWorkflow(ctx context.Context, temporalClient client.Client, input []service.Figure) (output []service.Figure, err error) {
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: temporal_microservices.FigureWorkflowQueue,
	}
	workflowReq := workflow.FigureWorkflowRequest{
		Figures: input,
	}
	workflowRun, err := temporalClient.ExecuteWorkflow(ctx, workflowOptions, workflow.FigureWorkflow, workflowReq)
	if err != nil {
		return
	}
	workflowResp := workflow.FigureWorkflowResponse{}
	err = workflowRun.Get(ctx, &workflowResp)
	if err != nil {
		return
	}
	return workflowResp.Figures, nil
}

func writeOutputFigures(rw http.ResponseWriter, output []service.Figure) (err error) {
	body, err := json.Marshal(output)
	if err != nil {
		return
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	_, err = rw.Write(body)
	if err != nil {
		return
	}
	return
}
