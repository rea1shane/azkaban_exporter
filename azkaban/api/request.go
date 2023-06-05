package api

import (
	"context"
	"fmt"
	"github.com/morikuni/failure"
	"github.com/rea1shane/gooooo/data"
	myHttp "github.com/rea1shane/gooooo/http"
	"net/http"
	"net/url"
	"strings"
)

const (
	NewRequestErrorMessage failure.Message    = "new request"
	AzkabanError           failure.StringCode = "AzkabanError"
)

// Authenticate return a sessionId
// doc https://github.com/azkaban/azkaban/blob/master/docs/ajaxApi.rst#authenticate
func Authenticate(params AuthenticateParams, ctx context.Context) (string, error) {
	payload := strings.NewReader(fmt.Sprintf("action=login&username=%s&password=%s", params.Username, params.Password))
	req, err := http.NewRequest(http.MethodPost, params.ServerUrl, payload)
	if err != nil {
		return "", failure.Wrap(err, NewRequestErrorMessage)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	preprocess(req, ctx)

	response := Auth{}
	err = myHttp.Load(http.DefaultClient, req, &response, data.JsonFormat)
	if err != nil {
		return "", failure.Wrap(err, newRequestContext(req))
	}
	if response.Error != "" {
		return "", failure.New(AzkabanError, newAzkabanErrorContext(req, response.Error))
	}
	return response.SessionId, nil
}

// FetchUserProjects
// doc https://github.com/azkaban/azkaban/blob/master/docs/ajaxApi.rst#fetch-user-projects
func FetchUserProjects(params FetchUserProjectsParams, ctx context.Context) ([]Project, error) {
	u := fmt.Sprintf("%s/index?ajax=fetchuserprojects&session.id=%s", params.ServerUrl, params.SessionId)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, failure.Wrap(err, NewRequestErrorMessage)
	}
	preprocess(req, ctx)

	response := UserProjects{}
	err = myHttp.Load(http.DefaultClient, req, &response, data.JsonFormat)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, failure.New(AzkabanError, newAzkabanErrorContext(req, response.Error))
	}
	return response.Projects, nil
}

// FetchFlowsOfAProject
// doc https://github.com/azkaban/azkaban/blob/master/docs/ajaxApi.rst#fetch-flows-of-a-project
func FetchFlowsOfAProject(params FetchFlowsOfAProjectParams, ctx context.Context) ([]Flow, error) {
	u := fmt.Sprintf("%s/manager?ajax=fetchprojectflows&session.id=%s&project=%s", params.ServerUrl, params.SessionId, params.ProjectName)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, failure.Wrap(err, NewRequestErrorMessage)
	}
	preprocess(req, ctx)

	response := ProjectFlows{}
	err = myHttp.Load(http.DefaultClient, req, &response, data.JsonFormat)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, failure.New(AzkabanError, newAzkabanErrorContext(req, response.Error))
	}
	return response.Flows, nil
}

// FetchExecutionsOfAFlow
// doc https://github.com/azkaban/azkaban/blob/master/docs/ajaxApi.rst#fetch-executions-of-a-flow
func FetchExecutionsOfAFlow(params FetchExecutionsOfAFlowParams, ctx context.Context) (Executions, error) {
	u := fmt.Sprintf("%s/manager?ajax=fetchFlowExecutions&session.id=%s&project=%s&flow=%s&start=%d&length=%d", params.ServerUrl, params.SessionId, params.ProjectName, params.FlowId, params.StartIndex, params.ListLength)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return Executions{}, failure.Wrap(err, NewRequestErrorMessage)
	}
	preprocess(req, ctx)

	response := Executions{}
	err = myHttp.Load(http.DefaultClient, req, &response, data.JsonFormat)
	if err != nil {
		return Executions{}, err
	}
	if response.Error != "" {
		return Executions{}, failure.New(AzkabanError, newAzkabanErrorContext(req, response.Error))
	}
	return response, nil
}

/* --------------------------------------------- */

func preprocess(req *http.Request, ctx context.Context) {
	req = req.WithContext(ctx)
	req.URL.RawQuery = url.PathEscape(req.URL.RawQuery)
}

func newRequestContext(req *http.Request) failure.Context {
	return failure.Context{
		"url": req.URL.String(),
	}
}

func newAzkabanErrorContext(req *http.Request, errMsg string) failure.Context {
	requestContext := newRequestContext(req)
	requestContext["err_msg"] = errMsg
	return requestContext
}
