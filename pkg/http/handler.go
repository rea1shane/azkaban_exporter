package http

import (
	"azkaban_exporter/require"
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

// PrometheusHandler is a common handler which implement net/http.Handler
type PrometheusHandler struct {
	Handler http.Handler
	Logger  log.Logger
}

func NewPrometheusHandler(logger log.Logger, exporter require.Exporter, target require.Target) *PrometheusHandler {
	h := &PrometheusHandler{
		Logger: logger,
	}
	if innerHandler, err := h.InnerHandler(exporter, target); err != nil {
		panic(fmt.Sprintf("Couldn't create metrics handler: %s", err))
	} else {
		h.Handler = innerHandler
	}
	return h
}

// ServeHTTP implements net/http.Handler.
func (h *PrometheusHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	h.Handler.ServeHTTP(writer, request)
}

func (h *PrometheusHandler) InnerHandler(exporter require.Exporter, target require.Target) (http.Handler, error) {
	collector, err := target.NewCollector(h.Logger)
	if err != nil {
		return nil, fmt.Errorf("couldn't create collector: %s", err)
	}

	r := prometheus.NewRegistry()
	r.MustRegister(version.NewCollector(exporter.AppName))
	if err := r.Register(collector); err != nil {
		return nil, fmt.Errorf("couldn't register %s collector: %s", strings.ToLower(exporter.TargetName), err)
	}
	handler := promhttp.HandlerFor(
		prometheus.Gatherers{r},
		promhttp.HandlerOpts{
			ErrorLog:            stdlog.New(log.NewStdlibAdapter(level.Error(h.Logger)), "", 0),
			ErrorHandling:       promhttp.ContinueOnError,
			MaxRequestsInFlight: 40, // TODO 变为读取命令行参数
		},
	)
	return handler, nil
}
