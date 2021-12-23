package main

import (
	"azkaban_exporter/azkaban"
	"azkaban_exporter/required/functions"
	"azkaban_exporter/required/structs"
	"azkaban_exporter/util"
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
