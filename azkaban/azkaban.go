package azkaban

import (
	"azkaban_exporter/util"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var singletonHttp = util.GetSingletonHttp()

type Server struct {
	Protocol string `yaml:"protocol"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	url      string
}

type User struct {
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	sessionId string // sessionId 有效期 24 小时 // TODO 逾期自动续期
}

type Azkaban struct {
	Server Server `yaml:"server"`
	Users  []User `yaml:"users"`
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
		instance.Server.url = instance.Server.Protocol + "://" + instance.Server.Host + ":" + instance.Server.Port
	})
	return instance
}

func (a *Azkaban) Login() {
	method := "POST"
	for index, user := range a.Users {
		response := LoginResponse{}
		payload := strings.NewReader("action=login&username=" + user.Username + "&password=" + user.Password)
		req, err := http.NewRequest(method, a.Server.url, payload)
		if err != nil {
			fmt.Println(err)
			return
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		singletonHttp.Request(req, &response)
		a.Users[index].sessionId = response.Sessionid
	}
}

func (a Azkaban) GetProjectIds() []string {
	method := "GET"
	var ids []string
	for _, user := range a.Users {
		response := GetProjectsResponse{}
		url := a.Server.url + "/index?ajax=fetchuserprojects&session.id=" + user.sessionId
		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		singletonHttp.Request(req, &response)
		for _, project := range response.Projects {
			ids = append(ids, strconv.Itoa(project.ProjectId))
		}
	}
	return ids
}
