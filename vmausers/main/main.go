package main

import (
	"context"
	"flag"
	"fmt"
	"vmausers/helper"
	"vmausers/models"

	"go.mongodb.org/mongo-driver/bson"
)

func ConnectAndTestMongo() {
	configFile := flag.String("config", "config.json", "a config file")
	flag.Parse()
	config, err := helper.LoadConfig(*configFile)
	if err != nil || len(config.Mongodb.Serveruri) == 0 {
		fmt.Printf("Error opening the config file: %v", err)
		return
	}

	client, err := helper.Connect(
		config.Mongodb.Serveruri,
		config.Mongodb.CaFilePath,
		config.Mongodb.CaKeyFilePath,
		config.Mongodb.ReplicaSet)
	if err != nil {
		panic(err)
	}

	db := client.Database("test_db")

	// Create a new user
	user := models.NewUser(
		"John", "Doe", "john.doe@example.com",
		"Nowhere", "Caddo", "Oklahoma", "United States",
	)

	err = user.Create(context.Background(), db, "users", &user)
	if err != nil {
		panic(err)
	}
	fmt.Printf("User created: %v\n", user)

	// Read a user by ID
	var readUser models.User
	err = readUser.Read(context.Background(), db, "users", bson.M{"_id": user.ID}, &readUser)
	if err != nil {
		panic(err)
	}
	fmt.Printf("User read: %v\n", readUser)

	// Update a user's email
	//update := bson.M{"$set": bson.M{"email": "john.doe_updated@example.com", "updated_at": primitive.NewDateTimeFromTime(user.UpdatedAt)}}
	//err = user.Update(context.Background(), db, "users", bson.M{"_id": user.ID}, update)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("User updated: %v\n", user)

	// Delete a user by ID
	//err = user.Delete(context.Background(), db, "users", bson.M{"_id": user.ID})
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("User deleted")
}

func main() {
	ConnectAndTestMongo()
}
