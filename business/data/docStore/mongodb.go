package documentStore

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
