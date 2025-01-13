package usergrp

import (
	"time"

	"github.com/Zanda256/commitsmart-task/business/core/user"
)

// AppNewUser contains information needed to create a new user.
type AppNewUser struct {
	Name       string `json:"name" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Department string `json:"department"`
	CreditCard string `json:"credit_card" validate:"required"`
}

// AppUser represents information about an individual user.
type AppUser struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Department  string `json:"department"`
	CreditCard  string `json:"credit_card"`
	DateCreated string `json:"date_created"`
	DateUpdated string `json:"date_updated"`
}

func toAppUser(usr user.User) AppUser {

	return AppUser{
		ID:          usr.ID.String(),
		Name:        usr.Name,
		CreditCard:  usr.CreditCard,
		Email:       usr.Email.Address,
		Department:  usr.Department,
		DateCreated: usr.DateCreated.Format(time.RFC3339),
		DateUpdated: usr.DateUpdated.Format(time.RFC3339),
	}
}

func toAppUsers(users []user.User) []AppUser {
	items := make([]AppUser, len(users))
	for i, usr := range users {
		items[i] = toAppUser(usr)
	}

	return items
}

// =============================================================================
