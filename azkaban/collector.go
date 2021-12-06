package azkaban

import (
	"azkaban_exporter/required"
	"azkaban_exporter/util"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	execSubsystem = "execution"
)

var (
	azkaban     = GetAzkaban()
	runningExec []int
)

func init() {
	util.RegisterCollector(execSubsystem, util.DefaultEnabled, NewAzkabanCollector)
}

type azkabanCollector struct {
	logger          log.Logger
	runningExecDesc util.TypedDesc
}

func NewAzkabanCollector(namespace string, logger log.Logger) (required.Collector, error) {
	return &azkabanCollector{
		logger: logger,
		runningExecDesc: util.TypedDesc{
			Desc:      prometheus.NewDesc(prometheus.BuildFQName(namespace, execSubsystem, "running"), "The number of running execution.", nil, nil),
			ValueType: prometheus.GaugeValue},
	}, nil
}

func (c azkabanCollector) Update(ch chan<- prometheus.Metric) error {
	err := recordRunningExec()
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	ch <- c.runningExecDesc.MustNewConstMetric(float64(len(runningExec)))
	return nil
}

func recordRunningExec() error {
	ids, err := azkaban.GetRunningExecIds()
	if err != nil {
		return err
	}
	runningExec = ids
	return nil
}
