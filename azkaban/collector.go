package azkaban

import (
	"azkaban_exporter/required"
	"azkaban_exporter/util"
	"encoding/json"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

const (
	subsystem = "execution"
)

var (
	azkaban = GetAzkaban()
)

func init() {
	util.RegisterCollector(subsystem, util.DefaultEnabled, NewAzkabanCollector)
}

type azkabanCollector struct {
	logger      log.Logger
	preparing   util.TypedDesc
	running     util.TypedDesc
	running0    util.TypedDesc
	running60   util.TypedDesc
	running300  util.TypedDesc
	running1440 util.TypedDesc
}

func NewAzkabanCollector(namespace string, logger log.Logger) (required.Collector, error) {
	var (
		labelNames = []string{"project"}
	)

	return &azkabanCollector{
		logger: logger,
		preparing: util.TypedDesc{
			Desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "preparing"),
				"The number of preparing executions.", labelNames, nil),
			ValueType: prometheus.GaugeValue,
		},
		running: util.TypedDesc{
			Desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "running"),
				"The number of running executions.", labelNames, nil),
			ValueType: prometheus.GaugeValue,
		},
		running0: util.TypedDesc{
			Desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "running_0"),
				"The number of running executions which running time in [0, 60) mins.", labelNames, nil),
			ValueType: prometheus.GaugeValue,
		},
		running60: util.TypedDesc{
			Desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "running_60"),
				"The number of running executions which running time in [60, 300) mins.", labelNames, nil),
			ValueType: prometheus.GaugeValue,
		},
		running300: util.TypedDesc{
			Desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "running_300"),
				"The number of running executions which running time in [300, 1440) mins.", labelNames, nil),
			ValueType: prometheus.GaugeValue,
		},
		running1440: util.TypedDesc{
			Desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "running_1440"),
				"The number of running executions which running time over 1440 mins.", labelNames, nil),
			ValueType: prometheus.GaugeValue,
		},
	}, nil
}

func (c azkabanCollector) Update(ch chan<- prometheus.Metric) error {
	projectNames, err := azkaban.GetProjectNames()
	if err != nil {
		return err
	}
	var (
		preparingCounter   = map[string]int{}
		runningCounter     = map[string]int{}
		running0Counter    = map[string]int{}
		running60Counter   = map[string]int{}
		running300Counter  = map[string]int{}
		running1440Counter = map[string]int{}
	)
	for _, projectName := range projectNames {
		preparingCounter[projectName] = 0
		runningCounter[projectName] = 0
		running0Counter[projectName] = 0
		running60Counter[projectName] = 0
		running300Counter[projectName] = 0
		running1440Counter[projectName] = 0
	}
	ids, err := azkaban.GetRunningExecIds()
	if err != nil {
		return err
	}
	_ = level.Debug(c.logger).Log("msg", "running exec ids", "ids", func(ids []int) []byte {
		idsString, _ := json.Marshal(ids)
		return idsString
	}(ids))
	infos, err := azkaban.GetExecInfos(ids)
	if err != nil {
		return err
	}
	for _, info := range infos {
		projectName := info.Project
		if info.StartTime == -1 {
			preparingCounter[projectName]++
		} else {
			runningTime := time.Now().UnixMilli() - info.StartTime
			if inRange(runningTime, 0, 3600000) {
				running0Counter[projectName]++
			} else if inRange(runningTime, 3600000, 18000000) {
				running60Counter[projectName]++
			} else if inRange(runningTime, 18000000, 86400000) {
				running300Counter[projectName]++
			} else {
				running1440Counter[projectName]++
			}
			runningCounter[projectName]++
		}
	}
	for projectName, num := range preparingCounter {
		ch <- c.preparing.MustNewConstMetric(float64(num), projectName)
	}
	for projectName, num := range runningCounter {
		ch <- c.running.MustNewConstMetric(float64(num), projectName)
	}
	for projectName, num := range running0Counter {
		ch <- c.running0.MustNewConstMetric(float64(num), projectName)
	}
	for projectName, num := range running60Counter {
		ch <- c.running60.MustNewConstMetric(float64(num), projectName)
	}
	for projectName, num := range running300Counter {
		ch <- c.running300.MustNewConstMetric(float64(num), projectName)
	}
	for projectName, num := range running1440Counter {
		ch <- c.running1440.MustNewConstMetric(float64(num), projectName)
	}
	return nil
}

// inRange determine whether a number belongs to a range.
// Will determine target number in [start number, end number)
func inRange(target int64, start int64, end int64) bool {
	if end <= start {
		panic("Wrong value of arguments.")
	} else {
		return target >= start && target < end
	}
}
