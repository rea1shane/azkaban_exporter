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

func (a *Azkaban) GetProjectNames(ch chan<- string) error {
	err := a.auth()
	if err != nil {
		return err
	}
	projects, err := api.FetchUserProjects(a.Server.Url, a.User.Session.SessionId)
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	for _, project := range projects {
		wg.Add(1)
		go func(project api.Project) {
			defer wg.Done()
			ch <- project.ProjectName
		}(project)
	}
	wg.Wait()
	close(ch)
	return nil
}

func (a *Azkaban) GetRunningExecIds(ch chan<- int) error {
	err := a.auth()
	if err != nil {
		return err
	}
	projects, err := api.FetchUserProjects(a.Server.Url, a.User.Session.SessionId)
	if err != nil {
		return err
	}
	wgProjects := sync.WaitGroup{}
	for _, project := range projects {
		wgProjects.Add(1)
		go func(project api.Project) {
			defer wgProjects.Done()
			flows, err := api.FetchFlowsOfAProject(a.Server.Url, a.User.Session.SessionId, project.ProjectName)
			if err != nil {
				// TODO 处理 panic
				panic(fmt.Errorf(err.Error()))
			}
			wgFlows := sync.WaitGroup{}
			for _, flow := range flows {
				wgFlows.Add(1)
				go func(flow api.Flow) {
					defer wgFlows.Done()
					runningExecutions, err := api.FetchRunningExecutionsOfAFlow(a.Server.Url, a.User.Session.SessionId, project.ProjectName, flow.FlowId)
					if err != nil {
						// TODO 处理 panic
						panic(fmt.Errorf(err.Error()))
					}
					wgExecs := sync.WaitGroup{}
					for _, execId := range runningExecutions.ExecIds {
						wgExecs.Add(1)
						go func(execId int) {
							wgExecs.Done()
							ch <- execId
						}(execId)
					}
					wgExecs.Wait()
				}(flow)
			}
			wgFlows.Wait()
		}(project)
	}
	wgProjects.Wait()
	close(ch)
	return nil
}

func (a *Azkaban) GetExecInfos(execIds <-chan int, execInfos chan<- api.ExecInfo) error {
	err := a.auth()
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	for execId := range execIds {
		wg.Add(1)
		go func(execId int) {
			defer wg.Done()
			execInfo, err := api.FetchAFlowExecution(a.Server.Url, a.User.Session.SessionId, execId)
			if err != nil {
				// TODO 处理 panic
				panic(fmt.Errorf(err.Error()))
			}
			execInfos <- execInfo
		}(execId)
	}
	wg.Wait()
	close(execInfos)
	return nil
}
