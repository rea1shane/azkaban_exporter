package api

import "azkaban_exporter/azkaban"

type LoginResponse struct {
	Sessionid string `json:"session.id"`
	Status    string `json:"status"`
	Error     string `json:"error"`
}

type ProjectsResponse struct {
	Projects []azkaban.Project `json:"projects"`
	Error    string            `json:"error"`
}

type ProjectFlowsResponse struct {
	ProjectName string         `json:"project"`
	ProjectId   int            `json:"projectId"`
	Flows       []azkaban.Flow `json:"flows"`
	Error       string         `json:"error"`
}
