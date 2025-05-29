package main

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

//go:embed all:static/*
var gui embed.FS

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
}

func FileServer(w http.ResponseWriter, r *http.Request) {
	path := "static" + r.URL.Path

	if path[len(path)-1] == '/' {
		path += "index.html"
	}

	println(path)

	http.ServeFileFS(w, r, gui, path)
}

func main() {
	// Create a CLI app which takes a port option.
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		// Create a new router & API
		router := chi.NewMux()
		router.Use(middleware.Logger)

		router.Route("/api", func(api_route chi.Router) {
			api := humachi.New(api_route, huma.DefaultConfig("My API", "1.0.0"))
			addRoutes(api)
		})

		router.Get("/*", FileServer)

		hooks.OnStart(func() {
			fmt.Printf("Starting server on port %d...\n", options.Port)
			http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", options.Port), router)
		})
	})

	// Run the CLI. When passed no commands, it starts the server.
	cli.Run()
}
