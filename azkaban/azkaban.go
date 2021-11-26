package azkaban

import (
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type Azkaban struct {
	Address []string
}

func (a Azkaban) NewCollector(logger log.Logger) (prometheus.Collector, error) {
	return Collector{}, nil
}
