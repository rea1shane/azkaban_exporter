package api

type AuthenticateParams struct {
	ServerUrl string
	Username  string
	Password  string
}

type FetchUserProjectsParams struct {
	ServerUrl string
	SessionId string
}

type FetchFlowsOfAProjectParams struct {
	ServerUrl   string
	SessionId   string
	ProjectName string
}

type FetchExecutionsOfAFlowParams struct {
	ServerUrl   string
	SessionId   string
	ProjectName string
	FlowId      string
	StartIndex  int
	ListLength  int
}
