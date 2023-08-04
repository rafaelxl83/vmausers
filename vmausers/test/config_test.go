package test

import (
	_ "fmt" // no more error
	"testing"
	"vmausers/helper"
)

func TestConfigLoadNoFile(t *testing.T) {
	data, err := helper.LoadConfig("config.json")
	_ = err
	AssertEqual(t, len(data.Mongodb.Serveruri), 0)
}

func TestConfigLoad(t *testing.T) {
	file := "../main/config.json"
	data, err := helper.LoadConfig(file)
	_ = err
	AssertNotEqual(t, len(data.Mongodb.Serveruri), 0)
}
