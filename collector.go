package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
)

type entry struct {
	ID      string
	Headers interface{}
	Params  interface{}
	Time    time.Time
}

func setupCollector(config Config, db *bolt.DB) {
	http.HandleFunc("/ka.php", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		if agentID := query.Get("id"); isValidAgent(config, agentID) {
			now := time.Now()
			digest := sha256.Sum256([]byte(r.URL.RawQuery))
			id := fmt.Sprintf("%s/%d/%x", agentID, now.Unix(), digest[:8])
			entry := entry{
				ID: id, Params: query, Headers: r.Header, Time: now,
			}
			err := db.Update(func(tx *bolt.Tx) error {
				bucket := tx.Bucket([]byte("items"))
				data, _ := json.Marshal(entry)
				return bucket.Put([]byte(id), []byte(data))
			})
			if err != nil {
				responseMessage(w, http.StatusInternalServerError, "Save failure")
			} else {
				responseMessage(w, http.StatusCreated, "success!")
			}
		} else {
			responseMessage(w, http.StatusBadRequest, "invalid id")
		}
	})

}
