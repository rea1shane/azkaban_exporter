package api

import (
	"azkaban_exporter/azkaban"
	"azkaban_exporter/util"
	"fmt"
	"net/http"
	"strings"
)

var singletonHttp = util.GetSingletonHttp()

// Login
// doc https://github.com/azkaban/azkaban/blob/master/docs/ajaxApi.rst#authenticate
// TODO 判断 session.id 是否过期, 没有过期的话跳过执行
func Login(url string, users []azkaban.User) []azkaban.User {
	method := "POST"
	for index, user := range users {
		response := LoginResponse{}
		payload := strings.NewReader("action=login&username=" + user.Username + "&password=" + user.Password)
		req, err := http.NewRequest(method, url, payload)
		if err != nil {
			panic(err)
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		singletonHttp.Request(req, &response)
		users[index].SessionId = response.Sessionid
	}
	return users
}

// GetProjects
// doc https://github.com/azkaban/azkaban/blob/master/docs/ajaxApi.rst#fetch-user-projects
func GetProjects(url string, users []azkaban.User) []azkaban.Project {
	method := "GET"
	var projects []azkaban.Project
	for _, user := range users {
		response := ProjectsResponse{}
		url := url + "/index?ajax=fetchuserprojects&session.id=" + user.SessionId
		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		singletonHttp.Request(req, &response)
		for _, project := range response.Projects {
			projects = append(projects, project)
		}
	}
	return projects
}

// GetFlows
// doc https://github.com/azkaban/azkaban/blob/master/docs/ajaxApi.rst#fetch-flows-of-a-project
func GetFlows(url string, projects []azkaban.Project) []azkaban.Flow {
	method := "GET"
	var flows []azkaban.Flow
	for _, project := range projects {
		response := ProjectFlowsResponse{}
		url := url + "/manager?ajax=fetchprojectflows&session.id=" + "3b311335-4716-45fb-9814-b2a4710297ac" + "&project=" + project.ProjectName
		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		singletonHttp.Request(req, &response)
		for _, flow := range response.Flows {
			flows = append(flows, flow)
		}
	}
	return flows
}
