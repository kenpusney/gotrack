package main

import (
	"fmt"
	"net/http"
	"text/template"
)

type kaJSModel struct {
	ID   string
	Host string
}

func setupScript() {
	http.HandleFunc("/ka.js", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		host := r.Host
		ID := query.Get("id")
		model := kaJSModel{ID: ID, Host: host}
		tmpl, err := template.New("ka.js").ParseFiles("ka.js.tmpl")
		if err == nil {
			w.Header().Add("Content-Type", "text/javascript")
			w.WriteHeader(200)
			tmpl.Templates()[0].Execute(w, model)
		} else {
			responseMessage(w, http.StatusNotFound, fmt.Sprint(err))
		}
	})
}
