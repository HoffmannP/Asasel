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

func RegisterConfigOperation(api huma.API) {
	huma.Get(api, "/servers", func(ctx context.Context, input *NoInput) (*ServlistOutput, error) {
		resp := &ServlistOutput{}
		resp.Body.Servers = []string{
			"yoga",
			"nura",
			"hackmack",
			"bereisen",
		}
		return resp, nil
	})
}
