package test

import (
	_ "fmt" // no more error
	"testing"
	"vmausers/helper"
	"vmausers/models"
)

var MinimalPassword string = "abAB12!#"
var InvalidPasswords = [...]string{
	"123",
	"1234567890",
	"abc",
	"abcdefghij",
	"ABC",
	"ABCDEFJHIJ",
	"abcABC",
	"abcdeABCDE",
	"aA1",
	"abcABC123",
	"!@#",
	"!@#$%&+*",
	"aA1@",
}

var TestPassword models.Password = *models.NewPassword(MinimalPassword)

func TestCreatePassword(t *testing.T) {
	helper.AppConfig = Config

	pass := models.NewPassword("")

	AssertNotEqual(t, pass, nil)
}

func TestCheckPassword(t *testing.T) {
	helper.AppConfig = Config

	err := TestPassword.CheckPassword(MinimalPassword)
	_ = err

	AssertEqual(t, err, nil)
}

func TestValidateDefault(t *testing.T) {

	err := TestPassword.ValidatePasswordRestrictions(MinimalPassword)
	_ = err

	AssertEqual(t, err, nil)
}

func TestValidateMinLengh(t *testing.T) {
	helper.AppConfig = Config

	helper.AppConfig.PasswordStrength.MinSize = 8
	helper.AppConfig.PasswordStrength.MustLowerUpper = false
	helper.AppConfig.PasswordStrength.MustNumeric = false
	helper.AppConfig.PasswordStrength.MustSpecialChars = false

	err := TestPassword.ValidatePasswordRestrictions(MinimalPassword)
	_ = err

	AssertEqual(t, err, nil)
}

func TestValidateMinLenghErr(t *testing.T) {
	helper.AppConfig = Config

	helper.AppConfig.PasswordStrength.MinSize = 16
	helper.AppConfig.PasswordStrength.MustLowerUpper = false
	helper.AppConfig.PasswordStrength.MustNumeric = false
	helper.AppConfig.PasswordStrength.MustSpecialChars = false

	err := TestPassword.ValidatePasswordRestrictions(MinimalPassword)
	_ = err

	AssertNotEqual(t, err, nil)
}

func TestValidateLowerUpper(t *testing.T) {
	helper.AppConfig = Config
	helper.AppConfig.PasswordStrength.MustNumeric = false
	helper.AppConfig.PasswordStrength.MustSpecialChars = false

	err := TestPassword.ValidatePasswordRestrictions(MinimalPassword)
	_ = err

	AssertEqual(t, err, nil)
}

func TestValidateNumeric(t *testing.T) {
	helper.AppConfig = Config
	helper.AppConfig.PasswordStrength.MustLowerUpper = false
	helper.AppConfig.PasswordStrength.MustSpecialChars = false

	err := TestPassword.ValidatePasswordRestrictions(MinimalPassword)
	_ = err

	AssertEqual(t, err, nil)
}

func TestValidateSpecialChar(t *testing.T) {
	helper.AppConfig = Config
	helper.AppConfig.PasswordStrength.MustLowerUpper = false
	helper.AppConfig.PasswordStrength.MustNumeric = false

	err := TestPassword.ValidatePasswordRestrictions(MinimalPassword)
	_ = err

	AssertEqual(t, err, nil)
}

func TestDefaultValidation(t *testing.T) {
	helper.AppConfig = Config

	fails := 0
	for _, pass := range InvalidPasswords {
		err := TestPassword.ValidatePasswordRestrictions(pass)
		if err != nil {
			println("Error", pass, err.Error())
			fails++
		}
	}

	AssertEqual(t, len(InvalidPasswords), fails)
}
