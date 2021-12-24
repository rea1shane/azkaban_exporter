package util

import (
	"context"
	"encoding/json"
	"github.com/morikuni/failure"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
)

type SingletonHttp struct {
	client *http.Client
}

var (
	instance *SingletonHttp
	onceHttp sync.Once
)

func GetSingletonHttp() *SingletonHttp {
	onceHttp.Do(func() {
		instance = &SingletonHttp{
			client: &http.Client{},
		}
	})
	return instance
}

func (h *SingletonHttp) Request(req *http.Request, ctx context.Context, responseStruct interface{}) error {
	req = req.WithContext(ctx)
	escape(req)
	res, err := h.client.Do(req)
	if err != nil {
		return failure.Wrap(err, failure.Context{
			"protocol":    req.Proto,
			"host":        req.URL.Hostname(),
			"port":        req.URL.Port(),
			"request_url": req.URL.RequestURI(),
		})
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return failure.Wrap(err, failure.Context{
			"protocol":      req.Proto,
			"host":          req.URL.Hostname(),
			"port":          req.URL.Port(),
			"request_url":   req.URL.RequestURI(),
			"response_body": string(responseBody),
		})
	}
	if err = json.Unmarshal(responseBody, &responseStruct); err != nil {
		return failure.Wrap(err, failure.Context{
			"protocol":      req.Proto,
			"host":          req.URL.Hostname(),
			"port":          req.URL.Port(),
			"request_url":   req.URL.RequestURI(),
			"response_body": string(responseBody),
		})
	}
	return nil
}

func escape(req *http.Request) {
	req.URL.RawQuery = url.PathEscape(req.URL.RawQuery)
}
