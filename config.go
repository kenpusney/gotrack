package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/boltdb/bolt"
)

// Config - Kana configurations
type Config struct {
	User        string
	Pass        string
	AcceptedIDs []string
}

func loadConfig(configFile string) (Config, error) {

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

func isValidAgent(config Config, agentID string) bool {
	for _, a := range config.AcceptedIDs {
		if a == agentID {
			return true
		}
	}
	return false
}

func setupConfigReload(config *Config, db *bolt.DB) {
	http.HandleFunc("/reload", func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		query := r.URL.Query()

		if !ok || (username != config.User && password != config.Pass) {
			responseMessage(w, http.StatusUnauthorized, "invalid auth info")
			return
		}

		if query.Get("config") == "" {
			responseMessage(w, http.StatusBadRequest, "invalid `config` parameter")
			return
		}
		newConfig, err := loadConfig(query.Get("config"))
		if err != nil {
			responseMessage(w, http.StatusInternalServerError, fmt.Sprintf("%s", err))
			return
		}
		config = &newConfig
		responseMessage(w, http.StatusAccepted, "config reloaded")
	})
}
