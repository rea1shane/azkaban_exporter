package util

import (
	"net/http"
	"sync"
)

type SingletonHttp struct {
	Client *http.Client
}

var instance *SingletonHttp
var once sync.Once

func GetSingletonHttp() *SingletonHttp {
	once.Do(func() {
		instance = &SingletonHttp{
			Client: &http.Client{},
		}
	})
	return instance
}
