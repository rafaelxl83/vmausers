package helper

import (
	"encoding/json"
	"os"
)

var AppConfig Config

type Config struct {
	Mongodb struct {
		Serveruri     string `json:"serveruri"`
		ReplicaSet    string `json:"replicaSet"`
		Database      string `json:"database"`
		CaFilePath    string `json:"caFilePath"`
		CaKeyFilePath string `json:"certificateKeyFilePath"`
	} `json:"mongodb"`
	PasswordStrength struct {
		MinSize          int  `json:"minSize"`
		MustSpecialChars bool `json:"mustSpecialChars"`
		MustNumeric      bool `json:"mustNumeric"`
		MustLowerUpper   bool `json:"mustLowerUpper"`
	} `json:"passwordStrength"`
}

func NewConfig(Uri string, Replica string, Database string, caFilePath string, caKeyFile string) *Config {
	config := Config{}
	config.Mongodb.Serveruri = Uri
	config.Mongodb.ReplicaSet = Replica
	config.Mongodb.Database = Database
	config.Mongodb.CaFilePath = caFilePath
	config.Mongodb.CaKeyFilePath = caKeyFile

	/// password default
	config.PasswordStrength.MinSize = 8
	config.PasswordStrength.MustLowerUpper = true
	config.PasswordStrength.MustNumeric = true
	config.PasswordStrength.MustSpecialChars = true

	return &config
}

func LoadConfig(configFile string) (Config, error) {
	f, err := os.ReadFile(configFile)
	if err != nil {
		return Config{}, err
	}

	data := Config{}
	err = json.Unmarshal([]byte(f), &data)
	if err != nil {
		return Config{}, err
	}

	return data, nil
}
