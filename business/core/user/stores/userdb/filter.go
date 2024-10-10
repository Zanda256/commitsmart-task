package userdb

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Zanda256/commitsmart-task/business/core/user"
	"github.com/Zanda256/commitsmart-task/foundation/keystore"
)

func (s *Store) ApplyFilter(filter user.QueryFilter) bson.D {

	fmt.Printf("\ndb.ApplyFilter : filter : %#v\n", filter)
	var (
		userIDQ     uuid.UUID
		nameQ       string
		emailQ      primitive.Binary
		departmentQ string
	)

	if filter.UserID != nil {
		userIDQ = *filter.UserID
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

	//return bson.D{
	//	{"user_id", bson.D{{"$eq", userIDQ}}},
	//	{"email", bson.D{{"$eq", emailQ}}},
	//	{"name", bson.D{{"$eq", nameQ}}},
	//	{"department", bson.D{{"$eq", departmentQ}}},
	//}
	t := 1
	if t > 0 {
		return bson.D{}
	} else {
		return bson.D{
			{"user_id", userIDQ},
			{"email", emailQ},
			{"name", nameQ},
			{"department", departmentQ},
		}
	}
}
