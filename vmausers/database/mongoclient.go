package database

import (
	"context"
	"vmausers/helper"
)

var user_collection = "users"

func (b BaseModel) CreateUser(user interface{}) error {
	client, err := NewConnection(&helper.AppConfig)
	if err != nil {
		return err
	}

	db := client.Database(helper.AppConfig.Mongodb.Database)
	err = b.Create(context.Background(), db, user_collection, &user)
	client.Disconnect(context.Background())

	return err
}
