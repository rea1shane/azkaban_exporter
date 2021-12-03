package azkaban

type LoginResponse struct {
	Sessionid string `json:"session.id"`
	Status    string `json:"status"`
	Error     string `json:"error"`
}

type Project struct {
	ProjectId   int    `json:"projectId"`
	ProjectName string `json:"projectName"`
	CreatedBy   string `json:"createdBy"`
	//CreatedTime time.Time `json:"createdTime"`
	//userPermissions
	//groupPermissions
}

type GetProjectsResponse struct {
	Projects []Project `json:"projects"`
}
