package api

import (
	"azkaban_exporter/azkaban"
	"azkaban_exporter/util"
	"fmt"
	"net/http"
	"strings"
)

var singletonHttp = util.GetSingletonHttp()

// Login return azkaban.Session's Id
// doc https://github.com/azkaban/azkaban/blob/master/docs/ajaxApi.rst#authenticate
// TODO 传入一个 time.Time 检测 session.id 是否过期, 没有过期的话跳过执行
// TODO 返回一个 time.Time 代表登录时间
func Login(serverUrl string, user azkaban.User) (string, error) {
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

// GetProjects
// doc https://github.com/azkaban/azkaban/blob/master/docs/ajaxApi.rst#fetch-user-projects
func GetProjects(serverUrl string, sessionId string) ([]azkaban.Project, error) {
	method := "GET"
	response := ProjectsResponse{}
	url := serverUrl + "/index?ajax=fetchuserprojects&session.id=" + sessionId
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	err = singletonHttp.Request(req, &response)
	if err != nil {
		return nil, err
	}
	return response.Projects, nil
}

// GetFlows
// doc https://github.com/azkaban/azkaban/blob/master/docs/ajaxApi.rst#fetch-flows-of-a-project
func GetFlows(serverUrl string, sessionId string, projects []azkaban.Project) ([]azkaban.Flow, error) {
	method := "GET"
	var flows []azkaban.Flow
	for _, project := range projects {
		response := ProjectFlowsResponse{}
		url := serverUrl + "/manager?ajax=fetchprojectflows&session.id=" + sessionId + "&project=" + project.ProjectName
		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		err = singletonHttp.Request(req, &response)
		if err != nil {
			return nil, err
		}
		for _, flow := range response.Flows {
			flows = append(flows, flow)
		}
	}
	return flows, nil
}
