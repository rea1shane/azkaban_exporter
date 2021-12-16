package util

import (
	"context"
	"encoding/json"
	"github.com/morikuni/failure"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
)

type SingletonHttp struct {
	client *http.Client
}

var instance *SingletonHttp
var once sync.Once

func GetSingletonHttp() *SingletonHttp {
	once.Do(func() {
		instance = &SingletonHttp{
			client: &http.Client{},
		}
	})
	return instance
}

func (h *SingletonHttp) Request(req *http.Request, ctx context.Context, responseStruct interface{}) error {
	req = req.WithContext(ctx)
	res, err := h.client.Do(req)
	if err != nil {
		return failure.Wrap(err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return failure.Wrap(err)
	}
	if err = json.Unmarshal(responseBody, &responseStruct); err != nil {
		return failure.Wrap(err)
	}
	return nil
}
