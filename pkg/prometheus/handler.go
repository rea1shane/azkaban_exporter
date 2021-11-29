package prometheus

import (
	"azkaban_exporter/required"
	"fmt"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	stdlog "log"
	"net/http"
	"strings"
)

// Handler is a common handler which implement http.Handler
type Handler struct {
	Handler http.Handler
	Logger  log.Logger
}

func NewPrometheusHandler(logger log.Logger, exporter required.Exporter, target required.Target) *Handler {
	h := &Handler{
		Logger: logger,
	}
	if innerHandler, err := h.InnerHandler(exporter, target); err != nil {
		panic(fmt.Sprintf("Couldn't create metrics handler: %s", err))
	} else {
		h.Handler = innerHandler
	}
	return h
}

// ServeHTTP implements http.Handler.
func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	h.Handler.ServeHTTP(writer, request)
}

// InnerHandler create a http.Handler in Handler
func (h *Handler) InnerHandler(exporter required.Exporter, target required.Target) (http.Handler, error) {
	collector, err := NewTargetCollector(target.GetNamespace(), exporter.AppName, h.Logger)
	if err != nil {
		return nil, fmt.Errorf("couldn't create collector: %s", err)
	}

	registry := prometheus.NewRegistry()
	registry.MustRegister(version.NewCollector(exporter.AppName))
	if err := registry.Register(collector); err != nil {
		return nil, fmt.Errorf("couldn't register %s collector: %s", strings.ToLower(exporter.TargetName), err)
	}
	handler := promhttp.HandlerFor(
		prometheus.Gatherers{registry},
		promhttp.HandlerOpts{
			ErrorLog:            stdlog.New(log.NewStdlibAdapter(level.Error(h.Logger)), "", 0),
			ErrorHandling:       promhttp.ContinueOnError,
			MaxRequestsInFlight: 40, // TODO 变为读取命令行参数
		},
	)
	return handler, nil
}
