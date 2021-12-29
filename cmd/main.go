package main

import (
	"github.com/rea1shane/azkaban_exporter/pkg"
	"github.com/rea1shane/basexporter/required/functions"
	"github.com/rea1shane/basexporter/required/structs"
	"github.com/rea1shane/basexporter/util"
)

func main() {
	logger := util.GetLogger()
	azkabanExporter := structs.Exporter{
		MetricNamespace: "azkaban",
		ExporterName:    "azkaban_exporter",
		DefaultPort:     9900,
	}
	functions.Start(logger, azkabanExporter, pkg.ParseArgs(azkabanExporter))
}
