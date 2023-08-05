package models

import (
	"time"
	"vmausers/database"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	database.BaseModel
	Address
	Password
	Name     string             `bson:"name"`
	LastName string             `bson:"last_name"`
	Email    string             `bson:"email"`
	ID       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
}

func NewUser(
	Name string, LastName string, Email string,
	Street string, City string, State string, Country string) *User {
	user := User{
		ID:       primitive.NewObjectID(),
		Name:     Name,
		LastName: LastName,
		Email:    Email,
		Address: Address{
			Street:  Street,
			City:    City,
			State:   State,
			Country: Country,
		},
		Password: *NewPassword(""),
		BaseModel: database.BaseModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	return &user
}
