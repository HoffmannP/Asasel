package main

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

type NoInput struct {
}

type ServlistOutput struct {
	Body struct {
		Servers []string `json:"server" example:"List of servers"`
	}
}

type ClientConfigOutput struct {
	Body struct {
		Account string `json:"account" example:"linus"`
	}
}

func RegisterConfigOperation(api huma.API, serverProvider func() []string, defaultAccount string) {
	huma.Get(api, "/servers", func(ctx context.Context, input *NoInput) (*ServlistOutput, error) {
		resp := &ServlistOutput{}
		resp.Body.Servers = serverProvider()
		return resp, nil
	})

	huma.Get(api, "/client", func(ctx context.Context, input *NoInput) (*ClientConfigOutput, error) {
		resp := &ClientConfigOutput{}
		resp.Body.Account = defaultAccount
		return resp, nil
	})
}
