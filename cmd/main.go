package main

import (
	_ "azkaban_exporter/azkaban"
	"azkaban_exporter/required/functions"
	"azkaban_exporter/required/structs"
)

func main() {
	azkabanExporter := structs.Exporter{
		MonitorTargetName: "Azkaban",
		DefaultPort:       9900,
	}
	functions.Run(azkabanExporter)
}
