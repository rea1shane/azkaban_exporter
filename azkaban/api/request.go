package api

import (
	"azkaban_exporter/util"
	"net/http"
	"strconv"
	"strings"
)

var singletonHttp = util.GetSingletonHttp()

// Authenticate return azkaban.Session's SessionId
// doc https://github.com/azkaban/azkaban/blob/master/docs/ajaxApi.rst#authenticate
func Authenticate(serverUrl string, username string, password string) (string, error) {
	method := "POST"
	response := Auth{}
	payload := strings.NewReader("action=login&username=" + username + "&password=" + password)
	req, err := http.NewRequest(method, serverUrl, payload)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	err = singletonHttp.Request(req, &response)
	if err != nil {
		return "", err
	}
	if response.Error != "" {
		return "", util.RequestFailureError("authenticate", response.Error)
	}
	return response.SessionId, nil
}

// FetchUserProjects
// doc https://github.com/azkaban/azkaban/blob/master/docs/ajaxApi.rst#fetch-user-projects
func FetchUserProjects(serverUrl string, sessionId string) ([]Project, error) {
	method := "GET"
	response := UserProjects{}
	url := serverUrl + "/index?ajax=fetchuserprojects&session.id=" + sessionId
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	err = singletonHttp.Request(req, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, util.RequestFailureError("fetch-user-projects", response.Error)
	}
	return response.Projects, nil
}

// FetchFlowsOfAProject
// doc https://github.com/azkaban/azkaban/blob/master/docs/ajaxApi.rst#fetch-flows-of-a-project
func FetchFlowsOfAProject(serverUrl string, sessionId string, projectName string) ([]Flow, error) {
	method := "GET"
	response := ProjectFlows{}
	url := serverUrl + "/manager?ajax=fetchprojectflows&session.id=" + sessionId + "&project=" + projectName
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	err = singletonHttp.Request(req, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, util.RequestFailureError("fetch-flows-of-a-project", response.Error)
	}
	return response.Flows, nil
}

// FetchExecutionsOfAFlow
// doc https://github.com/azkaban/azkaban/blob/master/docs/ajaxApi.rst#fetch-executions-of-a-flow
func FetchExecutionsOfAFlow(serverUrl string, sessionId string, projectName string, flowId string, startIndex int, listLength int) (Executions, error) {
	method := "GET"
	response := Executions{}
	url := serverUrl + "/manager?ajax=fetchFlowExecutions&session.id=" + sessionId + "&project=" + projectName + "&flow=" + flowId +
		"&start=" + strconv.Itoa(startIndex) + "&length=" + strconv.Itoa(listLength)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return Executions{}, err
	}
	err = singletonHttp.Request(req, &response)
	if err != nil {
		return Executions{}, err
	}
	if response.Error != "" {
		return Executions{}, util.RequestFailureError("fetch-running-executions-of-a-flow", response.Error)
	}
	return response, nil
}
