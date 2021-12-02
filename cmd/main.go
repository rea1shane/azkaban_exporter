package main

import (
	"azkaban_exporter/azkaban"
	exporterinfo "azkaban_exporter/pkg/exporter"
	"azkaban_exporter/pkg/prometheus"
	"azkaban_exporter/required"
	"fmt"
	"github.com/go-kit/log/level"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"os"
	"os/user"
	"strings"
)

func enter(exporter required.Exporter) {
	exporterInfo := exporterinfo.Exporter{
		Namespace:    strings.ToLower(exporter.MonitorTargetName),
		ExporterName: strings.ToLower(exporter.MonitorTargetName) + "_exporter",
	}
	var (
		listenAddress = kingpin.Flag(
			"web.listen-address",
			"Address on which to expose metrics and web interface.",
		).Default(fmt.Sprintf(":%d", exporter.DefaultPort)).String()
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
		configFile = kingpin.Flag(
			"web.config",
			"[EXPERIMENTAL] Path to config yaml file that can enable TLS or authentication.",
		).Default("").String()
	)

	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Version(version.Print(exporterInfo.ExporterName))
	kingpin.CommandLine.UsageWriter(os.Stdout)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promlog.New(promlogConfig)

	_ = level.Info(logger).Log("msg", "Starting "+exporterInfo.ExporterName, "version", version.Info())
	_ = level.Info(logger).Log("msg", "Build context", "build_context", version.BuildContext())
	if userCurrent, err := user.Current(); err == nil && userCurrent.Uid == "0" {
		_ = level.Warn(logger).Log("msg", exporter.MonitorTargetName+" Exporter is running as root user. This exporter is designed to run as unpriviledged user, root is not required.")
	}

	http.Handle(*metricsPath, prometheus.NewHandler(exporterInfo, !*disableExporterMetrics, *maxRequests, logger))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
			<head><title>` + exporter.MonitorTargetName + ` Exporter</title></head>
			<body>
			<h1>` + exporter.MonitorTargetName + ` Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	_ = level.Info(logger).Log("msg", "Listening on", "address", *listenAddress)
	server := &http.Server{Addr: *listenAddress}
	if err := web.ListenAndServe(server, *configFile, logger); err != nil {
		_ = level.Error(logger).Log("err", err)
		os.Exit(1)
	}
}

//func main() {
//	azkabanExporter := required.Exporter{
//		MonitorTargetName: "Azkaban",
//		DefaultPort:       9900,
//	}
//	enter(azkabanExporter)
//}

func main() {
	az := azkaban.GetInstance()
	fmt.Printf("%+v", az)
}