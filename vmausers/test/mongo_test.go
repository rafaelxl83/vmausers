package test

import (
	"context"
	_ "fmt" // no more error
	"testing"
	"vmausers/helper"
	"vmausers/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Config helper.Config = *helper.NewConfig(
	"mongodb+srv://cluster0.t4xtjka.mongodb.net/",
	"atlas-xk9ebu-shard-0",
	"../../certificate/X509-cert-5347953578960200531.crt",
	"../../certificate/X509-key-5347953578960200531.pem",
)

var TestUser models.User = models.User{
	ID:       SetObjectID(),
	Name:     "John",
	LastName: "Doe",
	Email:    "john.doe@example.com",
	Address: models.Address{
		Street:  "Nowhere",
		City:    "Caddo",
		State:   "Oklahoma",
		Country: "United States",
	},
	Password: *models.NewPassword(""),
}

func SetObjectID() primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex("64cd5d112a05e4baff910531")
	return id
}

func TestConnected(t *testing.T) {
	client, err := helper.Connect(
		Config.Mongodb.Serveruri,
		Config.Mongodb.CaFilePath,
		Config.Mongodb.CaKeyFilePath,
		Config.Mongodb.ReplicaSet)

	client.Disconnect(context.Background())
	_ = err
	AssertEqual(t, err, nil)
}

func TestCreateUser(t *testing.T) {
	client, _ := helper.Connect(
		Config.Mongodb.Serveruri,
		Config.Mongodb.CaFilePath,
		Config.Mongodb.CaKeyFilePath,
		Config.Mongodb.ReplicaSet)

	db := client.Database("test_db")

	// Create a new user
	err := TestUser.Create(context.Background(), db, "users", &TestUser)

	client.Disconnect(context.Background())
	_ = err
	AssertEqual(t, err, nil)
}

func TestReadUser(t *testing.T) {
	client, _ := helper.Connect(
		Config.Mongodb.Serveruri,
		Config.Mongodb.CaFilePath,
		Config.Mongodb.CaKeyFilePath,
		Config.Mongodb.ReplicaSet)

	db := client.Database("test_db")

	var readUser models.User
	err := readUser.Read(context.Background(), db, "users", bson.M{"_id": TestUser.ID}, &readUser)
	if err != nil {
		_ = TestUser.Create(context.Background(), db, "users", &TestUser)
		_ = readUser.Read(context.Background(), db, "users", bson.M{"_id": TestUser.ID}, &readUser)
	}

	client.Disconnect(context.Background())
	AssertEqual(t, TestUser.ID, readUser.ID)
}

func TestUpdateUser(t *testing.T) {
	client, _ := helper.Connect(
		Config.Mongodb.Serveruri,
		Config.Mongodb.CaFilePath,
		Config.Mongodb.CaKeyFilePath,
		Config.Mongodb.ReplicaSet)

	db := client.Database("test_db")

	// Update a user's email
	email := "john.doe_updated@example.com"
	update := bson.M{"$set": bson.M{"email": email, "updated_at": primitive.NewDateTimeFromTime(TestUser.UpdatedAt)}}
	err := TestUser.Update(context.Background(), db, "users", bson.M{"_id": TestUser.ID}, update)

	if err != nil {
		_ = TestUser.Create(context.Background(), db, "users", &TestUser)
		err = TestUser.Update(context.Background(), db, "users", bson.M{"_id": TestUser.ID}, update)
	}

	client.Disconnect(context.Background())
	AssertEqual(t, err, nil)
}

func TestDeleteUser(t *testing.T) {
	client, _ := helper.Connect(
		Config.Mongodb.Serveruri,
		Config.Mongodb.CaFilePath,
		Config.Mongodb.CaKeyFilePath,
		Config.Mongodb.ReplicaSet)

	db := client.Database("test_db")

	// Delete a user by ID
	err := TestUser.Delete(context.Background(), db, "users", bson.M{"_id": TestUser.ID})

	if err != nil {
		_ = TestUser.Create(context.Background(), db, "users", &TestUser)
		err = TestUser.Delete(context.Background(), db, "users", bson.M{"_id": TestUser.ID})
	}

	client.Disconnect(context.Background())
	AssertEqual(t, err, nil)
}
