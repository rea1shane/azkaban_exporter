package util

import (
	"azkaban_exporter/pkg/prometheus"
	"azkaban_exporter/required"
	"github.com/go-kit/log"
)

func RegisterCollector(collector string, isDefaultEnabled bool, factory func(logger log.Logger) (required.Collector, error)) {
	prometheus.RegisterCollector(collector, isDefaultEnabled, factory)
}
