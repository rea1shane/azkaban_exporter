package api

type AuthenticateParam struct {
	ServerUrl string
	Username  string
	Password  string
}

type FetchUserProjectsParam struct {
	ServerUrl string
	SessionId string
}

type FetchFlowsOfAProjectParam struct {
	ServerUrl   string
	SessionId   string
	ProjectName string
}

type FetchExecutionsOfAFlowParam struct {
	ServerUrl   string
	SessionId   string
	ProjectName string
	FlowId      string
	StartIndex  int
	ListLength  int
}
