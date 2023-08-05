package test

import (
	_ "fmt" // no more error
	"testing"
	"vmausers/helper"
)

func TestClientCreateUser(t *testing.T) {
	helper.AppConfig = Config

	// Create a new user
	err := TestUser.BaseModel.CreateUser(&TestUser)
	_ = err

	AssertEqual(t, err, nil)
}
