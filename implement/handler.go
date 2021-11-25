package implement

import (
	"fmt"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	stdlog "log"
	"net/http"
)

// AzkabanHandler wraps an unfiltered http.Handler
type AzkabanHandler struct {
	Logger log.Logger
}

func (h AzkabanHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	collector, err := NewAzkabanCollector(h.Logger)
	if err != nil {
		return
	}

	r := prometheus.NewRegistry()
	r.MustRegister(version.NewCollector("node_exporter"))
	if err := r.Register(collector); err != nil {
		_ = fmt.Errorf("couldn't register node collector: %s", err)
		return
	}
	handler := promhttp.HandlerFor(
		prometheus.Gatherers{r},
		promhttp.HandlerOpts{
			ErrorLog:            stdlog.New(log.NewStdlibAdapter(level.Error(h.Logger)), "", 0),
			ErrorHandling:       promhttp.ContinueOnError,
			MaxRequestsInFlight: 40, // TODO 变为读取命令行参数
		},
	)
	handler.ServeHTTP(writer, request)
}
