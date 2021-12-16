package azkaban

import (
	"azkaban_exporter/required/functions"
	"azkaban_exporter/required/structs"
	"azkaban_exporter/util"
	"context"
	"github.com/go-kit/log"
	"github.com/go-kratos/kratos/pkg/sync/errgroup"
	"github.com/morikuni/failure"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

const (
	subsystem  = "flow"
	startIndex = 0
	listLength = 1
)

var (
	azkaban = GetAzkaban()
)

func init() {
	functions.RegisterCollector(subsystem, util.DefaultEnabled, NewAzkabanCollector)
}

type azkabanCollector struct {
	logger      log.Logger
	projects    util.TypedDesc
	preparing   util.TypedDesc
	running     util.TypedDesc
	running0    util.TypedDesc
	running60   util.TypedDesc
	running300  util.TypedDesc
	running1440 util.TypedDesc
	lastStatus  util.TypedDesc
}

func NewAzkabanCollector(namespace string, logger log.Logger) (structs.Collector, error) {
	var (
		labelProject     = []string{"project"}
		labelProjectFlow = []string{"project", "flow"}
	)

	return &azkabanCollector{
		logger: logger,
		projects: util.TypedDesc{
			Desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "projects"),
				"The number of projects", nil, nil),
			ValueType: prometheus.GaugeValue,
		},
		preparing: util.TypedDesc{
			Desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "preparing"),
				"The number of preparing start flows", labelProject, nil),
			ValueType: prometheus.GaugeValue,
		},
		running: util.TypedDesc{
			Desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "running"),
				"The number of running flows.", labelProject, nil),
			ValueType: prometheus.GaugeValue,
		},
		running0: util.TypedDesc{
			Desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "running_0"),
				"The number of running flows which running time in [0, 60) mins.", labelProject, nil),
			ValueType: prometheus.GaugeValue,
		},
		running60: util.TypedDesc{
			Desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "running_60"),
				"The number of running flows which running time in [60, 300) mins.", labelProject, nil),
			ValueType: prometheus.GaugeValue,
		},
		running300: util.TypedDesc{
			Desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "running_300"),
				"The number of running flows which running time in [300, 1440) mins.", labelProject, nil),
			ValueType: prometheus.GaugeValue,
		},
		running1440: util.TypedDesc{
			Desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "running_1440"),
				"The number of running flows which running time over 1440 mins.", labelProject, nil),
			ValueType: prometheus.GaugeValue,
		},
		lastStatus: util.TypedDesc{
			Desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "last_status"),
				"flow last execution status flag, (-1=UNKNOW / 0=FAILED / 1=SUCCEEDED / 2=RUNNING / 3=PREPARING)", labelProjectFlow, nil),
			ValueType: prometheus.GaugeValue,
		},
	}, nil
}

func (c azkabanCollector) Update(ch chan<- prometheus.Metric) error {
	var (
		projectsWithFlows = make(chan ProjectWithFlows)
		executions        = make(chan Execution)

		projects = 0

		preparingCounter      = map[string]int{}
		runningCounter        = map[string]int{}
		running0Counter       = map[string]int{}
		running60Counter      = map[string]int{}
		running300Counter     = map[string]int{}
		running1440Counter    = map[string]int{}
		runningAttemptCounter = map[string]int{}

		lastStatusRecorder = map[string]map[string]int{}
	)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancelFunc()
	group := errgroup.WithCancel(ctx)
	group.Go(func(ctx context.Context) error {
		err := azkaban.GetProjectWithFlows(ctx, projectsWithFlows)
		close(projectsWithFlows)
		return err
	})
	group.Go(func(ctx context.Context) error {
		g := errgroup.WithCancel(ctx)
		for projectWithFlows := range projectsWithFlows {
			projectName := projectWithFlows.ProjectName
			flowIds := projectWithFlows.FlowIds
			projects++
			preparingCounter[projectName] = 0
			runningCounter[projectName] = 0
			running0Counter[projectName] = 0
			running60Counter[projectName] = 0
			running300Counter[projectName] = 0
			running1440Counter[projectName] = 0
			runningAttemptCounter[projectName] = 0
			lastStatusRecorder[projectName] = map[string]int{}
			g.Go(func(ctx context.Context) error {
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					for _, flowId := range flowIds {
						fid := flowId
						g.Go(func(ctx context.Context) error {
							select {
							case <-ctx.Done():
								return ctx.Err()
							default:
								return azkaban.GetExecutions(ctx, projectName, fid, startIndex, listLength, executions)
							}
						})
					}
					return nil
				}
			})
		}
		err := g.Wait()
		close(executions)
		return err
	})
	for execution := range executions {
		switch execution.Status {
		case "PREPARING":
			lastStatusRecorder[execution.ProjectName][execution.FlowID] = 3
			preparingCounter[execution.ProjectName]++
		case "RUNNING":
			lastStatusRecorder[execution.ProjectName][execution.FlowID] = 2
			runningTime := time.Now().UnixMilli() - execution.StartTime
			if inRange(runningTime, 0, 3600000) {
				running0Counter[execution.ProjectName]++
			} else if inRange(runningTime, 3600000, 18000000) {
				running60Counter[execution.ProjectName]++
			} else if inRange(runningTime, 18000000, 86400000) {
				running300Counter[execution.ProjectName]++
			} else {
				running1440Counter[execution.ProjectName]++
			}
			runningCounter[execution.ProjectName]++
		case "SUCCEEDED":
			lastStatusRecorder[execution.ProjectName][execution.FlowID] = 1
		case "FAILED":
			lastStatusRecorder[execution.ProjectName][execution.FlowID] = 0
		default:
			lastStatusRecorder[execution.ProjectName][execution.FlowID] = -1
		}
	}
	group.Go(func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			ch <- c.projects.MustNewConstMetric(float64(projects))
			return nil
		}
	})
	group.Go(func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			for projectName, num := range preparingCounter {
				ch <- c.preparing.MustNewConstMetric(float64(num), projectName)
			}
			return nil
		}
	})
	group.Go(func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			for projectName, num := range runningCounter {
				ch <- c.running.MustNewConstMetric(float64(num), projectName)
			}
			return nil
		}
	})
	group.Go(func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			for projectName, num := range running0Counter {
				ch <- c.running0.MustNewConstMetric(float64(num), projectName)
			}
			return nil
		}
	})
	group.Go(func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			for projectName, num := range running60Counter {
				ch <- c.running60.MustNewConstMetric(float64(num), projectName)
			}
			return nil
		}
	})
	group.Go(func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			for projectName, num := range running300Counter {
				ch <- c.running300.MustNewConstMetric(float64(num), projectName)
			}
			return nil
		}
	})
	group.Go(func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			for projectName, num := range running1440Counter {
				ch <- c.running1440.MustNewConstMetric(float64(num), projectName)
			}
			return nil
		}
	})
	group.Go(func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			for projectName, flowAndflag := range lastStatusRecorder {
				for flowId, flag := range flowAndflag {
					ch <- c.lastStatus.MustNewConstMetric(float64(flag), projectName, flowId)
				}
			}
			return nil
		}
	})
	err := group.Wait()
	if ctx.Err() != nil {
		return failure.Wrap(ctx.Err())
	}
	return err
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
