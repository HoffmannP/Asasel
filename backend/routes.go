package main

import "github.com/danielgtaylor/huma/v2"

type MessageOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Return message"`
	}
}

func addRoutes(api huma.API, app *App) {
	RegisterAccountOperations(huma.NewGroup(api, "/accounts"), app)
	RegisterTimeoutOperations(huma.NewGroup(api, "/timeouts"), app)
	RegisterConfigOperation(huma.NewGroup(api, "/config"), app.listServers, app.cfg.DefaultAccount, app.cfg.Mode, app.cfg.ControllerURL)
}
