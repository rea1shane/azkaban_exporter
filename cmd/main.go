package main

import (
	"azkaban_exporter/pkg/azkaban"
	"azkaban_exporter/require"
	"fmt"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"os"
	"os/user"
)

func enter(t require.Target) {
	target := t.GetTargetName()
	appName := t.GetAppName()
	defaultListenAddress := fmt.Sprintf(":%d", t.GetDefaultListenPort())

	var (
		listenAddress = kingpin.Flag(
			"web.listen-address",
			"Address on which to expose metrics and web interface.",
		).Default(defaultListenAddress).String()
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
	kingpin.Version(version.Print(appName))
	kingpin.CommandLine.UsageWriter(os.Stdout)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promlog.New(promlogConfig)

	_ = level.Info(logger).Log("msg", "Starting "+appName, "version", version.Info())
	_ = level.Info(logger).Log("msg", "Build context", "build_context", version.BuildContext())
	if userCurrent, err := user.Current(); err == nil && userCurrent.Uid == "0" {
		_ = level.Warn(logger).Log("msg", target+" Exporter is running as root user. This exporter is designed to run as unpriviledged user, root is not required.")
	}

	// TODO 实现特定的 http Handler
	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
			<head><title>` + target + ` Exporter</title></head>
			<body>
			<h1>` + target + ` Exporter</h1>
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
	var target = azkaban.Azkaban{}
	enter(target)
}
