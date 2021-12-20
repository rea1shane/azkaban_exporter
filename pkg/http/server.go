package http

import (
	"azkaban_exporter/pkg/args"
	"azkaban_exporter/pkg/exporter"
	"azkaban_exporter/pkg/middleware"
	"azkaban_exporter/pkg/prometheus"
	"azkaban_exporter/required/structs"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/version"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var srv *http.Server

func Start(logger *log.Logger, e structs.Exporter) {
	exporterInfo := exporter.Exporter{
		Namespace:    strings.ToLower(e.MonitorTargetName),
		ExporterName: strings.ToLower(e.MonitorTargetName) + "_exporter",
		DefaultPort:  e.DefaultPort,
	}

	a := args.ParseArgs(exporterInfo)

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
			<p><a href="` + *a.MetricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	}))
	app.GET(*a.MetricsPath, gin.WrapH(prometheus.NewHandler(exporterInfo, !*a.DisableExporterMetrics, *a.MaxRequests, logger)))

	logger.Info("Listening on address ", *a.ListenAddress)
	srv = &http.Server{
		Addr:    *a.ListenAddress,
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
