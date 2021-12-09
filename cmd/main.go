package main

import (
	_ "azkaban_exporter/azkaban"
	"azkaban_exporter/required/function"
	"azkaban_exporter/required/structs"
	"azkaban_exporter/util"
)

func main() {
	azkabanExporter := structs.Exporter{
		MonitorTargetName: "Azkaban",
		DefaultPort:       9900,
	}
	function.Run(azkabanExporter, util.GetErrorChannel())
}
