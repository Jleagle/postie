package main

import (
	"database/sql"
	"net/http"

	"github.com/pressly/chi"
)

func main() {

	r := chi.NewRouter()

	r.Get("/", homeRoute)
	r.Get("/new", newRoute)
	r.Get("/{url}/x", requestsRoute)
	r.Get("/{url}/ws", wsHandler)
	r.HandleFunc("/{url}", endpointRoute)

	http.ListenAndServe(":8080", r)
}

func connectToSQL() (*sql.DB, error) {

	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/postie")
	if err != nil {
		panic(err.Error())
	}

	return db, err
}

// request is the database row
type request struct {
	id      int
	url     string
	time    int
	method  string
	ip      string
	post    string
	headers string
	body    string
}

// url is the database row
type url struct {
	id  int
	url string
}
