package http

import (
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
	"syscall"
	"time"
)

var srv *http.Server

func Start(logger *log.Logger, e structs.Exporter, args structs.Args) {
	displayName := camelString(e.ExporterName)

	logger.Info("Starting "+e.ExporterName, " version", version.Info())
	logger.Info("Build context", version.BuildContext())

	app := gin.New()
	app.Use(
		middleware.ToStdout(logger),
		gin.Recovery(),
	)
	app.GET("/", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
			<head><title>` + displayName + `</title></head>
			<body>
			<h1>` + displayName + `</h1>
			<p><a href="` + args.MetricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	}))
	app.GET(args.MetricsPath, gin.WrapH(prometheus.NewHandler(e, !args.DisableExporterMetrics, args.MaxRequests, logger)))

	logger.Info("Listening on address ", args.ListenAddress)
	srv = &http.Server{
		Addr:    args.ListenAddress,
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

func camelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			data = append(data, ' ')
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}
