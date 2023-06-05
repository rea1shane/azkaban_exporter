package api

// TODO 不同版本的返回值
// TODO 返回值格式校验

type Project struct {
	ProjectId   int    `json:"projectId"`
	ProjectName string `json:"projectName"`
	CreatedBy   string `json:"createdBy"`
	//CreatedTime int64  `json:"createdTime"`
	//userPermissions
	//groupPermissions
}

type Flow struct {
	FlowId string `json:"flowId"`
}

type Auth struct {
	SessionId string `json:"session.id"`
	Status    string `json:"status"`
	Error     string `json:"error"`
}

type UserProjects struct {
	Projects []Project `json:"projects"`
	Error    string    `json:"error"`
}

type ProjectFlows struct {
	ProjectName string `json:"project"`
	ProjectId   int    `json:"projectId"`
	Flows       []Flow `json:"flows"`
	Error       string `json:"error"`
}

type Executions struct {
	Total      int         `json:"total"`
	Executions []Execution `json:"executions"`
	Length     int         `json:"length"`
	Project    string      `json:"project"`
	From       int         `json:"from"`
	ProjectID  int         `json:"projectId"`
	Flow       string      `json:"flow"`
	Error      string      `json:"error"`
}

type Execution struct {
	SubmitTime int64  `json:"submitTime"`
	SubmitUser string `json:"submitUser"`
	StartTime  int64  `json:"startTime"`
	EndTime    int64  `json:"endTime"`
	FlowID     string `json:"flowId"`
	ProjectID  int    `json:"projectId"`
	ExecID     int    `json:"execId"`
	Status     string `json:"status"`
}
