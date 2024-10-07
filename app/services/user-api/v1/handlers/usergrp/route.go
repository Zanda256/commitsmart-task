package usergrp

import (
	"net/http"

	"github.com/Zanda256/commitsmart-task/business/core/user"
	"github.com/Zanda256/commitsmart-task/business/core/user/stores/userdb"
	documentStore "github.com/Zanda256/commitsmart-task/business/data/docStore"
	"github.com/Zanda256/commitsmart-task/foundation/logger"
	"github.com/Zanda256/commitsmart-task/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build        string
	Store        *documentStore.DocStorage
	Log          *logger.Logger
	UserDbName   string
	UserCollName string
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"
	db := userdb.NewStore(cfg.Log, cfg.UserDbName, cfg.Store, cfg.UserCollName)
	hdl := New(user.NewCore(cfg.Log, db))
	app.HandlePath(http.MethodPost, version, "/users", hdl.Create)
	app.HandlePath(http.MethodGet, version, "/users/:user_id", hdl.Query)
}
