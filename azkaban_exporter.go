package main

import (
	"github.com/rea1shane/exporter"
	"github.com/rea1shane/gooooo/log"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	formatter := log.GetFormatter()
	formatter.FieldsOrder = []string{"StatusCode", "Latency", "Collector", "Duration"}
	logger.SetFormatter(formatter)
	exporter.Register("azkaban_exporter", "azkaban", "Exporter for <a href=\"https://github.com/azkaban/azkaban\">Azkaban</a> workflow manager.", ":9900", logger)
	exporter.Run()
}
