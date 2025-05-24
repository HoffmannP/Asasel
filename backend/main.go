package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"

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

type GreetingOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
	}
}

type ReviewInput struct {
	Body struct {
		Author  string `json:"author" maxLength:"10" doc:"Author of the review"`
		Rating  int    `json:"rating" minimum:"1" maximum:"5" doc:"Rating from 1 to 5"`
		Message string `json:"message,omitempty" maxLength:"100" doc:"Review message"`
	}
}

type AccountOutput struct {
	Body struct {
		Message string `json:"message" example:"Locked" doc:"Account status message"`
	}
}

func addRoutes(api huma.API) {
	// Register GET /greeting/{name}
	huma.Register(api, huma.Operation{
		OperationID: "get-greeting",
		Method:      http.MethodGet,
		Path:        "/greeting/{name}",
		Summary:     "Greet a greeting",
		Description: "Get a greeting for a person by name.",
		Tags:        []string{"Greetings"},
	}, func(ctx context.Context, input *struct {
		Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
	}) (*GreetingOutput, error) {
		resp := &GreetingOutput{}
		resp.Body.Message = fmt.Sprintf("Hello %s!", input.Name)
		return resp, nil
	})

	// Register POST /reviews handler.
	huma.Register(api, huma.Operation{
		OperationID:   "post-review",
		Method:        http.MethodPost,
		Path:          "/reviews",
		Summary:       "Gret a greeting",
		Description:   "Post a review.",
		Tags:          []string{"Reviews"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, i *ReviewInput) (*struct{}, error) {
		// TODO: save review in data store
		return nil, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "lock-acocunt",
		Method:      http.MethodGet,
		Path:        "/account/lock/{account}",
		Summary:     "Locks account",
		Description: "Locks the account of an existing user",
		Tags:        []string{"Account"},
	}, func(ctx context.Context, input *struct {
		Account string `path:"account" maxLength:"30" example:"linus" doc:"Account to lock"`
	}) (*AccountOutput, error) {
		resp := &AccountOutput{}
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
		OperationID: "unlock-acocunt",
		Method:      http.MethodGet,
		Path:        "/account/unlock/{account}",
		Summary:     "Unlocks account",
		Description: "Unlocks the account of an existing user",
		Tags:        []string{"Account"},
	}, func(ctx context.Context, input *struct {
		Account string `path:"account" maxLength:"30" example:"linus" doc:"Account to unlock"`
	}) (*AccountOutput, error) {
		resp := &AccountOutput{}
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
		Summary:     "Kills all processes of account",
		Description: "Kills every single process of an existing account",
		Tags:        []string{"Account"},
	}, func(ctx context.Context, input *struct {
		Account string `path:"account" maxLength:"30" example:"linus" doc:"Account to kill"`
	}) (*AccountOutput, error) {
		resp := &AccountOutput{}
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
		OperationID: "account-timeout",
		Method:      http.MethodGet,
		Path:        "/account/timeout/{account}/{duration}",
		Summary:     "Timeout accounts session",
		Description: "Kills every single process of an existing account after some time",
		Tags:        []string{"Account"},
	}, func(ctx context.Context, input *struct {
		Account string `path:"account" maxLength:"30" example:"linus" doc:"Account to kill"`
		Minutes int    `path:"duration" maxLength:"30" example:"30" doc:"Time to wait"`
	}) (*AccountOutput, error) {
		resp := &AccountOutput{}
		cmd := exec.Command("at", "now", "+", fmt.Sprintf("%d", input.Minutes), "min")
		cmd.Stdin = strings.NewReader(fmt.Sprintf("killall -u '%s'; usermod -e 1 '%s'", input.Account, input.Account))
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
			resp.Body.Message = "Error setting timeout"
			return resp, err
		}
		resp.Body.Message = "Timeout set"
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
