package userdb

import (
	"context"
	"net/mail"
	"time"

	"github.com/Zanda256/commitsmart-task/business/core/user"
	documentStore "github.com/Zanda256/commitsmart-task/business/data/docStore"
	"github.com/Zanda256/commitsmart-task/foundation/keystore"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// between the app and the database.
type dbUser struct {
	ID          uuid.UUID        `json:"user_id"  bson:"user_id"`
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
		ID:          usr.ID,
		Name:        usr.Name,
		Email:       em,
		Department:  usr.Department,
		CreditCard:  cc,
		DateCreated: usr.DateCreated.UTC(),
		DateUpdated: usr.DateCreated.UTC(),
	}
	return dbUsr, nil
}

func toCoreUser(dbUsr dbUser) (user.User, error) {
	addr := mail.Address{
		Address: string(dbUsr.Email.Data),
	}
	usr := user.User{
		ID:          dbUsr.ID,
		Name:        dbUsr.Name,
		Email:       addr,
		Department:  dbUsr.Department,
		CreditCard:  string(dbUsr.CreditCard.Data),
		DateCreated: dbUsr.DateCreated.In(time.Local),
		DateUpdated: dbUsr.DateUpdated.In(time.Local),
	}
	return usr, nil
}
