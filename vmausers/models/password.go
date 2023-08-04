package models

import (
	"time"
	"vmausers/helper"
)

const key string = "PhotOp0lym3rRXL4"

type Password struct {
	EncryptedPass string    `bson:"password"`
	CreatedAt     time.Time `bson:"created_at"`
	Expire        time.Time `bson:"expire"`
}

func NewPassword(Key string) *Password {
	keyStr := key
	if Key != "" {
		keyStr = Key
	}

	cryptoText := helper.Encrypt(keyStr, Key)
	pass := Password{cryptoText, time.Now(), time.Now()}
	return &pass
}
