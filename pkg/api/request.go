package api

import (
	"context"
	"github.com/morikuni/failure"
	"github.com/rea1shane/basexporter/util"
	"net/http"
	"strconv"
	"strings"
)

var h = util.GetHttp(3, 500)

const (
	RequestError failure.StringCode = "request error"
)

// Authenticate return pkg.Session's SessionId
// doc https://github.com/azkaban/azkaban/blob/master/docs/ajaxApi.rst#authenticate
func Authenticate(p AuthenticateParam, ctx context.Context) (string, error) {
	method := "POST"
	response := Auth{}
	payload := strings.NewReader("action=login&username=" + p.Username + "&password=" + p.Password)
	req, err := http.NewRequest(method, p.ServerUrl, payload)
	if err != nil {
		return "", failure.Wrap(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	err = h.Request(req, ctx, &response)
	if err != nil {
		return "", err
	}
	if response.Error != "" {
		return "", failure.New(RequestError, failure.Context{
			"reason": response.Error,
		})
	}
	return response.SessionId, nil
}

// FetchUserProjects
// doc https://github.com/azkaban/azkaban/blob/master/docs/ajaxApi.rst#fetch-user-projects
func FetchUserProjects(p FetchUserProjectsParam, ctx context.Context) ([]Project, error) {
	method := "GET"
	response := UserProjects{}
	url := p.ServerUrl + "/index?ajax=fetchuserprojects&session.id=" + p.SessionId
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, failure.Wrap(err)
	}
	err = h.Request(req, ctx, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, failure.New(RequestError, failure.Context{
			"reason": response.Error,
		})
	}
	return response.Projects, nil
}

// FetchFlowsOfAProject
// doc https://github.com/azkaban/azkaban/blob/master/docs/ajaxApi.rst#fetch-flows-of-a-project
func FetchFlowsOfAProject(p FetchFlowsOfAProjectParam, ctx context.Context) ([]Flow, error) {
	method := "GET"
	response := ProjectFlows{}
	url := p.ServerUrl + "/manager?ajax=fetchprojectflows&session.id=" + p.SessionId + "&project=" + p.ProjectName
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, failure.Wrap(err)
	}
	err = h.Request(req, ctx, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, failure.New(RequestError, failure.Context{
			"reason": response.Error,
		})
	}
	return response.Flows, nil
}

// FetchExecutionsOfAFlow
// doc https://github.com/azkaban/azkaban/blob/master/docs/ajaxApi.rst#fetch-executions-of-a-flow
func FetchExecutionsOfAFlow(p FetchExecutionsOfAFlowParam, ctx context.Context) (Executions, error) {
	method := "GET"
	response := Executions{}
	url := p.ServerUrl + "/manager?ajax=fetchFlowExecutions&session.id=" + p.SessionId + "&project=" + p.ProjectName + "&flow=" + p.FlowId +
		"&start=" + strconv.Itoa(p.StartIndex) + "&length=" + strconv.Itoa(p.ListLength)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return Executions{}, failure.Wrap(err)
	}
	err = h.Request(req, ctx, &response)
	if err != nil {
		return Executions{}, err
	}
	if response.Error != "" {
		return Executions{}, failure.New(RequestError, failure.Context{
			"reason": response.Error,
		})
	}
	return response, nil
}
