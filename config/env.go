package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Exporter struct {
	Port int `yaml:"port"`
}

func (e *Exporter) GetConf(filePath string) *Exporter {
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, e)
	if err != nil {
		fmt.Println(err.Error())
	}
	return e
}
