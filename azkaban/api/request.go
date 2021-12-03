package api

import (
	"azkaban_exporter/azkaban"
	"azkaban_exporter/util"
	"net/http"
	"strings"
)

var singletonHttp = util.GetSingletonHttp()

// Authenticate return azkaban.Session's Id
// doc https://github.com/azkaban/azkaban/blob/master/docs/ajaxApi.rst#authenticate
// TODO 传入一个 time.Time 检测 session.id 是否过期, 没有过期的话跳过执行
// TODO 返回一个 time.Time 代表登录时间
func Authenticate(serverUrl string, user azkaban.User) (string, error) {
	method := "POST"
	response := LoginResponse{}
	payload := strings.NewReader("action=login&username=" + user.Username + "&password=" + user.Password)
	req, err := http.NewRequest(method, serverUrl, payload)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	err = singletonHttp.Request(req, &response)
	if err != nil {
		return "", err
	}
	return response.SessionId, nil
}

// FetchUserProjects
// doc https://github.com/azkaban/azkaban/blob/master/docs/ajaxApi.rst#fetch-user-projects
func FetchUserProjects(serverUrl string, sessionId string) ([]azkaban.Project, error) {
	method := "GET"
	response := ProjectsResponse{}
	url := serverUrl + "/index?ajax=fetchuserprojects&session.id=" + sessionId
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	err = singletonHttp.Request(req, &response)
	if err != nil {
		return nil, err
	}
	return response.Projects, nil
}

// FetchFlowsOfAProject
// doc https://github.com/azkaban/azkaban/blob/master/docs/ajaxApi.rst#fetch-flows-of-a-project
func FetchFlowsOfAProject(serverUrl string, sessionId string, project azkaban.Project) ([]azkaban.Flow, error) {
	method := "GET"
	response := ProjectFlowsResponse{}
	url := serverUrl + "/manager?ajax=fetchprojectflows&session.id=" + sessionId + "&project=" + project.ProjectName
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	err = singletonHttp.Request(req, &response)
	if err != nil {
		return nil, err
	}
	return response.Flows, nil
}
