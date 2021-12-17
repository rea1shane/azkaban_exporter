package functions

import (
	"azkaban_exporter/pkg/http"
	"azkaban_exporter/pkg/prometheus"
	"azkaban_exporter/required/structs"
	log "github.com/sirupsen/logrus"
)

func Run(logger *log.Logger, e structs.Exporter) {
	http.Run(logger, e)
}

// RegisterCollector After you implement the structs.Collector, you should call this func to regist it.
func RegisterCollector(collector string, isDefaultEnabled bool, factory func(namespace string, logger *log.Entry) (structs.Collector, error)) {
	prometheus.RegisterCollector(collector, isDefaultEnabled, factory)
}
