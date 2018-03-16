package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

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
