package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"temporal_microservices"
	"temporal_microservices/domain/workflow"
	"temporal_microservices/tracing"

	"go.temporal.io/sdk/client"
)

const bearerPrefix = "bearer"

func MakeFiguresHandleFunc(temporalClient client.Client, tracer tracing.Tracer) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx := context.Background()
		var err error
		span, ctx := tracer.StartSpan(ctx, "HTTP starter for FiguresWorkflow")
		defer func() {
			span.Finish(err)
		}()
		wReq, err := getWorkflowRequest(req)
		if err != nil {
			writeError(rw, err)
			return
		}
		ctxWithRequestMetadata := injectFromHeaders(ctx, req)
		output, err := executeWorkflow(ctxWithRequestMetadata, temporalClient, wReq)
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

func getWorkflowRequest(req *http.Request) (wr workflow.CalculateParallelepipedWorkflowRequest, err error) {
	defer func() {
		if closeErr := req.Body.Close(); closeErr != nil {
			log.Print("error closing HTTP body: " + closeErr.Error())
		}
	}()
	err = json.NewDecoder(req.Body).Decode(&wr)
	return
}

func injectFromHeaders(ctx context.Context, req *http.Request) context.Context {
	jwt := extractTokenFromHeaders(req, temporal_microservices.AuthorizationHTTPHeader)
	processId := req.Header.Get(temporal_microservices.ProcessIDHTTPHeader)
	//nolint:staticcheck
	ctx1 := context.WithValue(ctx, temporal_microservices.ProcessIDContextField, processId)
	//nolint:staticcheck
	ctxOut := context.WithValue(ctx1, temporal_microservices.JWTContextField, jwt)
	return ctxOut
}

func extractTokenFromHeaders(req *http.Request, tokenHeaderName string) string {
	bearer := req.Header.Get(tokenHeaderName)
	authHeaderParts := strings.Split(bearer, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != bearerPrefix {
		return ""
	}
	return authHeaderParts[1]
}

func executeWorkflow(ctx context.Context, temporalClient client.Client, wReq workflow.CalculateParallelepipedWorkflowRequest) (output []workflow.Parallelepiped, err error) {
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: temporal_microservices.FigureWorkflowQueue,
	}
	workflowRun, err := temporalClient.ExecuteWorkflow(ctx, workflowOptions, workflow.CalculateParallelepipedWorkflow, wReq)
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

func writeError(rw http.ResponseWriter, err error) {
	log.Print(err.Error())
	rw.WriteHeader(http.StatusInternalServerError)
	if _, errWrite := rw.Write([]byte(err.Error())); errWrite != nil {
		log.Print("error writing the HTTP response: " + errWrite.Error())
	}
}
