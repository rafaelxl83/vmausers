package test

import (
	_ "fmt" // no more error
	"testing"
	"vmausers/helper"
	"vmausers/middlewares"
	"vmausers/models"
)

func TestClientCreateUser(t *testing.T) {
	helper.AppConfig = Config

	// Create a new user
	err := middlewares.CreateUser(&TestUser)
	_ = err

	AssertEqual(t, err, nil)
}

func TestClientGetUserById(t *testing.T) {
	helper.AppConfig = Config

	readUser, err := middlewares.GetUserById(TestUserId)
	if err != nil {
		_ = middlewares.CreateUser(&TestUser)
		readUser, err = middlewares.GetUserById(TestUserId)
	}
	_ = err

	AssertEqual(t, TestUser.ID, readUser.ID)
}

func TestClientGetUserByEmail(t *testing.T) {
	helper.AppConfig = Config

	readUser, err := middlewares.GetUserByEmail(TestUser.Email)
	if err != nil {
		_ = middlewares.CreateUser(&TestUser)
		readUser, err = middlewares.GetUserByEmail(TestUser.Email)
	}
	_ = err

	AssertEqual(t, TestUser.ID, readUser.ID)
}

func TestClientDeleteUser(t *testing.T) {
	helper.AppConfig = Config

	err := middlewares.DeleteUser(TestUser)
	_ = err

	AssertEqual(t, err, nil)
}

func TestClientDeleteUserById(t *testing.T) {
	helper.AppConfig = Config

	err := middlewares.DeleteUserById(TestUserId)
	_ = err

	AssertEqual(t, err, nil)
}

func TestClientDeleteUserByEmail(t *testing.T) {
	helper.AppConfig = Config

	err := middlewares.DeleteUserByEmail(TestUser.Email)
	_ = err

	AssertEqual(t, err, nil)
}

func TestClientUpdateUser(t *testing.T) {
	helper.AppConfig = Config

	update := models.NewUser(
		"Jhonny", "", 41, "john.doe@example.com",
		"Nowhere, 333", "", "", "")

	TestUser.UpdateValues(*update)
	err := middlewares.UpdateUser(&TestUser)
	_ = err

	AssertEqual(t, err, nil)
}

func TestClientUpdateUserPassword(t *testing.T) {
	helper.AppConfig = Config

	pass := models.NewPassword(MinimalPassword)
	err := middlewares.UpdateUserPassword(&TestUser, *pass)
	_ = err

	AssertEqual(t, err, nil)
}
