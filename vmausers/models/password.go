package models

import (
	"errors"
	"regexp"
	"strconv"
	"time"
	"vmausers/helper"
)

const key string = "PhotOp0lym3rRXL4"
const defaultPass string = "defaultpassword"

type Password struct {
	EncryptedPass string    `bson:"encrypted" binding:"required,len=8"`
	CreatedAt     time.Time `bson:"created_at"`
	Expire        time.Time `bson:"expire"`
}

func NewPassword(pass string) *Password {
	if pass == "" {
		pass = defaultPass
	}

	cryptoText := helper.Encrypt(key, pass)
	encPass := Password{cryptoText, time.Now(), time.Now().AddDate(0, 2, 0)}
	return &encPass
}

func (p *Password) CheckPassword(pass string) error {
	userPass := helper.Decrypt(key, p.EncryptedPass)

	if userPass == pass {
		return nil
	}

	return errors.New("Invalid Password")
}

func (p *Password) ValidatePasswordRestrictions(pass string) error {
	if helper.AppConfig.PasswordStrength.MinSize == 0 {
		helper.AppConfig.PasswordStrength.MinSize = 8
		helper.AppConfig.PasswordStrength.MustLowerUpper = true
		helper.AppConfig.PasswordStrength.MustNumeric = true
		helper.AppConfig.PasswordStrength.MustSpecialChars = true
	}

	if pass == "" {
		pass = p.EncryptedPass
	}

	errorMessage := ""
	if helper.AppConfig.PasswordStrength.MinSize > 0 && len(pass) < helper.AppConfig.PasswordStrength.MinSize {
		errorMessage += " minimum size is " + strconv.Itoa(helper.AppConfig.PasswordStrength.MinSize) + ";"
	}

	if helper.AppConfig.PasswordStrength.MustLowerUpper {
		er := regexp.MustCompile(`[a-z]+[A-Z]+`)
		matched := er.MatchString(pass)
		if !matched {
			errorMessage += " should have at least one lower and upper case letter;"
		}
	}

	if helper.AppConfig.PasswordStrength.MustNumeric {
		er := regexp.MustCompile(`[0-9]+`)
		matched := er.MatchString(pass)
		if !matched {
			errorMessage += " should have at least one number;"
		}
	}

	if helper.AppConfig.PasswordStrength.MustSpecialChars {
		er := regexp.MustCompile(`[!"#$%&'()*+,\-./:;<=>?@[\\\]^_{|}~]+`)
		matched := er.MatchString(pass)
		if !matched {
			errorMessage += " should have at least one special character;"
		}
	}

	if len(errorMessage) > 0 {
		return errors.New("The password does not match the security policy: " + errorMessage)
	}

	return nil
}
