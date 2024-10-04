package usergrp

// AppNewUser contains information needed to create a new user.
type AppNewUser struct {
	Name       string `json:"name" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Department string `json:"department"`
	CreditCard string `json:"credit_card" validate:"required"`
}
