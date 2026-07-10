package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

func main() {
	cfg := parseConfig()
	app := newApp(cfg)

	if cfg.Mode != "local" && cfg.Mode != "control" {
		log.Fatalf("invalid mode %q, expected local|control", cfg.Mode)
	}
	if (cfg.AuthUser == "") != (cfg.AuthPass == "") {
		log.Fatalf("both -auth-user and -auth-pass must be set together")
	}
	if cfg.DefaultAccount == "" {
		log.Fatalf("-account must not be empty")
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(func(next http.Handler) http.Handler {
		return app.requireBasicAuth(next)
	})
	app.registerControlRoutes(router)

	router.Route("/api", func(api_route chi.Router) {
		api := humachi.New(api_route, huma.DefaultConfig("My API", "1.0.0"))
		addRoutes(api, app)
	})

	router.Get("/*", FileServer)

	if cfg.Mode == "local" && cfg.ControllerURL != "" {
		go app.runAgentLoop(context.Background())
	}

	fmt.Printf("Starting %s server on %s...\n", cfg.Mode, cfg.ListenAddr)
	http.ListenAndServe(cfg.ListenAddr, router)
}
