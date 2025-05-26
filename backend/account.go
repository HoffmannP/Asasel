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

type AccountInput struct {
	Account string `path:"account" maxLength:"30" example:"linus" doc:"Account to lock"`
}

type AccountLockInput struct {
	AccountInput
	Body struct {
		LockState bool `path:"lockstatee" doc:"Lockstate (true = locked)"`
	}
}

type AccountTimeOutput struct {
	Body struct {
		Message  string `json:"message" example:"Hello, world!" doc:"Return message"`
		Duration int    `json:"duration" example:"30" doc:"Login since"`
	}
}

func RegisterAccountOperations(api huma.API) {
	huma.Get(api, "/lock/{account}", func(ctx context.Context, input *AccountInput) (*MessageOutput, error) {
		cmd := exec.Command("passwd", "-S", input.Account)
		var out strings.Builder
		cmd.Stdout = &out
		err := cmd.Run()

		resp := &MessageOutput{}
		if err != nil {
			log.Fatal(err)
			resp.Body.Message = fmt.Sprintf("Error getting lockstate for %s", input.Account)
			return resp, err
		}

		lockprefix := "un"
		if strings.Split(out.String(), " ")[1] == "L" {
			lockprefix = ""
		}

		resp.Body.Message = fmt.Sprintf("Account %s %slocked", lockprefix, input.Account)
		return resp, nil
	})

	huma.Post(api, "/lock/{account}", func(ctx context.Context, input *AccountLockInput) (*MessageOutput, error) {
		// lock via lockCommand
		lockCommand := "-e ''"
		// lockCommand := "-U"
		lockprefix := "un"
		if input.Body.LockState {
			// lock via lockCommand
			lockCommand = "-e '1'"
			// lockCommand = "-L"
			lockprefix = ""
		}

		cmd := exec.Command("usermod", lockCommand, input.Account)
		err := cmd.Run()

		resp := &MessageOutput{}
		if err != nil {
			log.Fatal(err)
			resp.Body.Message = fmt.Sprintf("Error %slocking %s", lockprefix, input.Account)
			return resp, err
		}

		resp.Body.Message = fmt.Sprintf("Account %s %slocked", lockprefix, input.Account)
		return resp, nil
	})

	huma.Get(api, "/time/{account}", func(ctx context.Context, input *AccountInput) (*AccountTimeOutput, error) {
		cmd := exec.Command("who")
		var out strings.Builder
		cmd.Stdout = &out
		err := cmd.Run()

		resp := &AccountTimeOutput{}
		if err != nil {
			log.Fatal(err)
			resp.Body.Message = fmt.Sprintf("Error getting time for %s", input.Account)
			return resp, err
		}

		var firstLogin time.Time
		for _, line := range strings.Split(out.String(), "\n") {
			if strings.HasPrefix(line, input.Account) {
				logintime, err := time.ParseInLocation("2006-01-02 15:04", line[22:38], time.Local)
				if err != nil {
					resp.Body.Message = "Timeout unreadable"
					return resp, nil
				}
				if (firstLogin == time.Time{}) || (firstLogin.Compare(logintime) == 1) {
					firstLogin = logintime
				}
			}
		}

		if (firstLogin == time.Time{}) {
			resp.Body.Message = fmt.Sprintf("Account %s not logged in", input.Account)
		} else {
			duration := time.Since(firstLogin).Round(time.Minute)
			resp.Body.Message = fmt.Sprintf("Account logged in since %s", strings.TrimSuffix(duration.String(), "0s"))
			resp.Body.Duration = int(duration.Minutes())
		}
		return resp, nil
	})

	huma.Get(api, "/killall/{account}", func(ctx context.Context, input *AccountInput) (*MessageOutput, error) {
		resp := &MessageOutput{}
		cmd := exec.Command("killall", "-u", input.Account)
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
			resp.Body.Message = fmt.Sprintf("Error killing all processes of %s", input.Account)
			return resp, err
		}
		resp.Body.Message = fmt.Sprintf("All processes of %s killed", input.Account)
		return resp, nil
	})
}
