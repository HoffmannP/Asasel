package main

import (
	"context"
	"errors"
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
		LockState bool `json:"lockstate" doc:"Lockstate (true = locked)"`
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

type AccountStateOutput struct {
	Body struct {
		Message   string `json:"message" example:"ok" doc:"Combined state info"`
		LockState bool   `json:"lockstate" example:"true" doc:"Lockstate (true = locked)"`
		Duration  int    `json:"duration" example:"30" doc:"Login since"`
		Remaining int    `json:"remaining" example:"30" doc:"Timeout remaining"`
	}
}

func getLockstate(account string) (bool, error) {
	cmd := exec.Command("chage", "-i", "-l", account)
	cmd.Env = append(cmd.Env, "LC_ALL=C")
	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false, err
	}

	var expStr string
	for _, line := range strings.Split(out.String(), "\n") {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(parts[0]), "Account expires") {
			expStr = strings.TrimSpace(parts[1])
			break
		}
	}
	if expStr == "" {
		return false, errors.New("unable to parse account expiration from chage output")
	}
	if expStr == "never" {
		return false, nil
	}
	exp_date, err := time.Parse("2006-01-02", expStr)
	if err != nil {
		return false, err
	}
	days := time.Since(exp_date).Hours() / 24
	return days > 1, nil
}

func setLockstate(account string, lockstate bool) error {
	var expiration string
	if lockstate {
		expiration = "1970-01-02"
	} else {
		expiration = "-1"
	}

	cmd := exec.Command("chage", "-E", expiration, account)
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
		fields := strings.Fields(line)
		if len(fields) < 4 || fields[0] != account {
			continue
		}

		logintime, err := time.ParseInLocation("2006-01-02 15:04", fields[2]+" "+fields[3], time.Local)
		if err != nil {
			return 0, err
		}
		if (firstLogin == time.Time{}) || (firstLogin.Compare(logintime) == 1) {
			firstLogin = logintime
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
	if err := exec.Command("loginctl", "terminate-user", account).Run(); err == nil {
		return nil
	}
	return exec.Command("killall", "-u", account).Run()
}

func RegisterAccountOperations(api huma.API, app *App) {
	huma.Get(api, "/lock/{account}", func(ctx context.Context, input *AccountInput) (*AccountLockOutput, error) {
		resp := &AccountLockOutput{}
		if app.cfg.Mode == "control" {
			result, err := app.controlForward(RemoteCommand{Op: "get_lock", Account: input.Account})
			if err != nil {
				resp.Body.Message = "Error getting lockstate"
				return resp, err
			}
			resp.Body.LockState = result.LockState != nil && *result.LockState
			resp.Body.Message = result.Message
			return resp, nil
		}

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
		if app.cfg.Mode == "control" {
			result, err := app.controlForward(RemoteCommand{Op: "set_lock", Account: input.Account, LockState: &input.Body.LockState})
			if err != nil {
				resp.Body.Message = "Error setting account to locked: " + err.Error()
				if !input.Body.LockState {
					resp.Body.Message = "Error setting account to unlocked: " + err.Error()
				}
				return resp, err
			}
			resp.Body.Message = result.Message
			return resp, nil
		}

		err := setLockstate(input.Account, input.Body.LockState)

		var lockaction string
		if input.Body.LockState {
			lockaction = "locked"
		} else {
			lockaction = "unlocked"
		}

		if err != nil {
			resp.Body.Message = "Error setting account to " + lockaction + ": " + err.Error()
			return resp, nil
		}
		resp.Body.Message = "Account " + lockaction
		return resp, nil
	})

	huma.Get(api, "/time/{account}", func(ctx context.Context, input *AccountInput) (*AccountTimeOutput, error) {
		resp := &AccountTimeOutput{}
		if app.cfg.Mode == "control" {
			result, err := app.controlForward(RemoteCommand{Op: "get_time", Account: input.Account})
			if err != nil {
				resp.Body.Message = "Error getting logintime"
				return resp, err
			}
			resp.Body.Message = result.Message
			if result.Duration != nil {
				resp.Body.Duration = *result.Duration
			}
			return resp, nil
		}

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

	huma.Get(api, "/state/{account}", func(ctx context.Context, input *AccountInput) (*AccountStateOutput, error) {
		resp := &AccountStateOutput{}
		if app.cfg.Mode == "control" {
			result, err := app.controlForward(RemoteCommand{Op: "get_state", Account: input.Account})
			if err != nil {
				resp.Body.Message = "Error reading lockstate and logintime"
				return resp, err
			}
			resp.Body.Message = result.Message
			if result.LockState != nil {
				resp.Body.LockState = *result.LockState
			}
			if result.Duration != nil {
				resp.Body.Duration = *result.Duration
			}
			if result.Remaining != nil {
				resp.Body.Remaining = *result.Remaining
			} else {
				resp.Body.Remaining = -1
			}
			return resp, nil
		}

		locked, lockErr := getLockstate(input.Account)
		minutes, timeErr := getLogintime(input.Account)
		remaining, timeoutErr := getTimeout(input.Account)

		resp.Body.LockState = locked
		resp.Body.Duration = minutes
		resp.Body.Remaining = remaining

		if lockErr != nil || timeErr != nil || timeoutErr != nil {
			parts := make([]string, 0, 2)
			if lockErr != nil {
				parts = append(parts, "lockstate unavailable")
			}
			if timeErr != nil {
				parts = append(parts, "logintime unavailable")
			}
			if timeoutErr != nil {
				parts = append(parts, "timeout unavailable")
			}
			resp.Body.Message = strings.Join(parts, "; ")
			return resp, nil
		}

		if minutes == -1 {
			resp.Body.Message = "Account not logged in"
		} else {
			resp.Body.Message = fmt.Sprintf("Account logged in for %d min", minutes)
		}
		return resp, nil
	})



	huma.Post(api, "/killall/{account}", func(ctx context.Context, input *AccountInput) (*MessageOutput, error) {
		resp := &MessageOutput{}
		if app.cfg.Mode == "control" {
			result, err := app.controlForward(RemoteCommand{Op: "killall", Account: input.Account})
			if err != nil {
				resp.Body.Message = "Error killing all processes"
				return resp, err
			}
			resp.Body.Message = result.Message
			return resp, nil
		}

		err := killall(input.Account)
		if err != nil {
			resp.Body.Message = "Error killing all processes"
		} else {
			resp.Body.Message = "All processes killed"
		}
		return resp, err
	})
}
