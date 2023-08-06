package database

import (
	"context"
	"log"
	"time"
	"vmausers/helper"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BaseModel struct {
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func NewConnection(config *helper.Config) (*mongo.Client, error) {
	if config == nil {
		config = &helper.AppConfig
	}

	client, err := Connect(
		config.Mongodb.Serveruri,
		config.Mongodb.CaFilePath,
		config.Mongodb.CaKeyFilePath,
		config.Mongodb.ReplicaSet)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func Connect(uri string, caFilePath string, caKeyFilePath string, replicaSet string) (*mongo.Client, error) {
	uri += "$external"
	uri += "?retryWrites=true&readPreference=primary"
	uri += "&replicaSet=" + replicaSet
	uri += "&srvServiceName=mongodb&connectTimeoutMS=10000"
	uri += "&ssl=true&tlsPrivateKeyFile=" + caKeyFilePath + "&tlsCertificateFile=" + caFilePath

	credential := options.Credential{
		AuthMechanism: "MONGODB-X509",
		AuthSource:    "$external",
	}

	clientOptions := options.Client().ApplyURI(uri).SetAuth(credential)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return client, nil
}

func (m *BaseModel) Create(ctx context.Context, db *mongo.Database, collectionName string, model interface{}) error {
	collection := db.Collection(collectionName)

	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()

	_, err := collection.InsertOne(ctx, model)
	if err != nil {
		return err
	}

	return nil
}

func (m *BaseModel) ReadOne(ctx context.Context, db *mongo.Database, collectionName string, filter interface{}, result interface{}) error {
	collection := db.Collection(collectionName)

	err := collection.FindOne(ctx, filter).Decode(result)
	if err != nil {
		return err
	}

	return nil
}

func (m *BaseModel) UpdateOne(ctx context.Context, db *mongo.Database, collectionName string, filter interface{}, update interface{}) error {
	collection := db.Collection(collectionName)

	m.UpdatedAt = time.Now()
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (m *BaseModel) DeleteOne(ctx context.Context, db *mongo.Database, collectionName string, filter interface{}) error {
	collection := db.Collection(collectionName)

	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
