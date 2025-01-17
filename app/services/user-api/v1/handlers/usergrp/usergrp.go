package usergrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/mail"

	"github.com/Zanda256/commitsmart-task/business/web/v1/response"

	"github.com/Zanda256/commitsmart-task/business/core/user"
	"github.com/Zanda256/commitsmart-task/foundation/validate"
	"github.com/Zanda256/commitsmart-task/foundation/web"
)

// Handlers manages the set of user endpoints.
type Handlers struct {
	user *user.Core
}

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
		"Created":              toAppUser(usr),
	}
	return web.Respond(ctx, w, ret, http.StatusCreated)
}

func (h *Handlers) Query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	filter, err := parseFilter(ctx, r)
	if err != nil {
		return err
	}

	users, err := h.user.Query(ctx, filter)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	return web.Respond(ctx, w, response.NewPageDocument(toAppUsers(users), len(users), 1, 10), http.StatusOK)
}

func (h *Handlers) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	filter, err := parseFilter(ctx, r)
	if err != nil {
		return err
	}

	if filter.UserID == nil {
		return response.NewError(validate.NewFieldsError("user_id", errors.New("user_id field is required")), http.StatusBadRequest)
	}

	usr, err := h.user.QueryByID(ctx, filter)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	return web.Respond(ctx, w, response.NewPageDocument(toAppUsers([]user.User{usr}), 1, 1, 10), http.StatusOK)
}
