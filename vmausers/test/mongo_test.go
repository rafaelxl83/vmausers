package test

import (
	"context"
	_ "fmt" // no more error
	"testing"
	"time"
	"vmausers/database"
	"vmausers/helper"
	"vmausers/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Config helper.Config = *helper.NewConfig(
	"mongodb+srv://cluster0.t4xtjka.mongodb.net/",
	"atlas-xk9ebu-shard-0",
	"test_db",
	"../../certificate/X509-cert-5347953578960200531.crt",
	"../../certificate/X509-key-5347953578960200531.pem",
)

var TestUserId = "64cd5d112a05e4baff910531"

var TestUser models.User = models.User{
	ID:        SetObjectID(),
	FirstName: "John",
	LastName:  "Doe",
	Age:       40,
	Email:     "john.doe@example.com",
	Address: models.Address{
		Street:  "Nowhere",
		City:    "Caddo",
		State:   "Oklahoma",
		Country: "United States",
	},
	Password: *models.NewPassword(""),
}

func SetObjectID() primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(TestUserId)
	return id
}

func TestConnected(t *testing.T) {
	client, err := database.Connect(
		Config.Mongodb.Serveruri,
		Config.Mongodb.CaFilePath,
		Config.Mongodb.CaKeyFilePath,
		Config.Mongodb.ReplicaSet)

	client.Disconnect(context.Background())
	_ = err
	AssertEqual(t, err, nil)
}

func TestCreateUser(t *testing.T) {
	client, _ := database.NewConnection(&Config)

	db := client.Database(Config.Mongodb.Database)

	// Create a new user
	err := TestUser.Create(context.Background(), db, "users", &TestUser)

	client.Disconnect(context.Background())
	_ = err
	AssertEqual(t, err, nil)
}

func TestReadUser(t *testing.T) {
	client, _ := database.NewConnection(&Config)

	db := client.Database(Config.Mongodb.Database)

	var readUser models.User
	err := readUser.ReadOne(context.Background(), db, "users", bson.M{"_id": TestUser.ID}, &readUser)
	if err != nil {
		_ = TestUser.Create(context.Background(), db, "users", &TestUser)
		_ = readUser.ReadOne(context.Background(), db, "users", bson.M{"_id": TestUser.ID}, &readUser)
	}

	client.Disconnect(context.Background())
	AssertEqual(t, TestUser.ID, readUser.ID)
}

func TestReadManyUsers(t *testing.T) {
	client, err := database.NewConnection(&Config)

	db := client.Database(Config.Mongodb.Database)

	var listOfUsers []models.User
	err = TestUser.ReadMany(context.Background(), db, "users", bson.D{{}}, &listOfUsers)
	if err != nil {
		_ = TestUser.Create(context.Background(), db, "users", &TestUser)
		_ = TestUser.ReadMany(context.Background(), db, "users", bson.D{{}}, &listOfUsers)
	}

	client.Disconnect(context.Background())
	AssertNotEqual(t, len(listOfUsers), 0)
}

func TestUpdateUser(t *testing.T) {
	client, _ := database.NewConnection(&Config)

	db := client.Database(Config.Mongodb.Database)

	// Update a user's email
	email := "john.doe_updated@example.com"
	update := bson.M{"$set": bson.M{"email": email, "basemodel.updated_at": primitive.NewDateTimeFromTime(time.Now())}}
	err := TestUser.UpdateOne(context.Background(), db, "users", bson.M{"_id": TestUser.ID}, update)

	client.Disconnect(context.Background())
	AssertEqual(t, err, nil)
}

func TestDeleteUser(t *testing.T) {
	client, _ := database.NewConnection(&Config)

	db := client.Database(Config.Mongodb.Database)

	// Delete a user by ID
	err := TestUser.DeleteOne(context.Background(), db, "users", bson.M{"_id": TestUser.ID})

	client.Disconnect(context.Background())
	AssertEqual(t, err, nil)
}
