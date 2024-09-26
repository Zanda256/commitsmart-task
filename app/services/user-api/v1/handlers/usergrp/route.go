package usergrp

import (
	"github.com/Zanda256/commitsmart-task/foundation/logger"
	"github.com/Zanda256/commitsmart-task/foundation/web"
	"net/http"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build string
	Log   *logger.Logger
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	hdl := New()
	app.HandlePath(http.MethodPost, version, "/users", hdl.Create)
}
