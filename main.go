package main

import (
	"database/sql"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

var webSockets []webSocket

func main() {

	r := chi.NewRouter()

	r.Get("/", homeRoute)
	r.Get("/info", infoRoute)
	r.Get("/new", newRoute)
	r.Get("/{url}/x", requestsRoute)
	r.Get("/{url}/ws", webSocketRoute)
	r.HandleFunc("/{url}", endpointRoute)

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "assets")
	FileServer(r, "/assets", http.Dir(filesDir))

	http.ListenAndServe(":8080", r)
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
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
	Time    int64  `json:"time"`
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
