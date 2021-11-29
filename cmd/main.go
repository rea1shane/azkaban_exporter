package main

import (
	"azkaban_exporter/azkaban"
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
)

func enter(exporter required.Exporter, target required.Target) {
	var (
		listenAddress = kingpin.Flag(
			"web.listen-address",
			"Address on which to expose metrics and web interface.",
		).Default(fmt.Sprintf(":%d", exporter.DefaultPort)).String()
		metricsPath = kingpin.Flag(
			"web.telemetry-path",
			"Path under which to expose metrics.",
		).Default("/metrics").String()
		configFile = kingpin.Flag(
			"web.config",
			"[EXPERIMENTAL] Path to config yaml file that can enable TLS or authentication.",
		).Default("").String()
	)

	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Version(version.Print(exporter.AppName))
	kingpin.CommandLine.UsageWriter(os.Stdout)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promlog.New(promlogConfig)

	_ = level.Info(logger).Log("msg", "Starting "+exporter.AppName, "version", version.Info())
	_ = level.Info(logger).Log("msg", "Build context", "build_context", version.BuildContext())
	if userCurrent, err := user.Current(); err == nil && userCurrent.Uid == "0" {
		_ = level.Warn(logger).Log("msg", exporter.TargetName+" Exporter is running as root user. This exporter is designed to run as unpriviledged user, root is not required.")
	}

	http.Handle(*metricsPath, prometheus.NewPrometheusHandler(logger, exporter, target))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
			<head><title>` + exporter.TargetName + ` Exporter</title></head>
			<body>
			<h1>` + exporter.TargetName + ` Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	server := &http.Server{Addr: *listenAddress}
	if err := web.ListenAndServe(server, *configFile, logger); err != nil {
		_ = level.Error(logger).Log("err", err)
		os.Exit(1)
	}
}

func main() {
	azkabanExporter := required.Exporter{
		AppName:     "azkaban_exporter",
		TargetName:  "Azkaban",
		DefaultPort: 9900,
	}
	a := azkaban.Azkaban{
		Namespace: "azkaban",
		Address: []string{
			"127.0.0.1:10000",
			"127.0.0.2:10000",
			"127.0.0.3:10000",
		},
	}
	enter(azkabanExporter, a)
}
