package middlewares

import (
	"context"
	"time"
	"vmausers/database"
	"vmausers/helper"
	"vmausers/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var user_collection = "users"

func fillMissingBase(user *models.User) {
	if user.ID.IsZero() {
		user.ID = primitive.NewObjectID()
	}
}

func CreateUser(user *models.User) error {
	client, err := database.NewConnection(&helper.AppConfig)
	if err != nil {
		return err
	}

	db := client.Database(helper.AppConfig.Mongodb.Database)

	user.ID = primitive.NewObjectID()
	err = user.Create(context.Background(), db, user_collection, &user)

	client.Disconnect(context.Background())

	return err
}

func GetUserById(id primitive.ObjectID) (*models.User, error) {
	client, err := database.NewConnection(&helper.AppConfig)
	if err != nil {
		return nil, err
	}

	db := client.Database(helper.AppConfig.Mongodb.Database)

	user := models.User{}
	err = user.ReadOne(context.Background(), db, user_collection, bson.M{"_id": id}, &user)

	client.Disconnect(context.Background())
	return &user, err
}

func GetUserByEmail(email string) (*models.User, error) {
	client, err := database.NewConnection(&helper.AppConfig)
	if err != nil {
		return nil, err
	}

	db := client.Database(helper.AppConfig.Mongodb.Database)

	user := models.User{}
	err = user.ReadOne(context.Background(), db, user_collection, bson.M{"email": email}, &user)

	client.Disconnect(context.Background())
	return &user, err
}

func DeleteUser(user models.User) error {
	client, err := database.NewConnection(&helper.AppConfig)
	if err != nil {
		return err
	}

	db := client.Database(helper.AppConfig.Mongodb.Database)
	err = user.DeleteOne(context.Background(), db, user_collection, bson.M{"_id": user.ID})

	client.Disconnect(context.Background())
	return err
}

func DeleteUserById(id primitive.ObjectID) error {
	user := models.User{
		ID: id,
	}
	err := DeleteUser(user)
	return err
}

func DeleteUserByEmail(email string) error {
	user, err := GetUserByEmail(email)
	if err != nil {
		err = DeleteUser(*user)
	}

	return err
}

func UpdateUser(user *models.User) error {
	client, err := database.NewConnection(&helper.AppConfig)
	if err != nil {
		return err
	}

	db := client.Database(helper.AppConfig.Mongodb.Database)
	update := bson.M{
		"$set": bson.M{
			"first_name":           user.FirstName,
			"last_name":            user.LastName,
			"email":                user.Email,
			"address.street":       user.Address.Street,
			"address.city":         user.Address.City,
			"address.state":        user.Address.State,
			"address.country":      user.Address.Country,
			"basemodel.updated_at": primitive.NewDateTimeFromTime(time.Now()),
		},
	}
	err = user.UpdateOne(context.Background(), db, user_collection, bson.M{"_id": user.ID}, update)

	client.Disconnect(context.Background())
	return err
}

func UpdateUserPassword(user *models.User, newPass models.Password) error {
	client, err := database.NewConnection(&helper.AppConfig)
	if err != nil {
		return err
	}

	db := client.Database(helper.AppConfig.Mongodb.Database)
	update := bson.M{
		"$set": bson.M{
			"password.encrypted":   newPass.EncryptedPass,
			"password.created_at":  primitive.NewDateTimeFromTime(newPass.CreatedAt),
			"password.expire":      primitive.NewDateTimeFromTime(newPass.Expire),
			"basemodel.updated_at": primitive.NewDateTimeFromTime(time.Now()),
		},
	}
	err = user.UpdateOne(context.Background(), db, user_collection, bson.M{"_id": user.ID}, update)
	if err != nil {
		user.Password = newPass
	}

	client.Disconnect(context.Background())
	return err
}
