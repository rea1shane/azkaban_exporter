package server

type Server interface {
	CheckHealth() bool
	GetAllMetrics() string
	Convert() string
}
