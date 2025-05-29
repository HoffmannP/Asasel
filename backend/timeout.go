package main

import (
	"context"
	"fmt"
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

func setTimeout(duration int, account string) error {
	cmd := exec.Command("at", "now", "+", fmt.Sprintf("%d", duration), "min")
	cmd.Stdin = strings.NewReader(fmt.Sprintf(
		"#timeout %s\nkillall -u '%s'\nusermod -e 1 '%s'",
		account, account, account))
	return cmd.Run()
}

func getTimeout(account string) (int, error) {
	tag := fmt.Sprintf("#timeout %s", account)
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
		err := cmd.Run()
		if err != nil {
			return 0, err
		}
		lines := strings.Split(out.String(), "\n")
		if lines[len(lines)-4] == tag {
			timestamp := strings.Split(line, "\t")[1][0:24]
			exectime, err := time.ParseInLocation(time.ANSIC, timestamp, time.Local)
			if err != nil {
				return 0, err
			}
			duration := time.Until(exectime).Round(time.Minute)
			return int(duration.Minutes()), nil
		}
	}
	return -1, nil
}

func delTimeout(account string) (bool, error) {
	tag := fmt.Sprintf("#timeout %s", account)
	cmd := exec.Command("atq")
	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false, err
	}
	for _, line := range strings.Split(out.String(), "\n") {
		if line == "" {
			continue
		}
		atid := strings.Split(line, "\t")[0]
		cmd := exec.Command("at", "-c", atid)
		var out strings.Builder
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			return false, err
		}
		lines := strings.Split(out.String(), "\n")
		if lines[len(lines)-4] == tag {
			err := exec.Command("atrm", atid).Run()
			if err != nil {
				return false, err
			}
			return true, nil
		}
	}
	return false, nil
}

func RegisterTimeoutOperations(api huma.API) {
	huma.Post(api, "/{account}", func(ctx context.Context, input *TimeoutInput) (*TimeoutOutput, error) {
		resp := &TimeoutOutput{}
		err := setTimeout(input.Body.Duration, input.Account)
		if err != nil {
			resp.Body.Message = "Error setting timeout"
			return resp, err
		}
		resp.Body.Message = fmt.Sprintf("Timeout is %d", input.Body.Duration)
		resp.Body.Remaining = input.Body.Duration
		return resp, nil
	})

	huma.Get(api, "/{account}", func(ctx context.Context, input *AccountInput) (*TimeoutOutput, error) {
		resp := &TimeoutOutput{}
		remaining, err := getTimeout(input.Account)
		if err != nil {
			resp.Body.Message = "Timeout unreadable"
			return resp, nil
		}
		resp.Body.Remaining = remaining
		if remaining == -1 {
			resp.Body.Message = "Timeout not found"
		} else {
			resp.Body.Message = fmt.Sprintf("Timeout is %d", remaining)
		}
		return resp, nil
	})

	huma.Delete(api, "/{account}", func(ctx context.Context, input *AccountInput) (*MessageOutput, error) {
		resp := &MessageOutput{}
		success, err := delTimeout(input.Account)
		if err != nil {
			resp.Body.Message = "Timeout not deleted"
			return resp, err
		}
		if success {
			resp.Body.Message = "Timeout unset"
		} else {
			resp.Body.Message = "Timeout not found"
		}
		return resp, nil
	})
}
