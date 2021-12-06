package azkaban

import (
	"azkaban_exporter/required"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

var (

)

type azkabanCollector struct {
	namespace string
	logger    log.Logger
}

func NewAzkabanCollector(namespace string, logger log.Logger) (required.Collector, error) {
	return &azkabanCollector{
		namespace: namespace,
		logger:    logger,
	}, nil
}

func (c azkabanCollector) Update(ch chan<- prometheus.Metric) error {
	panic("implement me")
}

