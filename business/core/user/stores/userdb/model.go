package userdb

import (
	"context"
	"github.com/google/uuid"
	"net/mail"
	"time"

	"github.com/Zanda256/commitsmart-task/business/core/user"
	documentStore "github.com/Zanda256/commitsmart-task/business/data/docStore"
	"github.com/Zanda256/commitsmart-task/foundation/keystore"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// dbUser represent the structure we need for moving data
// between the app and the database.
type dbUser struct {
	UserID      string           `json:"user_id"  bson:"user_id"`
	Name        string           `json:"name"  bson:"name"`
	Email       primitive.Binary `json:"email"  bson:"email"`
	Department  string           `json:"department" bson:"department"`
	CreditCard  primitive.Binary `json:"credit_card" bson:"credit_card"`
	DateCreated time.Time        `json:"date_created" bson:"date_created"`
	DateUpdated time.Time        `json:"date_updated" bson:"date_updated"`
}

func ToDbUser(ctx context.Context, c *documentStore.DocStorage, usr user.User) (dbUser, error) {
	cc, err := c.EncryptStrVal(ctx, usr.CreditCard, keystore.UserDEKeyAlias)
	if err != nil {
		return dbUser{}, err
	}
	em, err := c.EncryptStrVal(ctx, usr.Email.String(), keystore.UserDEKeyAlias)
	if err != nil {
		return dbUser{}, err
	}
	dbUsr := dbUser{
		UserID:      usr.ID.String(),
		Name:        usr.Name,
		Email:       em,
		Department:  usr.Department,
		CreditCard:  cc,
		DateCreated: usr.DateCreated.UTC(),
		DateUpdated: usr.DateCreated.UTC(),
	}
	return dbUsr, nil
}

func toCoreUser(ctx context.Context, c *documentStore.DocStorage, dbUsr dbUser) (user.User, error) {
	addr := mail.Address{
		Address: string(dbUsr.Email.Data),
	}
	id, err := uuid.Parse(dbUsr.UserID)
	if err != nil {
		return user.User{}, err
	}

	usr := user.User{
		ID:          id,
		Name:        dbUsr.Name,
		Email:       addr,
		Department:  dbUsr.Department,
		CreditCard:  string(dbUsr.CreditCard.Data),
		DateCreated: dbUsr.DateCreated.In(time.Local),
		DateUpdated: dbUsr.DateUpdated.In(time.Local),
	}
	return usr, nil
}

func toCoreUserSlice(ctx context.Context, c *documentStore.DocStorage, dbUsers []dbUser) ([]user.User, error) {
	usrs := make([]user.User, len(dbUsers))
	for i, dbUsr := range dbUsers {
		var err error
		usrs[i], err = toCoreUser(ctx, c, dbUsr)
		if err != nil {
			return nil, err
		}
	}
	return usrs, nil
}
