package userdb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Zanda256/commitsmart-task/business/core/user"
	"github.com/Zanda256/commitsmart-task/foundation/keystore"
)

func (s *Store) ApplyFilter(filter user.QueryFilter) bson.D {
	var (
		userIDQ     string
		nameQ       string
		emailQ      primitive.Binary
		departmentQ string
	)

	if filter.UserID != nil {
		userIDQ = filter.UserID.String()
	}

	if filter.Email != nil {
		em, err := s.dbClient.EncryptVal(context.Background(), filter.Email.String(), keystore.UserDEKeyAlias)
		if err != nil {
			return bson.D{}
		}
		emailQ = em
	}

	if filter.Name != nil {
		nameQ = *filter.Name
	}

	if filter.Department != nil {
		departmentQ = *filter.Department
	}

	return bson.D{
		{"user_id", userIDQ},
		{"email", emailQ},
		{"name", nameQ},
		{"department", departmentQ},
	}
}
