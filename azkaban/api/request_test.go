package api

import (
	"context"
	"fmt"
	"os"
	"testing"
)

const (
	serverUrl = "http://azkaban:20000"
	username  = "metrics"
	password  = "metrics"
)

func TestAPI(t *testing.T) {
	// Authenticate
	sessionId, err := Authenticate(AuthenticateParams{
		ServerUrl: serverUrl,
		Username:  username,
		Password:  password,
	}, context.Background())
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("SESSION_ID: %s\n", sessionId)

	// FetchUserProjects
	projects, err := FetchUserProjects(FetchUserProjectsParams{
		ServerUrl: serverUrl,
		SessionId: sessionId,
	}, context.Background())
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", projects)

	if len(projects) == 0 {
		fmt.Println("No project")
		os.Exit(1)
	}

	// FetchFlowsOfAProject
	flows, err := FetchFlowsOfAProject(FetchFlowsOfAProjectParams{
		ServerUrl:   serverUrl,
		SessionId:   sessionId,
		ProjectName: projects[0].ProjectName, // or specify manually
	}, context.Background())
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", flows)

	if len(flows) == 0 {
		fmt.Println("No flow")
		os.Exit(1)
	}

	// FetchExecutionsOfAFlow
	executions, err := FetchExecutionsOfAFlow(FetchExecutionsOfAFlowParams{
		ServerUrl:   serverUrl,
		SessionId:   sessionId,
		ProjectName: projects[0].ProjectName, // or specify manually
		FlowId:      flows[0].FlowId,         // or specify manually
		StartIndex:  0,
		ListLength:  1,
	}, context.Background())
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", executions)
}
