package azkaban

type LoginResponse struct {
	Sessionid string `json:"session.id"`
	Status    string `json:"status"`
	Error     string `json:"error"`
}
