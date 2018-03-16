package main

import (
	"log"
	"net/http"

	"github.com/boltdb/bolt"
)

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

	config, _ := loadConfig("conf.json")

	setupReport(config, db)
	setupCollector(config, db)
	setupConfigReload(&config, db)
	setupScript()

	http.ListenAndServe(":10086", nil)
	select {}
}
