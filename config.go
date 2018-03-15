package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config - Kana configurations
type Config struct {
	User        string
	Pass        string
	AcceptedIDs []string
}

// LoadConfig - loading configurations from file
func LoadConfig(configFile string) (Config, error) {

	file, _ := os.Open(configFile)
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := Config{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("error: ", err)
		return config, err
	}
	return config, nil
}
