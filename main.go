package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

type entry struct {
	ID      string
	Headers interface{}
	Params  interface{}
	Time    time.Time
}

type queryDetail struct {
	Limit uint
	Skip  uint
}

type result struct {
	Size  uint
	Query queryDetail
	Data  []entry
}

type messageDetail struct {
	Code    int
	Message string
}

func responseMessage(w http.ResponseWriter, status int, msg string) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	data, _ := json.Marshal(messageDetail{Code: status, Message: msg})
	w.Write(data)
}

func getIntParam(query url.Values, key string, d int) uint {
	if l := query.Get(key); l != "" {
		result, err := strconv.Atoi(l)
		if err == nil {
			return uint(result)
		}
	}
	return uint(d)
}

func validAgent(config Config, agentID string) bool {
	for _, a := range config.AcceptedIDs {
		if a == agentID {
			return true
		}
	}
	return false
}

func main() {
	db, err := bolt.Open("data.db", 0600, nil)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("items"))
		return err
	})

	config, _ := LoadConfig("conf.json")

	http.HandleFunc("/ka.php", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		if agentID := query.Get("id"); validAgent(config, agentID) {
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
		newConfig, err := LoadConfig(query.Get("config"))
		if err != nil {
			responseMessage(w, http.StatusInternalServerError, fmt.Sprintf("%s", err))
			return
		}
		config = newConfig
		responseMessage(w, http.StatusAccepted, "config reloaded")
	})

	http.HandleFunc("/report.php", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		username, password, ok := r.BasicAuth()

		if !ok || (username != config.User && password != config.Pass) {
			responseMessage(w, http.StatusUnauthorized, "invalid auth info")
			return
		}

		limit := getIntParam(query, "limit", 20)
		skip := getIntParam(query, "skip", 0)
		res := result{
			Query: queryDetail{
				Skip:  skip,
				Limit: limit,
			},
			Size: 0,
		}
		db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte("items"))
			cursor := bucket.Cursor()
			for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
				if skip > 0 {
					skip--
					continue
				}
				if limit > 0 {
					limit--
					e := entry{}
					json.Unmarshal(v, &e)
					res.Data = append(res.Data, e)
					res.Size++
				} else {
					break
				}
			}
			return nil
		})
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		data, _ := json.Marshal(res)
		w.Write(data)
	})

	http.ListenAndServe(":10086", nil)
	select {}
}
