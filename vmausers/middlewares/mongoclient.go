package middlewares

import (
	"context"
	"time"
	"vmausers/constants"
	"vmausers/database"
	"vmausers/helper"
	"vmausers/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateUser(user *models.User) error {
	client, err := database.NewConnection(&helper.AppConfig)
	if err != nil {
		return err
	}

	db := client.Database(helper.AppConfig.Mongodb.Database)

	if user.ID.IsZero() {
		user.ID = primitive.NewObjectID()
	}

	err = user.Create(context.Background(), db, constants.User_collection, &user)

	client.Disconnect(context.Background())

	return err
}

func GetUserById(id string) (*models.User, error) {
	client, err := database.NewConnection(&helper.AppConfig)
	if err != nil {
		return nil, err
	}

	db := client.Database(helper.AppConfig.Mongodb.Database)

	objectId, _ := primitive.ObjectIDFromHex(id)
	user := models.User{}
	err = user.ReadOne(context.Background(), db, constants.User_collection, bson.M{"_id": objectId}, &user)

	client.Disconnect(context.Background())
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	client, err := database.NewConnection(&helper.AppConfig)
	if err != nil {
		return nil, err
	}

	db := client.Database(helper.AppConfig.Mongodb.Database)

	user := models.User{}
	err = user.ReadOne(context.Background(), db, constants.User_collection, bson.M{"email": email}, &user)

	client.Disconnect(context.Background())
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetManyUsers() (*[]models.User, error) {
	client, err := database.NewConnection(&helper.AppConfig)

	db := client.Database(helper.AppConfig.Mongodb.Database)

	user := models.User{}
	var listOfUsers []models.User
	err = user.ReadMany(context.Background(), db, constants.User_collection, bson.D{{}}, &listOfUsers)

	client.Disconnect(context.Background())
	if err != nil {
		return nil, err
	}

	return &listOfUsers, nil
}

func DeleteUser(user models.User) error {
	client, err := database.NewConnection(&helper.AppConfig)
	if err != nil {
		return err
	}

	db := client.Database(helper.AppConfig.Mongodb.Database)
	err = user.DeleteOne(context.Background(), db, constants.User_collection, bson.M{"_id": user.ID})

	client.Disconnect(context.Background())
	return err
}

func DeleteUserById(id string) error {
	user, err := GetUserById(id)
	if err != nil {
		err = DeleteUser(*user)
	}

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
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
			"age":        user.Age,
			"address": bson.M{
				"street":  user.Address.Street,
				"city":    user.Address.City,
				"state":   user.Address.State,
				"country": user.Address.Country,
			},
			"basemodel": bson.M{
				"updated_at": primitive.NewDateTimeFromTime(time.Now()),
			},
		},
	}
	err = user.UpdateOne(context.Background(), db, constants.User_collection, bson.M{"_id": user.ID}, update)

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
	err = user.UpdateOne(context.Background(), db, constants.User_collection, bson.M{"_id": user.ID}, update)
	if err != nil {
		user.Password = newPass
	}

	client.Disconnect(context.Background())
	return err
}

func UpdateUserEmail(user *models.User, newEmail string) error {
	client, err := database.NewConnection(&helper.AppConfig)
	if err != nil {
		return err
	}

	db := client.Database(helper.AppConfig.Mongodb.Database)
	update := bson.M{
		"$set": bson.M{
			"email":                newEmail,
			"basemodel.updated_at": primitive.NewDateTimeFromTime(time.Now()),
		},
	}
	err = user.UpdateOne(context.Background(), db, constants.User_collection, bson.M{"_id": user.ID}, update)
	if err != nil {
		user.Email = newEmail
	}

	client.Disconnect(context.Background())
	return err
}
