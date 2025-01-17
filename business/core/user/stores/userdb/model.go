package userdb

import (
	"context"
	"fmt"
	"net/mail"
	"time"

	"github.com/Zanda256/commitsmart-task/business/core/user"
	documentStore "github.com/Zanda256/commitsmart-task/business/data/docStore"
	"github.com/Zanda256/commitsmart-task/foundation/keystore"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DbUser represent the structure we need for moving data
// between the app and the database.
type DbUser struct {
	UserID      uuid.UUID   `json:"user_id"  bson:"user_id"`
	Name        string      `json:"name"  bson:"name"`
	Email       interface{} `json:"email"  bson:"email"`
	Department  string      `json:"department" bson:"department"`
	CreditCard  interface{} `json:"credit_card" bson:"credit_card"`
	DateCreated time.Time   `json:"date_created" bson:"date_created"`
	DateUpdated time.Time   `json:"date_updated" bson:"date_updated"`
}

func ToDbUser(ctx context.Context, c *documentStore.DocStorage, usr user.User) (DbUser, error) {
	cc, err := c.EncryptVal(ctx, usr.CreditCard, keystore.UserDEKeyAlias)
	if err != nil {
		return DbUser{}, err
	}
	em, err := c.EncryptVal(ctx, usr.Email.String(), keystore.UserDEKeyAlias)
	if err != nil {
		return DbUser{}, err
	}

	dbUsr := DbUser{
		UserID:      usr.ID,
		Name:        usr.Name,
		Email:       em,
		Department:  usr.Department,
		CreditCard:  cc,
		DateCreated: usr.DateCreated.UTC(),
		DateUpdated: usr.DateCreated.UTC(),
	}
	return dbUsr, nil
}

func toCoreUser(ctx context.Context, dbUsr DbUser) (user.User, error) {
	emailStr := ""

	switch dbUsr.Email.(type) {
	case primitive.Binary:
		em := dbUsr.Email.(primitive.Binary)
		emailStr = string(em.Data)

	case string:
		emailStr = dbUsr.Email.(string)

	case []byte:
		emailStr = string(dbUsr.Email.([]byte))
	}

	addr := mail.Address{
		Address: emailStr,
	}

	switch dbUsr.CreditCard.(type) {
	case primitive.Binary:
		c := dbUsr.CreditCard.(primitive.Binary)
		emailStr = string(c.Data)

	case string:
		emailStr = dbUsr.CreditCard.(string)

	case []byte:
		emailStr = string(dbUsr.CreditCard.([]byte))
	}

	usr := user.User{
		ID:          dbUsr.UserID,
		Name:        dbUsr.Name,
		Email:       addr,
		Department:  dbUsr.Department,
		CreditCard:  emailStr,
		DateCreated: dbUsr.DateCreated.In(time.Local),
		DateUpdated: dbUsr.DateUpdated.In(time.Local),
	}
	return usr, nil
}

func toCoreUserSlice(ctx context.Context, dbUsers []DbUser) ([]user.User, error) {
	usrs := make([]user.User, len(dbUsers))
	for i, dbUsr := range dbUsers {
		var err error
		usrs[i], err = toCoreUser(ctx, dbUsr)
		if err != nil {
			return nil, fmt.Errorf("toCoreUser error: %w", err)
		}
	}
	return usrs, nil
}
