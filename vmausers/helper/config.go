package helper

import (
	"encoding/json"
	"os"
)

type Config struct {
	Mongodb struct {
		Serveruri     string `json:"serveruri"`
		ReplicaSet    string `json:"replicaSet"`
		CaFilePath    string `json:"caFilePath"`
		CaKeyFilePath string `json:"certificateKeyFilePath"`
	} `json:"mongodb"`
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
