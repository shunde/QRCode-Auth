package main

import (
	"html/template"
	"log"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.WriteHeader()
	}
}

func main() {
	http.HandleFunc("/", login)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}
