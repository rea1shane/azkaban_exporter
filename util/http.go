package util

import (
	"encoding/json"
	"errors"
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

func (h *SingletonHttp) Request(req *http.Request, responseStruct interface{}) error {
	res, err := h.Client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(responseBody, &responseStruct); err != nil {
		return err
	}
	return nil
}

// TODO rename
func RequestFailureError(apiName string, reason string) error {
	return errors.New("request failure when call " + apiName + " api, reason: " + reason)
}
