package azkaban

import (
	"azkaban_exporter/azkaban/api"
	"context"
	"github.com/go-kratos/kratos/pkg/sync/errgroup"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
	"time"
)

type Server struct {
	Protocol string `yaml:"protocol"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Url      string
}

type Session struct {
	SessionId     string // SessionId 默认有效期 24 小时
	AuthTimestamp int64
}

type User struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Session  Session
}

type Azkaban struct {
	Server Server `yaml:"server"`
	User   User   `yaml:"user"`
}

type ProjectWithFlows struct {
	ProjectName string
	FlowIds     []string
}

type Execution struct {
	SubmitTime  int64
	SubmitUser  string
	StartTime   int64
	EndTime     int64
	ProjectName string
	FlowID      string
	ExecID      int
	Status      string
}

var instance *Azkaban
var once sync.Once

func GetAzkaban() *Azkaban {
	once.Do(func() {
		yamlFile, err := ioutil.ReadFile(getAzkabanConfPath())
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(yamlFile, &instance)
		if err != nil {
			panic(err)
		}
		instance.Server.Url = instance.Server.Protocol + "://" + instance.Server.Host + ":" + instance.Server.Port
	})
	return instance
}

func (a *Azkaban) GetProjectWithFlows(ctx context.Context, ch chan<- ProjectWithFlows) error {
	err := a.auth(ctx)
	if err != nil {
		return err
	}
	projects, err := api.FetchUserProjects(api.FetchUserProjectsParam{
		ServerUrl: a.Server.Url,
		SessionId: a.User.Session.SessionId,
	}, ctx)
	if err != nil {
		return err
	}
	group := errgroup.WithCancel(ctx)
	for _, project := range projects {
		p := project
		group.Go(func(ctx context.Context) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				flows, err := api.FetchFlowsOfAProject(api.FetchFlowsOfAProjectParam{
					ServerUrl:   a.Server.Url,
					SessionId:   a.User.Session.SessionId,
					ProjectName: p.ProjectName,
				}, ctx)
				if err != nil {
					return err
				}
				var ids []string
				for _, flow := range flows {
					ids = append(ids, flow.FlowId)
				}
				ch <- ProjectWithFlows{
					ProjectName: p.ProjectName,
					FlowIds:     ids,
				}
				return nil
			}
		})
	}
	return group.Wait()
}

func (a *Azkaban) GetExecutions(ctx context.Context, projectName string, flowId string, startIndex int, listLength int, ch chan<- Execution) error {
	Executions, err := api.FetchExecutionsOfAFlow(api.FetchExecutionsOfAFlowParam{
		ServerUrl:   a.Server.Url,
		SessionId:   a.User.Session.SessionId,
		ProjectName: projectName,
		FlowId:      flowId,
		StartIndex:  startIndex,
		ListLength:  listLength,
	}, ctx)
	if err != nil {
		return err
	}
	if startIndex == 0 && len(Executions.Executions) == 0 {
		ch <- Execution{
			ProjectName: projectName,
			FlowID:      flowId,
			Status:      "NEVER RUN",
		}
		return nil
	}
	for _, execution := range Executions.Executions {
		ch <- Execution{
			SubmitTime:  execution.SubmitTime,
			SubmitUser:  execution.SubmitUser,
			StartTime:   execution.StartTime,
			EndTime:     execution.EndTime,
			ProjectName: projectName,
			FlowID:      execution.FlowID,
			ExecID:      execution.ExecID,
			Status:      execution.Status,
		}
	}
	return nil
}

// auth and check session < 23h:50m
func (a *Azkaban) auth(ctx context.Context) error {
	if a.User.Session.AuthTimestamp != 0 && time.Now().Unix()-a.User.Session.AuthTimestamp < 85800 {
		return nil
	}
	sessionId, err := api.Authenticate(api.AuthenticateParam{
		ServerUrl: a.Server.Url,
		Username:  a.User.Username,
		Password:  a.User.Password,
	}, ctx)
	if err != nil {
		return err
	}
	a.User.Session.SessionId = sessionId
	a.User.Session.AuthTimestamp = time.Now().Unix()
	return nil
}
