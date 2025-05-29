package main

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

//go:embed all:static/*
var gui embed.FS

const PORT = 2727

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
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Route("/api", func(api_route chi.Router) {
		api := humachi.New(api_route, huma.DefaultConfig("My API", "1.0.0"))
		addRoutes(api)
	})

	router.Get("/*", FileServer)

	fmt.Printf("Starting server on port %d...\n", PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), router)
}
