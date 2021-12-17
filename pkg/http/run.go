package http

import (
	"azkaban_exporter/pkg/exporter"
	"azkaban_exporter/pkg/middleware"
	"azkaban_exporter/pkg/prometheus"
	"azkaban_exporter/required/structs"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/version"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var srv *http.Server

func Run(logger *log.Logger, e structs.Exporter) {
	exporterInfo := exporter.Exporter{
		Namespace:    strings.ToLower(e.MonitorTargetName),
		ExporterName: strings.ToLower(e.MonitorTargetName) + "_exporter",
	}
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

	kingpin.Version(version.Print(exporterInfo.ExporterName))
	kingpin.CommandLine.UsageWriter(os.Stdout)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	logger.Info("Starting "+exporterInfo.ExporterName, "version", version.Info())
	logger.Info("Build context", version.BuildContext())

	app := gin.New()
	app.Use(
		middleware.ToStdout(logger),
		gin.Recovery(),
	)
	app.GET("/", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
			<head><title>` + e.MonitorTargetName + ` Exporter</title></head>
			<body>
			<h1>` + e.MonitorTargetName + ` Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	}))
	app.GET(*metricsPath, gin.WrapH(prometheus.NewHandler(exporterInfo, !*disableExporterMetrics, *maxRequests, logger)))

	logger.Info("Listening on address ", *listenAddress)
	srv = &http.Server{
		Addr:    *listenAddress,
		Handler: app,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	shutdown(logger)
}

func shutdown(logger *log.Logger) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGQUIT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.WithError(err).Error("server shutdown")
		return
	}
}
