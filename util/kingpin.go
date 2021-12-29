package util

import (
	"fmt"
	"github.com/prometheus/common/version"
	"github.com/rea1shane/azkaban_exporter/required/structs"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

func ParseArgs(e structs.Exporter) structs.Args {
	var (
		listenAddress = kingpin.Flag(
			"web.listen-address",
			"Address on which to expose metrics and web interface.",
		).Default(fmt.Sprintf(":%d", e.DefaultPort)).String()
		metricsPath = kingpin.Flag(
			"web.telemetry-path",
			"Path under which to expose metrics.",
		).Default("/metrics").String()
		disableExporterMetrics = kingpin.Flag(
			"web.disable-exporter-metrics",
			"Exclude metrics about the exporter itself (promhttp_*, process_*, go_*).",
		).Default("false").Bool()
		maxRequests = kingpin.Flag(
			"web.max-requests",
			"Maximum number of parallel scrape requests. Use 0 to disable.",
		).Default("40").Int()
	)

	kingpin.Version(version.Print(e.ExporterName))
	kingpin.CommandLine.UsageWriter(os.Stdout)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	return structs.Args{
		ListenAddress:          *listenAddress,
		MetricsPath:            *metricsPath,
		DisableExporterMetrics: *disableExporterMetrics,
		MaxRequests:            *maxRequests,
	}
}
