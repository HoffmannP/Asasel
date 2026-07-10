package main

import "github.com/danielgtaylor/huma/v2"

type MessageOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Return message"`
	}
}

func addRoutes(api huma.API, app *App) {
	RegisterAccountOperations(huma.NewGroup(api, "/accounts"))
	RegisterTimeoutOperations(huma.NewGroup(api, "/timeouts"))
	RegisterConfigOperation(huma.NewGroup(api, "/config"), app.listServers, app.cfg.DefaultAccount)
}
