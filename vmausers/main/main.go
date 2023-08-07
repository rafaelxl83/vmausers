package main

import (
	"context"
	"os"
)

func main() {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	client, err := StartDatabase(path + "\\config.json")
	if err != nil {
		panic(err)
	}

	client.Disconnect(context.Background())

	router := StartRouter()
	router.Run(":8080")

}
