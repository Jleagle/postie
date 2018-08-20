package main

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
	"github.com/rollbar/rollbar-go"
)

func main() {

	rollbar.SetToken(os.Getenv("POSTIE_ROLLBAR_PRIVATE"))
	rollbar.SetEnvironment(os.Getenv("ENV"))           // defaults to "development"
	rollbar.SetCodeVersion("dev-master")               // optional Git hash/branch/tag (required for GitHub integration)
	rollbar.SetServerRoot("github.com/Jleagle/postie") // path of project (required for GitHub integration and non-project stacktrace collapsing)

	r := chi.NewRouter()

	r.Get("/", homeRoute)
	r.Get("/info", infoRoute)
	r.Get("/new", newRoute)
	r.Get("/send", sendRoute)
	r.Post("/send", postSendRoute)
	r.Get("/{url}/list", requestsRoute)
	r.Get("/{url}/ws", webSocketRoute)
	r.Get("/{url}/clear", clearRoute)
	r.HandleFunc("/{url}", endpointRoute)

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "assets")
	fileServer(r, "/assets", http.Dir(filesDir))

	http.ListenAndServe(":8080", r)
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func fileServer(r chi.Router, path string, root http.FileSystem) {

	if strings.ContainsAny(path, "{}*") {
		rollbar.Message(rollbar.ERR, "FileServer does not permit URL parameters.")
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

var db *sql.DB

func connectToSQL() (*sql.DB, error) {

	var err error

	if db == nil {

		password := os.Getenv("SQL_PW")
		if len(password) > 0 {
			password = ":" + password
		}

		db, err = sql.Open("mysql", "root"+password+"@tcp(127.0.0.1:3306)/postie")
		if err != nil {
			Error(err)
		}
	}

	return db, err
}

func returnTemplate(w http.ResponseWriter, page string, pageData interface{}) {

	// Load templates needed
	folder := os.Getenv("POSTIE_PATH")
	if folder == "" {
		folder = "/root"
	}

	t, err := template.ParseFiles(
		folder+"/templates/header.html",
		folder+"/templates/footer.html",
		folder+"/templates/"+page+".html",
	)
	if err != nil {
		Error(err)
	}

	// Write a respone
	err = t.ExecuteTemplate(w, page, pageData)
	if err != nil {
		Error(err)
	}
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
