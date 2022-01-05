package api

import (
	"context"
	"fmt"
	"testing"
)

func TestAuthenticate(t *testing.T) {
	params := AuthenticateParam{
		ServerUrl: "http://172.16.87.150:20000",
		Username:  "metricss",
		Password:  "metrics",
	}
	authenticate, err := Authenticate(params, context.Background())
	if err != nil {
		fmt.Printf("%+v", err)
	}
	fmt.Println(authenticate)
}
