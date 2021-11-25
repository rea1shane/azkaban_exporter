package require

import (
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
)

type Target interface {
	NewCollector(logger log.Logger) (*prometheus.Collector, error)
	NewHandler(logger log.Logger) *http.Handler
}
