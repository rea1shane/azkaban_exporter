package prometheus

import (
	"azkaban_exporter/required"
	"errors"
	"fmt"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/alecthomas/kingpin.v2"
	"sync"
	"time"
)

var (
	Factories              = make(map[string]func(logger log.Logger) (required.Collector, error))
	InitiatedCollectorsMtx = sync.Mutex{}
	InitiatedCollectors    = make(map[string]required.Collector)
	CollectorState         = make(map[string]*bool)
	ForcedCollectors       = map[string]bool{} // ForcedCollectors collectors which have been explicitly enabled or disabled
)

type TargetCollector struct {
	Collectors         map[string]required.Collector
	logger             log.Logger
	ScrapeDurationDesc *prometheus.Desc
	ScrapeSuccessDesc  *prometheus.Desc
}

func (t TargetCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- t.ScrapeDurationDesc
	ch <- t.ScrapeSuccessDesc
}

func (t TargetCollector) Collect(ch chan<- prometheus.Metric) {
	wg := sync.WaitGroup{}
	wg.Add(len(t.Collectors))
	for name, c := range t.Collectors {
		go func(name string, c required.Collector) {
			Execute(name, c, ch, t.logger, t.ScrapeDurationDesc, t.ScrapeSuccessDesc)
			wg.Done()
		}(name, c)
	}
	wg.Wait()
}

// NewTargetCollector creates a new TargetCollector.
func NewTargetCollector(namespace string, appName string, logger log.Logger) (*TargetCollector, error) {
	f := make(map[string]bool)
	collectors := make(map[string]required.Collector)
	InitiatedCollectorsMtx.Lock()
	defer InitiatedCollectorsMtx.Unlock()
	for key, enabled := range CollectorState {
		if !*enabled || (len(f) > 0 && !f[key]) {
			continue
		}
		if collector, ok := InitiatedCollectors[key]; ok {
			collectors[key] = collector
		} else {
			collector, err := Factories[key](log.With(logger, "collector", key))
			if err != nil {
				return nil, err
			}
			collectors[key] = collector
			InitiatedCollectors[key] = collector
		}
	}
	return &TargetCollector{
		Collectors: InitiatedCollectors,
		logger:     logger,
		ScrapeDurationDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "scrape", "collector_duration_seconds"),
			appName+": Duration of a collector scrape.",
			[]string{"collector"},
			nil,
		),
		ScrapeSuccessDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "scrape", "collector_success"),
			appName+": Whether a collector succeeded.",
			[]string{"collector"},
			nil,
		),
	}, nil
}

func RegisterCollector(collector string, isDefaultEnabled bool, factory func(logger log.Logger) (required.Collector, error)) {
	var helpDefaultState string
	if isDefaultEnabled {
		helpDefaultState = "enabled"
	} else {
		helpDefaultState = "disabled"
	}

	flagName := fmt.Sprintf("collector.%s", collector)
	flagHelp := fmt.Sprintf("Enable the %s collector (default: %s).", collector, helpDefaultState)
	defaultValue := fmt.Sprintf("%v", isDefaultEnabled)

	flag := kingpin.Flag(flagName, flagHelp).Default(defaultValue).Action(CollectorFlagAction(collector)).Bool()
	CollectorState[collector] = flag

	Factories[collector] = factory
}

func CollectorFlagAction(collector string) func(ctx *kingpin.ParseContext) error {
	return func(ctx *kingpin.ParseContext) error {
		ForcedCollectors[collector] = true
		return nil
	}
}

func Execute(name string, c required.Collector, ch chan<- prometheus.Metric, logger log.Logger, scrapeDurationDesc *prometheus.Desc, scrapeSuccessDesc *prometheus.Desc) {
	begin := time.Now()
	err := c.Update(ch)
	duration := time.Since(begin)
	var success float64

	if err != nil {
		if IsNoDataError(err) {
			_ = level.Debug(logger).Log("msg", "collector returned no data", "name", name, "duration_seconds", duration.Seconds(), "err", err)
		} else {
			_ = level.Error(logger).Log("msg", "collector failed", "name", name, "duration_seconds", duration.Seconds(), "err", err)
		}
		success = 0
	} else {
		_ = level.Debug(logger).Log("msg", "collector succeeded", "name", name, "duration_seconds", duration.Seconds())
		success = 1
	}
	ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, duration.Seconds(), name)
	ch <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, success, name)
}

// ErrNoData indicates the collector found no data to collect, but had no other error.
var ErrNoData = errors.New("collector returned no data")

func IsNoDataError(err error) bool {
	return err == ErrNoData
}
