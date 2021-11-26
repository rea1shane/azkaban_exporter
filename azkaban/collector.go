package azkaban

import (
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {
	Logger log.Logger
}

func (collector Collector) Describe(descs chan<- *prometheus.Desc) {
	panic("implement me")
}

func (collector Collector) Collect(metrics chan<- prometheus.Metric) {
	panic("implement me")
}
