package implement

import (
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type AzkabanCollector struct {
}

func NewAzkabanCollector(logger log.Logger) (*AzkabanCollector, error) {
	panic("implement me")
}

func (collector AzkabanCollector) Describe(descs chan<- *prometheus.Desc) {
	panic("implement me")
}

func (collector AzkabanCollector) Collect(metrics chan<- prometheus.Metric) {
	panic("implement me")
}
