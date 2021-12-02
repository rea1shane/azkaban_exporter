package azkaban

import (
	"azkaban_exporter/util"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"net/http"
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
		instance.Server.url = instance.Server.Protocol + "://" + instance.Server.Host + ":" + instance.Server.Port
	})
	return instance
}

func (a *Azkaban) Login() {
	response := LoginResponse{}
	method := "POST"
	payload := strings.NewReader("action=login&username=" + a.User.Username + "&password=" + a.User.Password)
	req, err := http.NewRequest(method, a.Server.url, payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := singletonHttp.Client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println(err.Error())
	}
	a.User.sessionId = response.Sessionid
}
