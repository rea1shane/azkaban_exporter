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
	SessionId     string // SessionId 有效期 24 小时
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

var instance *Azkaban
var once sync.Once

func GetAzkaban() *Azkaban {
	once.Do(func() {
		// TODO 使用传参传入配置文件路径
		yamlFile, err := ioutil.ReadFile("azkaban/conf/azkaban.yml")
		if err != nil {
			// TODO 程序结束
			fmt.Println(err.Error())
		}
		err = yaml.Unmarshal(yamlFile, &instance)
		if err != nil {
			// TODO 程序结束
			fmt.Println(err.Error())
		}
		instance.Server.Url = instance.Server.Protocol + "://" + instance.Server.Host + ":" + instance.Server.Port
	})
	return instance
}

func (a *Azkaban) auth() error {
	if a.User.Session.AuthTimestamp != 0 && time.Now().Unix()-a.User.Session.AuthTimestamp < 85800 { // session < 23h:50m
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

func (a *Azkaban) GetRunningExecIds() ([]int, error) {
	var runningExecIds []int
	err := a.auth()
	if err != nil {
		return nil, err
	}
	projects, err := api.FetchUserProjects(a.Server.Url, a.User.Session.SessionId)
	if err != nil {
		return nil, err
	}
	for _, project := range projects {
		flows, err := api.FetchFlowsOfAProject(a.Server.Url, a.User.Session.SessionId, project.ProjectName)
		if err != nil {
			return nil, err
		}
		for _, flow := range flows {
			runningExecutions, err := api.FetchRunningExecutionsOfAFlow(a.Server.Url, a.User.Session.SessionId, project.ProjectName, flow.FlowId)
			if err != nil {
				return nil, err
			}
			runningExecIds = append(runningExecIds, runningExecutions.ExecIds...)
		}
	}
	return runningExecIds, nil
}

func (a *Azkaban) GetExecInfo(execId int) (api.ExecutionInformationResponse, error) {
	err := a.auth()
	if err != nil {
		return api.ExecutionInformationResponse{}, err
	}
	execution, err := api.FetchAFlowExecution(a.Server.Url, a.User.Session.SessionId, execId)
	return execution, nil
}
