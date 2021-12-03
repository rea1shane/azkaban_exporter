package api

import "azkaban_exporter/azkaban"

type AuthenticateResponse struct {
	SessionId string `json:"session.id"`
	Status    string `json:"status"`
	Error     string `json:"error"`
}

type FetchUserProjectsResponse struct {
	Projects []azkaban.Project `json:"projects"`
	Error    string            `json:"error"`
}

type FetchFlowsOfAProjectResponse struct {
	ProjectName string         `json:"project"`
	ProjectId   int            `json:"projectId"`
	Flows       []azkaban.Flow `json:"flows"`
	Error       string         `json:"error"`
}
