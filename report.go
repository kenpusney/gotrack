package main

import (
	"encoding/json"
	"net/http"

	"github.com/boltdb/bolt"
)

type queryDetail struct {
	Limit uint
	Skip  uint
}

type result struct {
	Size  uint
	Query queryDetail
	Data  []entry
}

func setupReport(config Config, db *bolt.DB) {
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
}
