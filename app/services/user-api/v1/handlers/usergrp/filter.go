package usergrp

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/mail"

	"github.com/google/uuid"

	"github.com/Zanda256/commitsmart-task/business/core/user"
	"github.com/Zanda256/commitsmart-task/foundation/validate"
	"github.com/Zanda256/commitsmart-task/foundation/web"
)

func parseFilter(ctx context.Context, r *http.Request) (user.QueryFilter, error) {
	const (
		filterByUserID = "user_id"
		filterByEmail  = "email"
		filterByName   = "name"
	)

	values := r.URL.Query()

	var filter user.QueryFilter

	p := web.GetPathParams(ctx)
	pars, ok := p.(*httprouter.Params)
	if ok {
		userID := pars.ByName("user_id")
		if userID != "" {
			id, err := uuid.Parse(userID)
			if err != nil {
				return user.QueryFilter{}, validate.NewFieldsError(filterByUserID, err)
			}
			filter.WithUserID(id)
		}
	}

	if email := values.Get(filterByEmail); email != "" {
		addr, err := mail.ParseAddress(email)
		if err != nil {
			return user.QueryFilter{}, validate.NewFieldsError(filterByEmail, err)
		}
		filter.WithEmail(*addr)
	}

	if name := values.Get(filterByName); name != "" {
		filter.WithName(name)
	}

	if err := filter.Validate(); err != nil {
		return user.QueryFilter{}, err
	}

	return filter, nil
}
