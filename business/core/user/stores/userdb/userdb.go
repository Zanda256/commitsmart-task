package userdb

import (
	"context"
	"github.com/Zanda256/commitsmart-task/business/core/user"
	documentStore "github.com/Zanda256/commitsmart-task/business/data/docStore"
	"github.com/Zanda256/commitsmart-task/foundation/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

// Store manages the set of APIs for rate search-db access.
type Store struct {
	log      *logger.Logger
	db       *mongo.Database
	userColl *mongo.Collection
}

// NewStore constructs the api for data access.
func NewStore(log *logger.Logger, db *mongo.Database, userCollectionName string) *Store {
	return &Store{
		log:      log,
		db:       db,
		userColl: documentStore.OpenCollection(db, userCollectionName),
	}
}

// Create inserts a new user into the database.
func (s *Store) Create(ctx context.Context, usr user.User) error {

	return nil
}
