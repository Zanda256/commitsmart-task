package user

import (
	"fmt"
	"github.com/google/uuid"
	"net/mail"

	"github.com/Zanda256/commitsmart-task/foundation/validate"
)

// QueryFilter holds the available fields a query can be filtered on.
type QueryFilter struct {
	UserID     *uuid.UUID    `validate:"required"`
	Name       *string       `validate:"omitempty"`
	Email      *mail.Address `validate:"omitempty"`
	Department *string       `validate:"omitempty"`
}

// Validate checks the data in the model is considered clean.
func (qf *QueryFilter) Validate() error {
	if err := validate.Check(qf); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}

// WithUserID sets the UserID field of the QueryFilter value.
func (qf *QueryFilter) WithUserID(usrID uuid.UUID) {
	qf.UserID = &usrID
}

// WithName sets the Name field of the QueryFilter value.
func (qf *QueryFilter) WithName(name string) {
	qf.Name = &name
}

// WithEmail sets the Email field of the QueryFilter value.
func (qf *QueryFilter) WithEmail(email mail.Address) {
	qf.Email = &email
}
