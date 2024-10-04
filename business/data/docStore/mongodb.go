package documentStore

import (
	"context"
	"fmt"
	"time"

	"github.com/Zanda256/commitsmart-task/foundation/keystore"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DocStorage struct {
	Client     *mongo.Client
	EncryptMgr *mongo.ClientEncryption
}

func StartDB(opts *options.ClientOptions) (*mongo.Client, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to MongoDB Successfully!")

	return client, nil
}

func StatusCheck(ctx context.Context, db *mongo.Database) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	var pingError error
	for attempts := 1; ; attempts++ {
		var result bson.M
		pingError = db.RunCommand(ctx, bson.D{{"ping", 1}}).Decode(&result)
		if pingError == nil {
			break
		}
		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	fmt.Println("MongoDB connection status healthy")
	return nil
}

func OpenCollection(db *mongo.Database, collectionName string) *mongo.Collection {
	var collection = db.Collection(collectionName)
	return collection
}

func StartEncryptedDB(opts *options.ClientOptions) (*DocStorage, error) {

	// Create Customer Master Key
	localMasterKey, err := keystore.MustGetMKey()
	if err != nil {
		return nil, err
	}
	kmsProviders := map[string]map[string]interface{}{
		"local": {
			"key": localMasterKey,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	autoEncryptionOpts := options.AutoEncryption().
		SetKmsProviders(kmsProviders).
		SetKeyVaultNamespace(keyVaultNamespace).
		SetBypassAutoEncryption(true)
	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts.SetAutoEncryptionOptions(autoEncryptionOpts))

	if err != nil {
		return nil, err
	}

	csfleImpl, err := NewCSFLEObj(client, kmsProviders)
	if err != nil {
		return nil, err
	}

	_, err = CreateDEK(ctx, csfleImpl)
	if err != nil {
		return nil, err
	}

	// -------------------------------------------------------------------------------------------------
	// Define index on "keyAltNames" field of the "__keyVault" collection
	// Define the index model
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "keyAltNames", Value: 1},
		},
		Options: options.Index().SetName(fmt.Sprintf("%s.%s_index", keyVaultNamespace, keyAltNames)).
			SetUnique(true).
			SetPartialFilterExpression(bson.D{
				{Key: keyAltNames, Value: bson.D{
					{Key: "$exists", Value: true},
				}},
			}),
	}

	// Create the index
	keyVaultColl := client.Database("encryption").Collection("__keyVault")
	_, err = keyVaultColl.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		fmt.Println("Error creating index:", err)
		return nil, err
	}

	fmt.Println("Connected to MongoDB Successfully!")
	strg := &DocStorage{
		Client:     client,
		EncryptMgr: csfleImpl,
	}

	return strg, nil
}

// func localMasterKey() []byte {
// 	key := make([]byte, 96)
// 	if _, err := rand.Read(key); err != nil {
// 		logger.Fatalf("Unable to create a random 96 byte data key: %v", err)
// 	}
// 	if err := ioutil.WriteFile("master-key.txt", key, 0644); err != nil {
// 		log.Fatalf("Unable to write key to file: %v", err)
// 	}
// 	return key
// }
