package util

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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

func (h *SingletonHttp) Request(req *http.Request, responseStruct interface{}) {
	res, err := h.Client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := json.Unmarshal(responseBody, &responseStruct); err != nil {
		fmt.Println(err.Error())
	}
}
