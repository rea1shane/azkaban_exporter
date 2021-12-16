package functions

import (
	"azkaban_exporter/pkg/prometheus"
	"azkaban_exporter/pkg/run"
	"azkaban_exporter/required/structs"
	"github.com/go-kit/log"
)

func Run(e structs.Exporter) {
	run.Run(e)
}

// RegisterCollector After you implement the structs.Collector, you should call this func to regist it.
func RegisterCollector(collector string, isDefaultEnabled bool, factory func(namespace string, logger log.Logger) (structs.Collector, error)) {
	prometheus.RegisterCollector(collector, isDefaultEnabled, factory)
}
