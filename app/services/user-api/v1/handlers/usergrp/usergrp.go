package usergrp

import (
	"context"
	"fmt"
	"net/http"
	"net/mail"

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
	v := &AppNewUser{}
	err := web.Decode(r, v)
	if err != nil {
		return fmt.Errorf("web.Decode: %w", err)
	}

	addr, err := mail.ParseAddress(v.Email)
	if err != nil {
		return fmt.Errorf("parsing email: %w", err)
	}

	nu := user.NewUser{
		Name:       v.Name,
		Email:      *addr,
		Department: v.Department,
		CreditCard: v.CreditCard,
	}

	usr, err := h.user.Create(ctx, nu)
	if err != nil {
		return err
	}
	ret := map[string]any{
		"User Create endpoint": "success",
		"Created":              usr,
	}
	return web.Respond(ctx, w, ret, http.StatusCreated)
}
