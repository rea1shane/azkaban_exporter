package prometheus

import (
	"azkaban_exporter/pkg/exporter"
	"encoding/json"
	"fmt"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	promcollectors "github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	stdlog "log"
	"net/http"
	"sort"
)

// Handler is a common handler which implement http.Handler
type Handler struct {
	UnfilteredHandler       http.Handler
	ExporterMetricsRegistry *prometheus.Registry // ExporterMetricsRegistry is a separate registry for the metrics about the exporter itself.
	IncludeExporterMetrics  bool
	MaxRequests             int
	Logger                  log.Logger
	Exporter                exporter.Exporter
}

func NewHandler(exporter exporter.Exporter, includeExporterMetrics bool, maxRequests int, logger log.Logger) *Handler {
	h := &Handler{
		ExporterMetricsRegistry: prometheus.NewRegistry(),
		IncludeExporterMetrics:  includeExporterMetrics,
		MaxRequests:             maxRequests,
		Logger:                  logger,
		Exporter:                exporter,
	}
	if h.IncludeExporterMetrics {
		h.ExporterMetricsRegistry.MustRegister(
			promcollectors.NewProcessCollector(promcollectors.ProcessCollectorOpts{}),
			promcollectors.NewGoCollector(),
		)
	}
	if innerHandler, err := h.InnerHandler(); err != nil {
		panic(fmt.Sprintf("Couldn't create metrics handler: %s", err))
	} else {
		h.UnfilteredHandler = innerHandler
	}
	return h
}

// ServeHTTP implements http.Handler.
func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	filters := request.URL.Query()["collect[]"]
	_ = level.Debug(h.Logger).Log("msg", "collect query:", "filters", func(filters []string) []byte {
		filterBytes, err := json.Marshal(filters)
		if err != nil {
			// TODO 处理 panic
			panic(fmt.Errorf(err.Error()))
		}
		return filterBytes
	}(filters))

	if len(filters) == 0 {
		// No filters, use the prepared unfiltered handler.
		h.UnfilteredHandler.ServeHTTP(writer, request)
		return
	}
	// To serve filtered metrics, we create a filtering handler on the fly.
	filteredHandler, err := h.InnerHandler(filters...)
	if err != nil {
		_ = level.Warn(h.Logger).Log("msg", "Couldn't create filtered metrics handler:", "err", err)
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte(fmt.Sprintf("Couldn't create filtered metrics handler: %s", err)))
		return
	}
	filteredHandler.ServeHTTP(writer, request)
}

// InnerHandler is used to create both the one unfiltered http.Handler to be
// wrapped by the outer Handler and also the filtered handlers created on the
// fly. The former is accomplished by calling InnerHandler without any arguments
// (in which case it will log all the collectors enabled via command-line flags).
func (h *Handler) InnerHandler(filters ...string) (http.Handler, error) {
	targetCollector, err := NewTargetCollector(h.Exporter, h.Logger, filters...)
	if err != nil {
		return nil, fmt.Errorf("couldn't create collector: %s", err)
	}

	// Only log the creation of an unfiltered handler, which should happen
	// only once upon startup.
	if len(filters) == 0 {
		_ = level.Info(h.Logger).Log("msg", "Enabled collectors")
		var collectors []string
		for n := range targetCollector.Collectors {
			collectors = append(collectors, n)
		}
		sort.Strings(collectors)
		for _, c := range collectors {
			_ = level.Info(h.Logger).Log("collector", c)
		}
	}

	r := prometheus.NewRegistry()
	r.MustRegister(version.NewCollector(h.Exporter.ExporterName))
	if err := r.Register(targetCollector); err != nil {
		return nil, fmt.Errorf("couldn't register "+h.Exporter.Namespace+" collector: %s", err)
	}
	handler := promhttp.HandlerFor(
		prometheus.Gatherers{h.ExporterMetricsRegistry, r},
		promhttp.HandlerOpts{
			ErrorLog:            stdlog.New(log.NewStdlibAdapter(level.Error(h.Logger)), "", 0),
			ErrorHandling:       promhttp.ContinueOnError,
			MaxRequestsInFlight: h.MaxRequests,
			Registry:            h.ExporterMetricsRegistry,
		},
	)
	if h.IncludeExporterMetrics {
		// Note that we have to use h.exporterMetricsRegistry here to use the same promhttp metrics for all expositions.
		handler = promhttp.InstrumentMetricHandler(
			h.ExporterMetricsRegistry, handler,
		)
	}
	return handler, nil
}
