package pkg

import (
	"context"
	"github.com/go-kratos/kratos/pkg/sync/errgroup"
	"github.com/rea1shane/azkaban_exporter/pkg/api"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
	"time"
)

type Server struct {
	Protocol string `yaml:"protocol"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	url      string
}

type session struct {
	sessionId     string // sessionId 默认有效期 24 小时
	authTimestamp int64
}

type User struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	session  session
}

type Azkaban struct {
	Server Server `yaml:"server"`
	User   User   `yaml:"user"`
}

type projectWithFlows struct {
	projectName string
	flowIds     []string
}

type execution struct {
	submitTime  int64
	submitUser  string
	startTime   int64
	endTime     int64
	projectName string
	flowID      string
	execID      int
	status      string
}

var instance *Azkaban
var once sync.Once

func getAzkaban() *Azkaban {
	once.Do(func() {
		yamlFile, err := ioutil.ReadFile(getAzkabanConfPath())
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(yamlFile, &instance)
		if err != nil {
			panic(err)
		}
		instance.Server.url = instance.Server.Protocol + "://" + instance.Server.Host + ":" + instance.Server.Port
	})
	return instance
}

func (a *Azkaban) getProjectWithFlows(ctx context.Context, ch chan<- projectWithFlows) error {
	err := a.auth(ctx)
	if err != nil {
		return err
	}
	projects, err := api.FetchUserProjects(api.FetchUserProjectsParam{
		ServerUrl: a.Server.url,
		SessionId: a.User.session.sessionId,
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
					ServerUrl:   a.Server.url,
					SessionId:   a.User.session.sessionId,
					ProjectName: p.ProjectName,
				}, ctx)
				if err != nil {
					return err
				}
				var ids []string
				for _, flow := range flows {
					ids = append(ids, flow.FlowId)
				}
				ch <- projectWithFlows{
					projectName: p.ProjectName,
					flowIds:     ids,
				}
				return nil
			}
		})
	}
	return group.Wait()
}

func (a *Azkaban) getExecutions(ctx context.Context, projectName string, flowId string, startIndex int, listLength int, ch chan<- execution) error {
	Executions, err := api.FetchExecutionsOfAFlow(api.FetchExecutionsOfAFlowParam{
		ServerUrl:   a.Server.url,
		SessionId:   a.User.session.sessionId,
		ProjectName: projectName,
		FlowId:      flowId,
		StartIndex:  startIndex,
		ListLength:  listLength,
	}, ctx)
	if err != nil {
		return err
	}
	if startIndex == 0 && len(Executions.Executions) == 0 {
		ch <- execution{
			projectName: projectName,
			flowID:      flowId,
			status:      "NEW",
		}
		return nil
	}
	for _, exec := range Executions.Executions {
		ch <- execution{
			submitTime:  exec.SubmitTime,
			submitUser:  exec.SubmitUser,
			startTime:   exec.StartTime,
			endTime:     exec.EndTime,
			projectName: projectName,
			flowID:      exec.FlowID,
			execID:      exec.ExecID,
			status:      exec.Status,
		}
	}
	return nil
}

// auth and check session < 23h:50m
func (a *Azkaban) auth(ctx context.Context) error {
	if a.User.session.authTimestamp != 0 && time.Now().Unix()-a.User.session.authTimestamp < 85800 {
		return nil
	}
	sessionId, err := api.Authenticate(api.AuthenticateParam{
		ServerUrl: a.Server.url,
		Username:  a.User.Username,
		Password:  a.User.Password,
	}, ctx)
	if err != nil {
		return err
	}
	a.User.session.sessionId = sessionId
	a.User.session.authTimestamp = time.Now().Unix()
	return nil
}
