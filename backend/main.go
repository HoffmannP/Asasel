package main

import (
	"fmt"
	"net/http"

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

func addRoutes(api huma.API) {
	RegisterAccountOperations(huma.NewGroup(api, "/accounts"))
	RegisterTimeoutOperations(huma.NewGroup(api, "/timeouts"))

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

		router.Route("/api", func(api_route chi.Router) {
			api := humachi.New(api_route, huma.DefaultConfig("My API", "1.0.0"))
			addRoutes(api)
		})

		router.Get("/*", http.FileServer(http.Dir("/static")).ServeHTTP)

		hooks.OnStart(func() {
			fmt.Printf("Starting server on port %d...\n", options.Port)
			http.ListenAndServe(fmt.Sprintf(":%d", options.Port), router)
		})
	})

	// Run the CLI. When passed no commands, it starts the server.
	cli.Run()
}
