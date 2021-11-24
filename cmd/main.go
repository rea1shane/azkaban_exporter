package main

import (
	"azkaban-exporter/config"
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

const target = "Azkaban"

func main() {
	var port int
	var configFilePath string
	flag.IntVar(&port, "web.listen-address", 3001, `Server port, default: 3001`)
	flag.StringVar(&configFilePath, "config.file", "conf/exporter.yml", fmt.Sprintf("%s exporter configuration file path.", target))
	flag.Parse()

	var e config.Exporter
	e.GetConf(configFilePath)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", e.Port), nil))
}
