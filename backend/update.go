package main

import (
	"context"
	"os/exec"
	"strings"

	"github.com/danielgtaylor/huma/v2"
)

type Noput struct {
}

func initUpdate() error {
	cmd := exec.Command("at", "now", "+", "1", "min")
	cmd.Stdin = strings.NewReader(
		"cd /tmp;" +
			"git clone https://github.com/HoffmannP/Asasel.git;" +
			"cd Asasel;" +
			"task;" +
			"rm -rf /tmp/Asasel")
	return cmd.Run()
}

func RegisterUpdateOperations(api huma.API) {
	huma.Post(api, "/update", func(ctx context.Context, input *Noput) (*Noput, error) {
		err := initUpdate()
		return nil, err
	})
}
