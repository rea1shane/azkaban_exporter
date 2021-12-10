package azkaban

import (
	"azkaban_exporter/azkaban/api"
	"fmt"
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
	FlowIds     chan string
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
		// TODO 使用传参传入配置文件路径
		yamlFile, err := ioutil.ReadFile("azkaban/conf/azkaban.yml")
		if err != nil {
			panic(fmt.Errorf(err.Error()))
		}
		err = yaml.Unmarshal(yamlFile, &instance)
		if err != nil {
			panic(fmt.Errorf(err.Error()))
		}
		instance.Server.Url = instance.Server.Protocol + "://" + instance.Server.Host + ":" + instance.Server.Port
	})
	return instance
}

func (a *Azkaban) GetProjectWithFlows(ch chan<- ProjectWithFlows) error {
	err := a.auth()
	if err != nil {
		return err
	}
	projects, err := api.FetchUserProjects(a.Server.Url, a.User.Session.SessionId)
	if err != nil {
		return err
	}
	wgProjects := sync.WaitGroup{}
	wgProjects.Add(len(projects))
	for _, project := range projects {
		go func(project api.Project) {
			defer wgProjects.Done()
			elem := ProjectWithFlows{
				ProjectName: project.ProjectName,
				FlowIds:     make(chan string),
			}
			ch <- elem
			flows, err := api.FetchFlowsOfAProject(a.Server.Url, a.User.Session.SessionId, elem.ProjectName)
			if err != nil {
				// TODO 处理 panic
				panic(fmt.Errorf(err.Error()))
			}
			wgFlows := sync.WaitGroup{}
			wgFlows.Add(len(flows))
			for _, flow := range flows {
				go func(flow api.Flow) {
					defer wgFlows.Done()
					elem.FlowIds <- flow.FlowId
				}(flow)
			}
			wgFlows.Wait()
			close(elem.FlowIds)
		}(project)
	}
	wgProjects.Wait()
	return nil
}

func (a *Azkaban) GetExecutions(projectName string, flowId string, startIndex int, listLength int, ch chan<- Execution) error {
	Executions, err := api.FetchExecutionsOfAFlow(a.Server.Url, a.User.Session.SessionId, projectName, flowId, startIndex, listLength)
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	wg.Add(len(Executions.Executions))
	for _, execution := range Executions.Executions {
		go func(execution api.Execution) {
			defer wg.Done()
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
		}(execution)
	}
	wg.Wait()
	return nil
}

// auth and check session < 23h:50m
func (a *Azkaban) auth() error {
	if a.User.Session.AuthTimestamp != 0 && time.Now().Unix()-a.User.Session.AuthTimestamp < 85800 {
		return nil
	}
	sessionId, err := api.Authenticate(a.Server.Url, a.User.Username, a.User.Password)
	if err != nil {
		return err
	}
	a.User.Session.SessionId = sessionId
	a.User.Session.AuthTimestamp = time.Now().Unix()
	return nil
}
