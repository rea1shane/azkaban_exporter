package main

import (
	_ "azkaban_exporter/azkaban"
	"azkaban_exporter/required/functions"
	"azkaban_exporter/required/structs"
	"azkaban_exporter/util"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	logger := util.GetLogger()
	logger.SetLevel(log.InfoLevel)
	gin.SetMode(gin.ReleaseMode)
	azkabanExporter := structs.Exporter{
		MonitorTargetName: "Azkaban",
		DefaultPort:       9900,
	}
	functions.Start(logger, azkabanExporter)
}
