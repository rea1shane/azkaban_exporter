package azkaban

import (
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

// Session
// TODO 逾期自动续期
// TODO api.Login 后记录时间
type Session struct {
	SessionId  string // SessionId 有效期 24 小时
	ExpiryTime time.Time
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
