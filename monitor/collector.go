package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
)

type AzkabanCollector struct {
}

func (collector AzkabanCollector) Describe(descs chan<- *prometheus.Desc) {
	panic("implement me")
}

func (collector AzkabanCollector) Collect(metrics chan<- prometheus.Metric) {
	panic("implement me")
}
