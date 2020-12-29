package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"temporal_microservices"
	"temporal_microservices/domain/workflow"

	"go.temporal.io/sdk/client"
)

func MakeFiguresHandleFunc(temporalClient client.Client) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx := context.Background()
		input, err := getInputParallelepipeds(req)
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
	if _, errWrite := rw.Write([]byte(err.Error())); errWrite != nil {
		log.Print("error writing the HTTP response: " + errWrite.Error())
	}
}

func getInputParallelepipeds(req *http.Request) (pp []workflow.Parallelepiped, err error) {
	defer func() {
		if closeErr := req.Body.Close(); closeErr != nil {
			log.Print("error closing HTTP body: " + closeErr.Error())
		}
	}()
	err = json.NewDecoder(req.Body).Decode(&pp)
	return
}

func executeWorkflow(ctx context.Context, temporalClient client.Client, input []workflow.Parallelepiped) (output []workflow.Parallelepiped, err error) {
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: temporal_microservices.FigureWorkflowQueue,
	}
	workflowReq := workflow.CalculateParallelepipedWorkflowRequest{Parallelepipeds: input}
	workflowRun, err := temporalClient.ExecuteWorkflow(ctx, workflowOptions, workflow.CalculateParallelepipedWorkflow, workflowReq)
	if err != nil {
		return
	}
	workflowResp := workflow.CalculateParallelepipedWorkflowResponse{}
	err = workflowRun.Get(ctx, &workflowResp)
	if err != nil {
		return
	}
	return workflowResp.Parallelepipeds, nil
}

func writeOutputFigures(rw http.ResponseWriter, output []workflow.Parallelepiped) (err error) {
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
