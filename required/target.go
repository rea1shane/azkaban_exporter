package required

import (
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type Target interface {
	// NewCollector create a new prometheus.Collector
	NewCollector(logger log.Logger) (prometheus.Collector, error)
}
