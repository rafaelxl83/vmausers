package models

import (
	"vmausers/database"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	database.BaseModel
	Address   Address            `bson:"address" binding:"required"`
	Password  Password           `bson:"password" binding:"required"`
	FirstName string             `bson:"first_name" binding:"required"`
	LastName  string             `bson:"last_name"`
	Email     string             `bson:"email" binding:"required,email"`
	ID        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
}

func NewUser(
	FirstName string, LastName string, Email string,
	Street string, City string, State string, Country string) *User {
	user := User{
		ID:        primitive.NewObjectID(),
		FirstName: FirstName,
		LastName:  LastName,
		Email:     Email,
		Address: Address{
			Street:  Street,
			City:    City,
			State:   State,
			Country: Country,
		},
		Password: *NewPassword(""),
	}
	return &user
}
