package userdb

import (
	"context"
	"fmt"
	"time"

	"github.com/Zanda256/commitsmart-task/foundation/keystore"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Zanda256/commitsmart-task/business/core/user"
	documentStore "github.com/Zanda256/commitsmart-task/business/data/docStore"
	"github.com/Zanda256/commitsmart-task/foundation/logger"
)

// Store manages the set of APIs for rate search-db access.
type Store struct {
	log      *logger.Logger
	dbClient *documentStore.DocStorage
	dbName   string
	userColl string
}

// NewStore constructs the api for data access.
func NewStore(log *logger.Logger, dbName string, dbClient *documentStore.DocStorage, userCollectionName string) *Store {
	return &Store{
		log:      log,
		dbClient: dbClient,
		userColl: userCollectionName,
		dbName:   dbName,
	}
}

// func (s *Store) encryptStrValue(val, dataKeyId string) error {
// 	nameRawValueType, nameRawValueData, err := bson.MarshalValue(val)
// 	if err != nil {
// 		panic(err)
// 	}
// 	nameRawValue := bson.RawValue{Type: nameRawValueType, Value: nameRawValueData}
// 	nameEncryptionOpts := options.Encrypt().
// 		SetAlgorithm("AEAD_AES_256_CBC_HMAC_SHA_512-Deterministic").
// 		SetKeyID(dataKeyId)
// 	nameEncryptedField, err := s.encryptMgr.Encrypt(
// 		context.TODO(),
// 		nameRawValue,
// 		nameEncryptionOpts)
// 	if err != nil {
// 		panic(err)
// 	}
// }

// Create inserts a new user into the database.
func (s *Store) Create(ctx context.Context, usr user.User) error {
	_, err := s.saveUser(ctx, usr)
	if err != nil {
		return fmt.Errorf("s.saveUser: %v", err)
	}
	return nil
}

// Create inserts a new user into the database.
func (s *Store) saveUser(ctx context.Context, usr user.User) (user.User, error) {
	email, err := s.dbClient.EncryptStrVal(ctx, usr.Email.String(), keystore.UserDEKeyAlias)
	if err != nil {
		return user.User{}, fmt.Errorf("EncryptStrVal error: %v", err)
	}

	cc, err := s.dbClient.EncryptStrVal(ctx, usr.CreditCard, keystore.UserDEKeyAlias)
	if err != nil {
		return user.User{}, fmt.Errorf("EncryptStrVal error: %v", err)
	}
	now := time.Now()
	var doc bson.D = bson.D{
		{"user_id", usr.ID},
		{"name", usr.Name},
		{"email", email},
		{"department", usr.Department},
		{"credit_card", cc},
		{"date_created", now},
		{"date_updated", now},
	}
	// Insert the document
	coll := documentStore.OpenCollection(s.dbClient.Client.Database(s.dbName), s.userColl)
	_, err = coll.InsertOne(ctx, doc)
	if err != nil {
		return user.User{}, fmt.Errorf("coll.InsertOne error %v", err)
	}
	return usr, nil
}

// Query retrieves a list of rates from the database.
func (s *Store) Query(ctx context.Context, filter user.QueryFilter, pageNumber int, rowsPerPage int) ([]user.User, error) {
	// filtering
	dbFilter := s.ApplyFilter(filter)

	s.log.Info(ctx, "Check filter", "dbFilter", dbFilter)

	// pagination
	skip := (pageNumber - 1) * rowsPerPage
	s.log.Info(ctx, "skip", "val", skip)
	s.log.Info(ctx, "limit", "val", rowsPerPage)

	var collection *mongo.Collection = documentStore.OpenCollection(s.dbClient.Client.Database(s.dbName), s.userColl)

	//sorting
	sortknob := bson.D{
		{
			"department", 1,
		},
		{
			"email", 1,
		},
		//{
		//	"user_id", 1,
		//},
		//{
		//	"start_date", 1,
		//},
	}

	// Set query opts
	opts := options.Find().SetSort(sortknob).SetSkip(int64(skip)).SetLimit(int64(rowsPerPage)) //pageNumber *

	// execute query
	cursor, err := collection.Find(ctx, dbFilter, opts)
	if err != nil {
		s.log.Error(ctx, "error in collection.Find", "msg", err.Error())
		return []user.User{}, err
	}

	// deserialize from binary json
	var results []dbUser
	if err = cursor.All(context.Background(), &results); err != nil {
		s.log.Error(ctx, "error in cursor.All", "msg", err.Error())
		return []user.User{}, err
	}
	s.log.Info(ctx, "Check dbResults", "len", len(results))

	return toCoreUserSlice(ctx, s.dbClient, results)
}
