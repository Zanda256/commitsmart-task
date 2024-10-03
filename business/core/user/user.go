package user

import (
	"context"

	"github.com/Zanda256/commitsmart-task/foundation/logger"
)

// Storer interface declares the behavior this package needs to perists and
// retrieve data.
type Storer interface {
	Create(ctx context.Context, usr User) error
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
	return User{}, nil
}
