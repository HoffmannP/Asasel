package main

import (
	"context"
	"fmt"
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
		LockState bool `path:"lockstate" doc:"Lockstate (true = locked)"`
	}
}

type AccountLockOutput struct {
	Body struct {
		Message   string `json:"message" example:"Hello, world!" doc:"Return message"`
		LockState bool   `json:"lockstate" example:"true" doc:"Lockstate (true = locked)"`
	}
}

type AccountTimeOutput struct {
	Body struct {
		Message  string `json:"message" example:"Hello, world!" doc:"Return message"`
		Duration int    `json:"duration" example:"30" doc:"Login since"`
	}
}

func getLockstate(account string) (bool, error) {
	cmd := exec.Command("chage", "-i", "-l", account)
	cmd.Env = append(cmd.Env, "LC_ALL=en_US")
	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false, err
	}
	exp_str := strings.Trim(strings.Split(strings.Split(out.String(), "\n")[3], ":")[1], " ")
	if exp_str == "never" {
		return false, nil
	}
	exp_date, err := time.Parse("2006-01-02", exp_str)
	if err != nil {
		return false, err
	}
	days := time.Since(exp_date).Hours() / 24
	return days > 1, nil
}

func setLockstate(account string, lockstate bool) error {
	var expiration string

	if lockstate {
		expiration = "1"
		// lockCommand = "-L"
	} else {
		expiration = ""
		// lockCommand = "-U"
	}

	cmd := exec.Command("usermod", "--expiredate", expiration, account)
	return cmd.Run()
}

func getLogintime(account string) (int, error) {
	var out strings.Builder
	cmd := exec.Command("who")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return -1, err
	}

	var firstLogin time.Time
	for _, line := range strings.Split(out.String(), "\n") {
		if strings.HasPrefix(line, account) {
			logintime, err := time.ParseInLocation("2006-01-02 15:04", line[22:38], time.Local)
			if err != nil {
				return 0, err
			}
			if (firstLogin == time.Time{}) || (firstLogin.Compare(logintime) == 1) {
				firstLogin = logintime
			}
		}
	}

	var minutes int
	if (firstLogin == time.Time{}) {
		minutes = -1
	} else {
		duration := time.Since(firstLogin).Round(time.Minute)
		minutes = int(duration.Minutes())
	}
	return minutes, nil
}

func killall(account string) error {
	return exec.Command("killall", "-u", account).Run()
}

func RegisterAccountOperations(api huma.API) {
	huma.Get(api, "/lock/{account}", func(ctx context.Context, input *AccountInput) (*AccountLockOutput, error) {
		resp := &AccountLockOutput{}
		locked, err := getLockstate(input.Account)
		if err != nil {
			resp.Body.Message = "Error getting lockstate"
			return resp, err
		}
		resp.Body.LockState = locked
		if locked {
			resp.Body.Message = "Account locked"
		} else {
			resp.Body.Message = "Account unlocked"
		}
		return resp, nil
	})

	huma.Post(api, "/lock/{account}", func(ctx context.Context, input *AccountLockInput) (*MessageOutput, error) {
		resp := &MessageOutput{}
		err := setLockstate(input.Account, input.Body.LockState)

		var lockaction string
		if input.Body.LockState {
			lockaction = "locked"
		} else {
			lockaction = "unlocked"
		}

		if err != nil {
			resp.Body.Message = "Error setting account to " + lockaction
			return resp, err
		}
		resp.Body.Message = "Account " + lockaction
		return resp, nil
	})

	huma.Get(api, "/time/{account}", func(ctx context.Context, input *AccountInput) (*AccountTimeOutput, error) {
		resp := &AccountTimeOutput{}
		minutes, err := getLogintime(input.Account)
		resp.Body.Duration = minutes
		if err != nil {
			if minutes == -1 {
				resp.Body.Message = "Error getting logintime"
			} else {
				resp.Body.Message = "Timeout unreadable"
			}
			return resp, err
		}
		if minutes == -1 {
			resp.Body.Message = "Account not logged in"
		} else {
			resp.Body.Message = fmt.Sprintf("Account logged in for %d min", minutes)
		}
		return resp, nil
	})

	huma.Get(api, "/killall/{account}", func(ctx context.Context, input *AccountInput) (*MessageOutput, error) {
		resp := &MessageOutput{}
		err := killall(input.Account)
		if err != nil {
			resp.Body.Message = "Error killing all processes"
		} else {
			resp.Body.Message = "All processes killed"
		}
		return resp, err
	})
}
