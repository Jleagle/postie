package main

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

var webSockets []webSocket

func main() {

	r := chi.NewRouter()

	r.Get("/", homeRoute)
	r.Get("/new", newRoute)
	r.Get("/{url}/x", requestsRoute)
	r.Get("/{url}/ws", webSocketRoute)
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
	ID      int    `json:"id"`
	URL     string `json:"url"`
	Time    int    `json:"time"`
	Method  string `json:"method"`
	IP      string `json:"ip"`
	Post    string `json:"post"`
	Headers string `json:"headers"`
	Body    string `json:"body"`
}

// url is the database row
type url struct {
	id  int
	url string
}

type webSocket struct {
	key        string
	time       int64
	connection *websocket.Conn
}
