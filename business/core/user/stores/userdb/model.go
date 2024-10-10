package userdb

import (
	"context"
	"fmt"
	"net/mail"
	"time"

	"github.com/google/uuid"

	"github.com/Zanda256/commitsmart-task/business/core/user"
	documentStore "github.com/Zanda256/commitsmart-task/business/data/docStore"
	"github.com/Zanda256/commitsmart-task/foundation/keystore"
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
	// byt, err := usr.ID.MarshalBinary()
	// if err != nil {
	// 	return DbUser{}, err
	// }
	//	usrID, err := bson.Marshal(byte)
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

func toCoreUser(ctx context.Context, c *documentStore.DocStorage, dbUsr DbUser) (user.User, error) {
	//em := dbUsr.Email.(primitive.Binary).Data
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
	fmt.Printf("\n\nmodel.go:71 addr %+v\n\n", addr)

	//	var uid []byte
	fmt.Printf("\n\nmodel.go:73 db.UserID %+v\ntype : %T\n", dbUsr.UserID, dbUsr.UserID)
	// err := bson.Unmarshal(dbUsr.UserID.([]byte), &uid)
	// if err != nil {
	// 	return user.User{}, err
	// }
	//k, v := bson.UnmarshalValue(bson, data []byte, val interface{})
	// var (
	// 	idbyts []byte
	// 	ok     bool
	// )
	// if idbyts, ok = dbUsr.UserID.([]byte); ok {
	// 	fmt.Printf("\n\nidbyts %+v\n\n", idbyts)
	// }

	// id, err := uuid.FromBytes(idbyts)
	// if err != nil {
	// 	return user.User{}, err
	// }

	cc := ""

	switch dbUsr.CreditCard.(type) {
	case primitive.Binary:
		c := dbUsr.CreditCard.(primitive.Binary)
		cc = string(c.Data)

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
		CreditCard:  cc,
		DateCreated: dbUsr.DateCreated.In(time.Local),
		DateUpdated: dbUsr.DateUpdated.In(time.Local),
	}
	return usr, nil
}

func toCoreUserSlice(ctx context.Context, c *documentStore.DocStorage, dbUsers []DbUser) ([]user.User, error) {
	usrs := make([]user.User, len(dbUsers))
	for i, dbUsr := range dbUsers {
		var err error
		usrs[i], err = toCoreUser(ctx, c, dbUsr)
		if err != nil {
			return nil, fmt.Errorf("toCoreUser error: %w", err)
		}
	}
	return usrs, nil
}
