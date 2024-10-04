package documentStore

import (
	"context"
	"errors"
	"fmt"

	"github.com/Zanda256/commitsmart-task/foundation/keystore"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	ENCRYPTION_ALGORITHM_DETERMINISTIC = "AEAD_AES_256_CBC_HMAC_SHA_512-Deterministic"

	ENCRYPTION_ALGORITHM_RANDOM = "AEAD_AES_256_CBC_HMAC_SHA_512-Random"

	keyAltNames = "keyAltNames"
)

// The MongoDB namespace (db.collection) used to store the encryption data
// keys.
var (
	keyVaultDBName, keyVaultCollName = "encryption", "__keyVault"
	keyVaultNamespace                = keyVaultDBName + "." + keyVaultCollName
)

func NewCSFLEObj(client *mongo.Client, kmsProviders map[string]map[string]interface{}) (*mongo.ClientEncryption, error) {

	// Create a ClientEncryption Instance
	clientEncryptionOpts := options.ClientEncryption().SetKeyVaultNamespace(keyVaultNamespace).SetKmsProviders(kmsProviders)
	clientEnc, err := mongo.NewClientEncryption(client, clientEncryptionOpts)
	if err != nil {
		return nil, fmt.Errorf("NewClientEncryption error %v", err)
	}
	return clientEnc, nil
}

func CreateDEK(ctx context.Context, clientEncryption *mongo.ClientEncryption) (primitive.Binary, error) {
	var (
		dataKeyData = bson.M{}
		err         error
	)
	keyResult := clientEncryption.GetKeyByAltName(ctx, keystore.UserDEKeyAlias)

	if err = keyResult.Decode(dataKeyData); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// Create a new data key for the encrypted field.
			dataKeyOpts := options.DataKey().
				SetKeyAltNames([]string{keystore.UserDEKeyAlias})

			NewKeyID, err := clientEncryption.CreateDataKey(
				ctx,
				"local",
				dataKeyOpts)
			if err != nil {
				return primitive.Binary{}, err
			}
			return NewKeyID, nil
		} else {
			return primitive.Binary{}, fmt.Errorf("error creating Data Encryption Key; %w", err)
		}
	}
	key := dataKeyData["_id"].(primitive.Binary)
	return key, nil
}

func (store *DocStorage) GetDataKeyID(ctx context.Context, alias string) (primitive.Binary, error) {
	var keyID primitive.Binary
	if result := store.EncryptMgr.GetKeyByAltName(ctx, alias); result.Err() == nil {
		err := result.Decode(keyID)
		if err != nil {
			return primitive.Binary{}, err
		}
	}
	return keyID, nil
}

func (store *DocStorage) EncryptStrVal(ctx context.Context, v string, dataKeyAltName string) (primitive.Binary, error) {
	// Create a bson.RawValue to encrypt and encrypt it using the key that was
	// just created.
	rawValueType, rawValueData, err := bson.MarshalValue(v)
	if err != nil {
		return primitive.Binary{}, err
	}

	rawValue := bson.RawValue{Type: rawValueType, Value: rawValueData}

	encryptionOpts := options.Encrypt().
		SetAlgorithm(ENCRYPTION_ALGORITHM_DETERMINISTIC).
		SetKeyAltName(dataKeyAltName)

	encryptedField, err := store.EncryptMgr.Encrypt(
		ctx,
		rawValue,
		encryptionOpts)

	if err != nil {
		return primitive.Binary{}, err
	}

	return encryptedField, nil
}

func (store *DocStorage) DencryptStrVal(ctx context.Context, encrypted primitive.Binary) (string, error) {
	decryptedRaw, err := store.EncryptMgr.Decrypt(
		ctx,
		encrypted)
	if err != nil {
		return "", nil
	}
	return decryptedRaw.String(), nil
}
