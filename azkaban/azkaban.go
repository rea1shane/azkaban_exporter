package azkaban

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
)

type Address struct {
	Protocol string `yaml:"protocol"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	url      string
}

type User struct {
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	sessionId string
}

type Azkaban struct {
	Address Address `yaml:"address"`
	Users   []User  `yaml:"users"`
}

var instance *Azkaban
var once sync.Once

func GetInstance() *Azkaban {
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
		instance.Address.url = instance.Address.Protocol + "://" + instance.Address.Host + ":" + instance.Address.Port
	})
	return instance
}
