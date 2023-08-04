package helper

import (
	"flag"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

func Start(configFile string) (*mongo.Client, error) {
	file := flag.String("config", configFile, "a config file")
	flag.Parse()

	config, err := LoadConfig(*file)
	if err != nil || len(config.Mongodb.Serveruri) == 0 {
		fmt.Printf("Error opening the config file: %v", err)
		return nil, err
	}

	client, err := Connect(
		config.Mongodb.Serveruri,
		config.Mongodb.CaFilePath,
		config.Mongodb.CaKeyFilePath,
		config.Mongodb.ReplicaSet)
	if err != nil {
		fmt.Printf("Error connecting to the database: %v", err)
		return nil, err
	}

	return client, nil
}
