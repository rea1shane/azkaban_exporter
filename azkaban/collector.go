package azkaban

import (
	"azkaban_exporter/required/functions"
	"azkaban_exporter/required/structs"
	"azkaban_exporter/util"
	"context"
	"github.com/go-kratos/kratos/pkg/sync/errgroup"
	"github.com/morikuni/failure"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"time"
)

var (
	runningExecs          []int
	totalSucceededCounter = cmap.New()
	totalFailedCounter    = cmap.New()
	totalKilledCounter    = cmap.New()
)

const (
	subsystem  = "flow"
	startIndex = 0
	listLength = 1
)

func init() {
	functions.RegisterCollector(subsystem, util.DefaultEnabled, NewAzkabanCollector)
}

type azkabanCollector struct {
	logger          *log.Entry
	new             util.TypedDesc
	preparing       util.TypedDesc
	running         util.TypedDesc
	succeeded       util.TypedDesc
	failed          util.TypedDesc
	unknow          util.TypedDesc
	killed          util.TypedDesc
	running0        util.TypedDesc
	running60       util.TypedDesc
	running300      util.TypedDesc
	running1440     util.TypedDesc
	runningDuration util.TypedDesc
	totalSucceeded  util.TypedDesc
	totalFailed     util.TypedDesc
	totalKilled     util.TypedDesc
	lastStatus      util.TypedDesc
	lastDuration    util.TypedDesc
}

func NewAzkabanCollector(namespace string, logger *log.Entry) (structs.Collector, error) {
	var (
		labelProject     = []string{"project"}
		labelProjectFlow = []string{"project", "flow"}
	)

	return &azkabanCollector{
		logger: logger,
		new: util.TypedDesc{
			Desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "new"),
				"The number of never run flows", labelProject, nil),
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
		succeeded: util.TypedDesc{
			Desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "succeeded"),
				"The number of flows that last status is succeeded.", labelProject, nil),
			ValueType: prometheus.GaugeValue,
		},
		failed: util.TypedDesc{
			Desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "failed"),
				"The number of flows that last status is failed.", labelProject, nil),
			ValueType: prometheus.GaugeValue,
		},
		unknow: util.TypedDesc{
			Desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "unknow"),
				"The number of flows that last status is unknow.", labelProject, nil),
			ValueType: prometheus.GaugeValue,
		},
		killed: util.TypedDesc{
			Desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "killed"),
				"The number of flows that last status is killed.", labelProject, nil),
			ValueType: prometheus.GaugeValue,
		},
		running0: util.TypedDesc{
			Desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "running_0"),
				"The number of running flows that duration in [0, 60) mins.", labelProject, nil),
			ValueType: prometheus.GaugeValue,
		},
		running60: util.TypedDesc{
			Desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "running_60"),
				"The number of running flows that duration in [60, 300) mins.", labelProject, nil),
			ValueType: prometheus.GaugeValue,
		},
		running300: util.TypedDesc{
			Desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "running_300"),
				"The number of running flows that duration in [300, 1440) mins.", labelProject, nil),
			ValueType: prometheus.GaugeValue,
		},
		running1440: util.TypedDesc{
			Desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "running_1440"),
				"The number of running flows that duration over 1440 mins.", labelProject, nil),
			ValueType: prometheus.GaugeValue,
		},
		totalSucceeded: util.TypedDesc{
			Desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "succeeded_total"),
				"The total number of succeeded.", labelProject, nil),
			ValueType: prometheus.CounterValue,
		},
		totalFailed: util.TypedDesc{
			Desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "failed_total"),
				"The total number of failed.", labelProject, nil),
			ValueType: prometheus.CounterValue,
		},
		totalKilled: util.TypedDesc{
			Desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "killed_total"),
				"The total number of killed.", labelProject, nil),
			ValueType: prometheus.CounterValue,
		},
		runningDuration: util.TypedDesc{
			Desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "running_duration_ms"),
				"Duration of each running flows. (unit: ms)", labelProjectFlow, nil),
			ValueType: prometheus.GaugeValue,
		},
		lastStatus: util.TypedDesc{
			Desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "last_status"),
				"Flow last execution status flag. (-2=KILLED / -1=UNKNOW / 0=FAILED / 1=SUCCEEDED / 2=RUNNING / 3=PREPARING / 4=NEW)", labelProjectFlow, nil),
			ValueType: prometheus.GaugeValue,
		},
		lastDuration: util.TypedDesc{
			Desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "last_duration_ms"),
				"Flow last execution duration which finished. (unit: ms)", labelProjectFlow, nil),
			ValueType: prometheus.GaugeValue,
		},
	}, nil
}

func (c azkabanCollector) Update(ch chan<- prometheus.Metric) error {
	var (
		azkaban = GetAzkaban()

		projectsWithFlows = make(chan ProjectWithFlows)
		executions        = make(chan Execution)

		newCounter       = map[string]int{}
		preparingCounter = map[string]int{}
		runningCounter   = map[string]int{}
		succeededCounter = map[string]int{}
		failedCounter    = map[string]int{}
		unknowCounter    = map[string]int{}
		killedCounter    = map[string]int{}

		running0Counter       = map[string]int{}
		running60Counter      = map[string]int{}
		running300Counter     = map[string]int{}
		running1440Counter    = map[string]int{}
		runningAttemptCounter = map[string]int{}

		runningDurationRecorder = map[string]map[string]int64{}
		lastStatusRecorder      = map[string]map[string]int{}
		lastDurationRecorder    = map[string]map[string]int64{}
	)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5000*time.Millisecond)
	defer cancelFunc()
	group := errgroup.WithCancel(ctx)
	group.Go(func(ctx context.Context) error {
		err := azkaban.GetProjectWithFlows(ctx, projectsWithFlows)
		close(projectsWithFlows)
		return err
	})
	group.Go(func(ctx context.Context) error {
		var projectNames []string
		g := errgroup.WithCancel(ctx)
		for projectWithFlows := range projectsWithFlows {
			projectName := projectWithFlows.ProjectName
			flowIds := projectWithFlows.FlowIds
			projectNames = append(projectNames, projectName)

			// TODO map 并发安全问题
			newCounter[projectName] = 0
			preparingCounter[projectName] = 0
			runningCounter[projectName] = 0
			succeededCounter[projectName] = 0
			failedCounter[projectName] = 0
			unknowCounter[projectName] = 0
			killedCounter[projectName] = 0

			running0Counter[projectName] = 0
			running60Counter[projectName] = 0
			running300Counter[projectName] = 0
			running1440Counter[projectName] = 0
			runningAttemptCounter[projectName] = 0

			runningDurationRecorder[projectName] = map[string]int64{}
			lastStatusRecorder[projectName] = map[string]int{}
			lastDurationRecorder[projectName] = map[string]int64{}

			totalSucceededCounter.SetIfAbsent(projectName, 0)
			totalFailedCounter.SetIfAbsent(projectName, 0)
			totalKilledCounter.SetIfAbsent(projectName, 0)
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
		removeKeys(totalSucceededCounter, projectNames)
		removeKeys(totalFailedCounter, projectNames)
		removeKeys(totalKilledCounter, projectNames)
		return err
	})
	for execution := range executions {
		switch execution.Status {
		case "NEW":
			lastStatusRecorder[execution.ProjectName][execution.FlowID] = 4
			newCounter[execution.ProjectName]++
		case "PREPARING":
			lastStatusRecorder[execution.ProjectName][execution.FlowID] = 3
			preparingCounter[execution.ProjectName]++
		case "RUNNING":
			if _, ok := findInt(runningExecs, execution.ExecID); !ok {
				runningExecs = append(runningExecs, execution.ExecID)
			}
			lastStatusRecorder[execution.ProjectName][execution.FlowID] = 2
			runningTime := time.Now().UnixMilli() - execution.StartTime
			runningDurationRecorder[execution.ProjectName][execution.FlowID] = runningTime
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
			if index, ok := findInt(runningExecs, execution.ExecID); ok {
				value, _ := totalSucceededCounter.Get(execution.ProjectName)
				totalSucceededCounter.Set(execution.ProjectName, value.(int)+1)
				runningExecs = append(runningExecs[:index], runningExecs[index+1:]...)
			}
			lastStatusRecorder[execution.ProjectName][execution.FlowID] = 1
			lastDurationRecorder[execution.ProjectName][execution.FlowID] = execution.EndTime - execution.StartTime
			succeededCounter[execution.ProjectName]++
		case "FAILED":
			if index, ok := findInt(runningExecs, execution.ExecID); ok {
				value, _ := totalFailedCounter.Get(execution.ProjectName)
				totalFailedCounter.Set(execution.ProjectName, value.(int)+1)
				runningExecs = append(runningExecs[:index], runningExecs[index+1:]...)
			}
			lastStatusRecorder[execution.ProjectName][execution.FlowID] = 0
			lastDurationRecorder[execution.ProjectName][execution.FlowID] = execution.EndTime - execution.StartTime
			failedCounter[execution.ProjectName]++
		case "KILLED":
			if index, ok := findInt(runningExecs, execution.ExecID); ok {
				value, _ := totalKilledCounter.Get(execution.ProjectName)
				totalKilledCounter.Set(execution.ProjectName, value.(int)+1)
				runningExecs = append(runningExecs[:index], runningExecs[index+1:]...)
			}
			lastStatusRecorder[execution.ProjectName][execution.FlowID] = -2
			lastDurationRecorder[execution.ProjectName][execution.FlowID] = execution.EndTime - execution.StartTime
			killedCounter[execution.ProjectName]++
		default:
			lastStatusRecorder[execution.ProjectName][execution.FlowID] = -1
			unknowCounter[execution.ProjectName]++
		}
	}
	group.Go(func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			for projectName, num := range newCounter {
				ch <- c.new.MustNewConstMetric(float64(num), projectName)
			}
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
			for projectName, num := range succeededCounter {
				ch <- c.succeeded.MustNewConstMetric(float64(num), projectName)
			}
			return nil
		}
	})
	group.Go(func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			for projectName, num := range failedCounter {
				ch <- c.failed.MustNewConstMetric(float64(num), projectName)
			}
			return nil
		}
	})
	group.Go(func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			for projectName, num := range unknowCounter {
				ch <- c.unknow.MustNewConstMetric(float64(num), projectName)
			}
			return nil
		}
	})
	group.Go(func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			for projectName, num := range killedCounter {
				ch <- c.killed.MustNewConstMetric(float64(num), projectName)
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
			totalSucceededCounter.IterCb(func(key string, v interface{}) {
				ch <- c.totalSucceeded.MustNewConstMetric(float64(v.(int)), key)
			})
			return nil
		}
	})
	group.Go(func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			totalFailedCounter.IterCb(func(key string, v interface{}) {
				ch <- c.totalFailed.MustNewConstMetric(float64(v.(int)), key)
			})
			return nil
		}
	})
	group.Go(func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			totalKilledCounter.IterCb(func(key string, v interface{}) {
				ch <- c.totalKilled.MustNewConstMetric(float64(v.(int)), key)
			})
			return nil
		}
	})
	group.Go(func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			for projectName, flowAndDuration := range runningDurationRecorder {
				for flowId, duration := range flowAndDuration {
					ch <- c.runningDuration.MustNewConstMetric(float64(duration), projectName, flowId)
				}
			}
			return nil
		}
	})
	group.Go(func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			for projectName, flowAndFlag := range lastStatusRecorder {
				for flowId, flag := range flowAndFlag {
					ch <- c.lastStatus.MustNewConstMetric(float64(flag), projectName, flowId)
				}
			}
			return nil
		}
	})
	group.Go(func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			for projectName, flowAndDuration := range lastDurationRecorder {
				for flowId, duration := range flowAndDuration {
					ch <- c.lastDuration.MustNewConstMetric(float64(duration), projectName, flowId)
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

func findInt(slice []int, val int) (int, bool) {
	for index, item := range slice {
		if item == val {
			return index, true
		}
	}
	return -1, false
}

func findString(slice []string, val string) (int, bool) {
	for index, item := range slice {
		if item == val {
			return index, true
		}
	}
	return -1, false
}

// removeKeys if key not in slice, delete.
func removeKeys(m cmap.ConcurrentMap, s []string) {
	keys := m.Keys()
	for _, key := range keys {
		if _, ok := findString(s, key); !ok {
			m.Remove(key)
		}
	}
}
