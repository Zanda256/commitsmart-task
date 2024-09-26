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
	DateCreated time.Time
	DateUpdated time.Time
}

// NewUser contains information needed to create a new user.
type NewUser struct {
	Name       string
	Email      mail.Address
	Department string
	CreditCard string
}
