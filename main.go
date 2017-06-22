package main

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
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

func makeAWebSocket(w http.ResponseWriter, r *http.Request) *websocket.Conn {

	if r.Header.Get("Origin") != "http://"+r.Host {
		http.Error(w, "Origin not allowed", 403)
	}

	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	return conn
}

// request is the database row
type request struct {
	ID      int    `json:"id"`
	URL     string `json:"url"`
	Time    int
	Method  string
	IP      string `json:"ip"`
	Post    string
	Headers string
	Body    string
}

// url is the database row
type url struct {
	id  int
	url string
}
