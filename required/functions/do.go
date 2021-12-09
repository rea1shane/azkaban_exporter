package functions

import (
	"azkaban_exporter/pkg/prometheus"
	"azkaban_exporter/pkg/run"
	"azkaban_exporter/required/structs"
	"github.com/go-kit/log"
)

func Run(e structs.Exporter, errCh chan error) {
	run.Run(e, errCh)
}

// RegisterCollector After you implement the structs.Collector, you should call this func to regist it.
func RegisterCollector(collector string, isDefaultEnabled bool, factory func(namespace string, logger log.Logger) (structs.Collector, error)) {
	prometheus.RegisterCollector(collector, isDefaultEnabled, factory)
}

// ErrNoData indicates the collector found no data to collect, but had no other error.
// When metric data is empty, return this error
func ErrNoData() error {
	return prometheus.ErrNoData
}
