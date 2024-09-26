package usergrp

import (
	"context"
	"github.com/Zanda256/commitsmart-task/foundation/web"
	"net/http"
)

// Handlers manages the set of user endpoints.
type Handlers struct {
	//user *user.Core
	//auth *auth.Auth
}

// New constructs a handlers for route access.
//func New(user *user.Core, auth *auth.Auth) *Handlers {
//	return &Handlers{
//		user: user,
//		auth: auth,
//	}
//}

//  cd ./app/services/ratesApi && \
//    go build -ldflags "-X main.build=$(BUILD_REF)

func New() *Handlers {
	return &Handlers{}
}

// Create adds a new user to the system.
func (h *Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	//var app AppNewUser
	//if err := web.Decode(r, &app); err != nil {
	//	return response.NewError(err, http.StatusBadRequest)
	//}
	//
	//nc, err := toCoreNewUser(app)
	//if err != nil {
	//	return response.NewError(err, http.StatusBadRequest)
	//}
	//
	//usr, err := h.user.Create(ctx, nc)
	//if err != nil {
	//	if errors.Is(err, user.ErrUniqueEmail) {
	//		return response.NewError(err, http.StatusConflict)
	//	}
	//	return fmt.Errorf("create: usr[%+v]: %w", usr, err)
	//}
	//
	//return web.Respond(ctx, w, toAppUser(usr), http.StatusCreated)

	return web.Respond(ctx, w, "User Create endpoint", http.StatusCreated)
}
