package util

import (
	prom2 "azkaban_exporter/pkg/prometheus"
	"azkaban_exporter/required"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	DefaultEnabled  = true
	DefaultDisabled = false
)

// RegisterCollector After you implement the required.Collector, you should call this func to regist it.
func RegisterCollector(collector string, isDefaultEnabled bool, factory func(namespace string, logger log.Logger) (required.Collector, error)) {
	prom2.RegisterCollector(collector, isDefaultEnabled, factory)
}

// ErrNoData indicates the collector found no data to collect, but had no other error.
// When metric data is empty, return this error
func ErrNoData() error {
	return prom2.ErrNoData
}

// TypedDesc is suggestion metric's type
type TypedDesc struct {
	Desc      *prometheus.Desc
	ValueType prometheus.ValueType
}

func (t *TypedDesc) MustNewConstMetric(value float64, labels ...string) prometheus.Metric {
	return prometheus.MustNewConstMetric(t.Desc, t.ValueType, value, labels...)
}
