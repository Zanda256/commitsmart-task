package usergrp

import (
	"context"
	"net/http"

	"github.com/Zanda256/commitsmart-task/business/core/user"
	"github.com/Zanda256/commitsmart-task/foundation/web"
)

// Handlers manages the set of user endpoints.
type Handlers struct {
	user *user.Core
}

// New constructs a handlers for route access.
//func New(user *user.Core, auth *auth.Auth) *Handlers {
//	return &Handlers{
//		user: user,
//		auth: auth,
//	}
//}

func New(user *user.Core) *Handlers {
	return &Handlers{
		user: user,
	}
}

// Create adds a new user to the system.
func (h *Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	return web.Respond(ctx, w, "User Create endpoint", http.StatusCreated)
}
