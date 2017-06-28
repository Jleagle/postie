package main

import (
	"database/sql"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"encoding/json"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

var webSockets []webSocket

func main() {

	r := chi.NewRouter()

	r.Get("/", homeRoute)
	r.Get("/info", infoRoute)
	r.Get("/new", newRoute)
	r.Get("/send", sendRoute)
	r.Get("/{url}/list", requestsRoute)
	r.Get("/{url}/ws", webSocketRoute)
	r.Get("/{url}/clear", clearRoute)
	r.HandleFunc("/{url}", endpointRoute)

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "assets")
	fileServer(r, "/assets", http.Dir(filesDir))

	http.ListenAndServe(":7070", r)
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func fileServer(r chi.Router, path string, root http.FileSystem) {

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
	URL     string `json:"url"`
	Time    int64  `json:"time"`
	Method  string `json:"method"`
	IP      string `json:"ip"`
	Post    string `json:"post"`
	Headers string `json:"headers"`
	Body    string `json:"body"`
	Referer string `json:"referer"`
}

func (r request) GetInfo() string {
	bytes, _ := json.Marshal(map[string]interface{}{"HTTP Method": r.Method, "IP": r.IP, "Time": r.Time, "Referer": r.Referer, "Body": r.Body})
	return string(bytes)
}

// url is the database row
type url struct {
	id  int
	url string
}

// webSocket is how we store websockets
type webSocket struct {
	key        string
	time       int64
	connection *websocket.Conn
}
