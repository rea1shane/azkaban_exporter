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
	Handler http.Handler
	Logger  log.Logger
}

func NewAzkabanHandler(logger log.Logger) (*AzkabanCollector, error) {
	panic("implement me")
}

func (h AzkabanHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	h.Handler.ServeHTTP(writer, request)
}

func (h *AzkabanHandler) InnerHandler() (http.Handler, error) {
	collector, err := NewAzkabanCollector(h.Logger)
	if err != nil {
		return nil, fmt.Errorf("couldn't create collector: %s", err)
	}

	r := prometheus.NewRegistry()
	r.MustRegister(version.NewCollector("node_exporter"))
	if err := r.Register(collector); err != nil {
		return nil, fmt.Errorf("couldn't register node collector: %s", err)
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
