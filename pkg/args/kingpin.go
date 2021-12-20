package args

import (
	"azkaban_exporter/pkg/exporter"
	"fmt"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

func ParseArgs(e exporter.Exporter) Args {
	var args = Args{
		ListenAddress: kingpin.Flag(
			"web.listen-address",
			"Address on which to expose metrics and web interface.",
		).Default(fmt.Sprintf(":%d", e.DefaultPort)).String(),
		MetricsPath: kingpin.Flag(
			"web.telemetry-path",
			"Path under which to expose metrics.",
		).Default("/metrics").String(),
		DisableExporterMetrics: kingpin.Flag(
			"web.disable-exporter-metrics",
			"Exclude metrics about the exporter itself (promhttp_*, process_*, go_*).",
		).Default("false").Bool(),
		MaxRequests: kingpin.Flag(
			"web.max-requests",
			"Maximum number of parallel scrape requests. Use 0 to disable.",
		).Default("40").Int(),
	}

	kingpin.Version(version.Print(e.ExporterName))
	kingpin.CommandLine.UsageWriter(os.Stdout)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	return args
}
