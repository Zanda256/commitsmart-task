package user

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"

	"github.com/Zanda256/commitsmart-task/foundation/logger"
)

// Storer interface declares the behavior this package needs to perists and
// retrieve data.
type Storer interface {
	Create(ctx context.Context, usr User) error
	Query(ctx context.Context, filter QueryFilter, pageNumber int, rowsPerPage int) ([]User, error)
}

// =============================================================================

// Core manages the set of APIs for user access.
type Core struct {
	storer Storer
	log    *logger.Logger
}

// NewCore constructs a core for user api access.
func NewCore(log *logger.Logger, storer Storer) *Core {
	return &Core{
		storer: storer,
		log:    log,
	}
}

func (uc *Core) Create(ctx context.Context, nu NewUser) (User, error) {

	now := time.Now()

	usr := User{
		ID:          uuid.New(),
		Name:        nu.Name,
		Email:       nu.Email,
		Department:  nu.Department,
		CreditCard:  nu.CreditCard,
		DateCreated: now,
		DateUpdated: now,
	}
	if err := uc.storer.Create(ctx, usr); err != nil {
		uc.log.Error(ctx, "error creating user", "message", err)
		return User{}, fmt.Errorf("create: %w", err)
	}

	return usr, nil
}

func (c *Core) Query(ctx context.Context, filter QueryFilter) ([]User, error) {
	return c.storer.Query(ctx, filter, 1, 10)
}
