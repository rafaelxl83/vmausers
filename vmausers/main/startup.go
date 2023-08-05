package main

import (
	"flag"
	"fmt"
	"os"
	"vmausers/helper"

	"go.mongodb.org/mongo-driver/mongo"
)

func StartDatabase(configFile string) (*mongo.Client, error) {
	file := flag.String("config", configFile, "a config file")
	flag.Parse()

	_, err := os.Stat(*file)
	if err != nil {
		fmt.Printf("File not found: %v", err)
		return nil, err
	}

	helper.DBConfig, err = helper.LoadConfig(*file)
	if err != nil || len(helper.DBConfig.Mongodb.Serveruri) == 0 {
		fmt.Printf("Error opening the config file: %v", err)
		return nil, err
	}

	client, err := helper.NewConnection(&helper.DBConfig)
	if err != nil {
		fmt.Printf("Error connecting to the database: %v", err)
		return nil, err
	}

	return client, nil
}
