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

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(func(next http.Handler) http.Handler {
		return app.requireBasicAuth(next)
	})
	app.registerControlRoutes(router)
	app.registerPublicRoutes(router)

	router.Route("/api", func(api_route chi.Router) {
		api := humachi.New(api_route, huma.DefaultConfig("My API", "1.0.0"))
		addRoutes(api, app)
	})

	router.Get("/*", FileServer)

	if cfg.Mode == LocalMode && cfg.ControllerURL != "" {
		go app.runAgentLoop(context.Background())
	}

	if len(cfg.Certs) == 2 {
		fmt.Printf("Starting %s server on %s (TLS)...\n", cfg.Mode, cfg.ListenAddr)
		log.Fatal(http.ListenAndServeTLS(cfg.ListenAddr, cfg.Certs[0], cfg.Certs[1], router))
	} else {
		fmt.Printf("Starting %s server on %s (insecure)...\n", cfg.Mode, cfg.ListenAddr)
		log.Fatal(http.ListenAndServe(cfg.ListenAddr, router))
	}
}
