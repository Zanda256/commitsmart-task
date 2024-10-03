package handlers

import (
	"github.com/Zanda256/commitsmart-task/app/services/user-api/v1/handlers/usergrp"
	v1 "github.com/Zanda256/commitsmart-task/business/web/v1"
	"github.com/Zanda256/commitsmart-task/foundation/web"
)

type Routes struct{}

// Add implements the RouterAdder interface.
func (Routes) Add(app *web.App, cfg v1.APIMuxConfig) {
	usergrp.Routes(app, usergrp.Config{
		Build:        cfg.Build,
		Store:        cfg.DbClients,
		Log:          cfg.Log,
		UserDb:       cfg.UserDb,
		UserCollName: cfg.UserCollectionName,
	})
}
