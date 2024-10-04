package user

import (
	"github.com/google/uuid"
	"net/mail"
	"time"
)

// User represents information about an individual user.
type User struct {
	ID          uuid.UUID
	Name        string
	Email       mail.Address
	Department  string
	CreditCard  string
	DateCreated time.Time
	DateUpdated time.Time
}

// NewUser contains information needed to create a new user.
type NewUser struct {
	Name       string       `json:"name"`
	Email      mail.Address `json:"email"`
	Department string       `json:"department"`
	CreditCard string       `json:"credit_card"`
}
