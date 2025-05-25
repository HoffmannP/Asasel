package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

// Options for the CLI
type Options struct {
	Port int `help:"Port to listen on" short:"p" default:"8888"`
}

// Call saved command

type MessageOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Return message"`
	}
}

type TimeoutOutput struct {
	Body struct {
		Message   string `json:"message" example:"Hello, world!" doc:"Return message"`
		Remaining int    `json:"remaining" example:"30" doc:"Remaining minutes"`
	}
}

type AccountInput struct {
	Account string `path:"account" maxLength:"30" example:"linus" doc:"Account to lock"`
}

type TimeoutInput struct {
	Account string `path:"account" maxLength:"30" example:"linus" doc:"Account to kill"`
	Body    struct {
		Duration int `path:"duration" maxLength:"30" example:"30" doc:"Time to wait"`
	}
}

func addRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "account-lock",
		Method:      http.MethodGet,
		Path:        "/account/lock/{account}",
		Summary:     "Lock account",
		Description: "Locks the account of an existing user",
		Tags:        []string{"Account"},
	}, func(ctx context.Context, input *AccountInput) (*MessageOutput, error) {
		resp := &MessageOutput{}
		cmd := exec.Command("usermod", "-e", "1", input.Account)
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
			resp.Body.Message = fmt.Sprintf("Error locking account %s", input.Account)
			return resp, err
		}
		resp.Body.Message = fmt.Sprintf("Account %s locked", input.Account)
		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "account-unlock",
		Method:      http.MethodGet,
		Path:        "/account/unlock/{account}",
		Summary:     "Unlock account",
		Description: "Unlocks the account of an existing user",
		Tags:        []string{"Account"},
	}, func(ctx context.Context, input *AccountInput) (*MessageOutput, error) {
		resp := &MessageOutput{}
		cmd := exec.Command("usermod", "-e", "", input.Account)
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
			resp.Body.Message = fmt.Sprintf("Error unlocking account %s", input.Account)
			return resp, err
		}
		resp.Body.Message = fmt.Sprintf("Account %s unlocked", input.Account)
		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "account-killall",
		Method:      http.MethodGet,
		Path:        "/account/killall/{account}",
		Summary:     "Kill all processes",
		Description: "Kills all processes of an existing account",
		Tags:        []string{"Account"},
	}, func(ctx context.Context, input *AccountInput) (*MessageOutput, error) {
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

	huma.Register(api, huma.Operation{
		OperationID: "account-settimeout",
		Method:      http.MethodPost,
		Path:        "/timeout/{account}",
		Summary:     "Set session timeout",
		Description: "Sets a session timeout for the specified account",
		Tags:        []string{"Timeout"},
	}, func(ctx context.Context, input *TimeoutInput) (*TimeoutOutput, error) {
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

	huma.Register(api, huma.Operation{
		OperationID: "account-gettimeout",
		Method:      http.MethodGet,
		Path:        "/timeout/{account}",
		Summary:     "Get session timeout",
		Description: "Gets the remaining duration until the next session timeout for the specified account",
		Tags:        []string{"Timeout"},
	}, func(ctx context.Context, input *AccountInput) (*TimeoutOutput, error) {
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

	huma.Register(api, huma.Operation{
		OperationID: "account-rmtimeout",
		Method:      http.MethodDelete,
		Path:        "/timeout/{account}",
		Summary:     "Remove session timeout",
		Description: "Removes the first found session timeout for the specified account",
		Tags:        []string{"Timeout"},
	}, func(ctx context.Context, input *AccountInput) (*MessageOutput, error) {
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

	/* For later for meta */
	// Add new Command

	// Update Command

	// Remove Command

	// Call Command
}

func main() {
	// Create a CLI app which takes a port option.
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		// Create a new router & API
		router := chi.NewMux()
		api := humachi.New(router, huma.DefaultConfig("My API", "1.0.0"))

		addRoutes(api)

		hooks.OnStart(func() {
			fmt.Printf("Starting server on port %d...\n", options.Port)
			http.ListenAndServe(fmt.Sprintf(":%d", options.Port), router)
		})
	})

	// Run the CLI. When passed no commands, it starts the server.
	cli.Run()
}
