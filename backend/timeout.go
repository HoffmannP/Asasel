package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

type TimeoutOutput struct {
	Body struct {
		Message   string `json:"message" example:"Hello, world!" doc:"Return message"`
		Remaining int    `json:"remaining" example:"30" doc:"Remaining minutes"`
	}
}

type TimeoutInput struct {
	AccountInput
	Body struct {
		Duration int `path:"duration" maxLength:"30" example:"30" doc:"Time to wait"`
	}
}

func RegisterTimeoutOperations(api huma.API) {
	huma.Post(api, "/{account}", func(ctx context.Context, input *TimeoutInput) (*TimeoutOutput, error) {
		resp := &TimeoutOutput{}
		cmd := exec.Command("at", "now", "+", fmt.Sprintf("%d", input.Body.Duration), "min")
		cmd.Stdin = strings.NewReader(fmt.Sprintf(
			"#timeout %s\nkillall -u '%s'\nusermod -e 1 '%s'",
			input.Account,
			input.Account,
			input.Account))
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
			resp.Body.Message = "Error setting timeout"
			return resp, err
		}
		resp.Body.Message = "Timeout set"
		resp.Body.Remaining = input.Body.Duration
		return resp, nil
	})

	huma.Get(api, "/{account}", func(ctx context.Context, input *AccountInput) (*TimeoutOutput, error) {
		resp := &TimeoutOutput{}
		tag := fmt.Sprintf("#timeout %s", input.Account)
		cmd := exec.Command("atq")
		var out strings.Builder
		cmd.Stdout = &out
		cmd.Run()
		for _, line := range strings.Split(out.String(), "\n") {
			if line == "" {
				continue
			}
			atid := strings.Split(line, "\t")[0]
			cmd := exec.Command("at", "-c", atid)
			var out strings.Builder
			cmd.Stdout = &out
			cmd.Run()
			lines := strings.Split(out.String(), "\n")
			if lines[len(lines)-4] == tag {
				timestamp := strings.Split(line, "\t")[1][0:24]
				exectime, err := time.ParseInLocation(time.ANSIC, timestamp, time.Local)
				if err != nil {
					resp.Body.Message = "Timeout unreadable"
					return resp, nil
				}
				duration := time.Until(exectime).Round(time.Minute)
				resp.Body.Message = fmt.Sprintf("Timeout is %s", strings.TrimSuffix(duration.String(), "0s"))
				resp.Body.Remaining = int(duration.Minutes())
				return resp, nil
			}
		}
		resp.Body.Message = "Timeout not found"
		return resp, nil
	})

	huma.Delete(api, "/{account}", func(ctx context.Context, input *AccountInput) (*MessageOutput, error) {
		resp := &MessageOutput{}
		tag := fmt.Sprintf("#timeout %s", input.Account)
		cmd := exec.Command("atq")
		var out strings.Builder
		cmd.Stdout = &out
		cmd.Run()
		for _, line := range strings.Split(out.String(), "\n") {
			if line == "" {
				continue
			}
			atid := strings.Split(line, "\t")[0]
			cmd := exec.Command("at", "-c", atid)
			var out strings.Builder
			cmd.Stdout = &out
			cmd.Run()
			lines := strings.Split(out.String(), "\n")
			if lines[len(lines)-4] == tag {
				exec.Command("atrm", atid).Run()
				resp.Body.Message = "Timeout unset"
				return resp, nil
			}
		}
		resp.Body.Message = "Timeout not found"
		return resp, nil
	})
}
