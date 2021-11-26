package require

import (
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type Target interface {
	NewCollector(logger log.Logger) (prometheus.Collector, error)
}
