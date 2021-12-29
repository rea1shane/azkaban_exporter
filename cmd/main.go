package main

import (
	"github.com/rea1shane/azkaban_exporter/azkaban"
	"github.com/rea1shane/azkaban_exporter/required/functions"
	"github.com/rea1shane/azkaban_exporter/required/structs"
	"github.com/rea1shane/azkaban_exporter/util"
)

func main() {
	logger := util.GetLogger()
	azkabanExporter := structs.Exporter{
		MetricNamespace: "azkaban",
		ExporterName:    "azkaban_exporter",
		DefaultPort:     9900,
	}
	functions.Start(logger, azkabanExporter, azkaban.ParseArgs(azkabanExporter))
}
