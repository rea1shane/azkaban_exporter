package api

import (
	"context"
	"fmt"
	"testing"
)

func TestAuthenticate(t *testing.T) {
	params := AuthenticateParams{
		ServerUrl: "http://azkaban:20000",
		Username:  "metrics",
		Password:  "metrics",
	}
	sessionId, err := Authenticate(params, context.Background())
	if err != nil {
		fmt.Printf("%+v", err)
		return
	}
	fmt.Printf("session id: %s\n", sessionId)
}
