package azkaban

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
)

type Server struct {
	Protocol string `yaml:"protocol"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Url      string
}

type User struct {
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	SessionId string // SessionId 有效期 24 小时 // TODO 逾期自动续期
}

type Project struct {
	ProjectId   int    `json:"projectId"`
	ProjectName string `json:"projectName"`
	CreatedBy   string `json:"createdBy"`
	//CreatedTime time.Time `json:"createdTime"`
	//userPermissions
	//groupPermissions
}

type Flow struct {
	FlowId string `json:"flowId"`
}

// Azkaban
// TODO 更改为单个用户, 减少请求次数, session 检测成本
type Azkaban struct {
	Server Server `yaml:"server"`
	Users  []User `yaml:"users"`
	//LoginTime time.Time
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
