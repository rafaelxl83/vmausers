package models

import (
	"net/mail"
	"reflect"
	"vmausers/database"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	database.BaseModel
	Address   Address            `bson:"address" binding:"required"`
	Password  Password           `bson:"password" binding:"required"`
	FirstName string             `bson:"first_name" binding:"required"`
	LastName  string             `bson:"last_name"`
	Age       int                `bson:"age" binding:"required"`
	Email     string             `bson:"email" binding:"required,email"`
	ID        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
}

func NewUser(
	FirstName string, LastName string, Age int, Email string,
	Street string, City string, State string, Country string) *User {
	user := User{
		ID:        primitive.NewObjectID(),
		FirstName: FirstName,
		LastName:  LastName,
		Age:       Age,
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

func IsEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (currUser *User) UpdateValues(newUser User) {
	updateProperty(&currUser.FirstName, newUser.FirstName)
	updateProperty(&currUser.LastName, newUser.LastName)
	updateProperty(&currUser.Age, newUser.Age)
	updateProperty(&currUser.Address.Street, newUser.Address.Street)
	updateProperty(&currUser.Address.City, newUser.Address.City)
	updateProperty(&currUser.Address.State, newUser.Address.State)
	updateProperty(&currUser.Address.Country, newUser.Address.Country)
}

func updateProperty[T any](oldValue *T, newValue T) {
	ov := reflect.ValueOf(oldValue)
	nv := reflect.ValueOf(newValue)

	switch reflect.TypeOf(newValue).String() {
	case "int":
		if !ov.Equal(nv) {
			*oldValue = newValue
		}
	case "string":
		if !nv.Equal(reflect.ValueOf("")) && !ov.Equal(nv) {
			*oldValue = newValue
		}
	}
}
