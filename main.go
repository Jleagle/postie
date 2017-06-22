package main

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

var webSocketConnections map[string]*websocket.Conn

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

func getWebSocket(w http.ResponseWriter, r *http.Request) *websocket.Conn {

	// Check if map is initialized
	if webSocketConnections == nil {
		webSocketConnections = make(map[string]*websocket.Conn)
	}

	// Check if we already have a connection with this key
	url := chi.URLParam(r, "url")
	conn, ok := webSocketConnections[url]

	if !ok {

		if r.Header.Get("Origin") != "http://"+r.Host {
			http.Error(w, "Origin not allowed", 403)
		}

		conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
		if err != nil {
			http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		}

		webSocketConnections[url] = conn
	}

	return conn
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
