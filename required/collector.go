package required

import (
	prom2 "azkaban_exporter/pkg/prometheus"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	DefaultEnabled  = true
	DefaultDisabled = false
)

// Collector is the interface a collector has to implement.
type Collector interface {
	// Update Get new metrics and expose them via prometheus registry.
	Update(ch chan<- prometheus.Metric) error
}

// RegisterCollector After you implement the Collector, you should call this func to regist it.
func RegisterCollector(collector string, isDefaultEnabled bool, factory func(namespace string, logger log.Logger) (Collector, error)) {
	prom2.RegisterCollector(collector, isDefaultEnabled, factory)
}

// ErrNoData indicates the collector found no data to collect, but had no other error.
// When metric data is empty, return this error
func ErrNoData() error {
	return prom2.ErrNoData
}
