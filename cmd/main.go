package main

import (
	"github.com/rea1shane/azkaban_exporter/pkg"
	"github.com/rea1shane/basexporter"
	"github.com/rea1shane/basexporter/util"
)

func main() {
	logger := util.GetLogger()
	defaultPort := 9900
	azkabanExporter := basexporter.BuildExporter(
		"azkaban",
		"azkaban_exporter",
		defaultPort,
		"v1.1.1")
	basexporter.Start(logger, azkabanExporter, pkg.ParseArgs(defaultPort))
}
