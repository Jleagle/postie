package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"

	"net/http/httputil"

	"github.com/elgs/gostrgen"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	e.Renderer = t

	e.GET("/", homeRoute)
	e.GET("/new", newRoute)
	e.Any("/:url", endpointRoute)
	e.GET("/:url/x", requestsRoute)

	e.Logger.Fatal(e.Start(":1323"))
}

func homeRoute(c echo.Context) error {
	return c.Render(http.StatusOK, "home", "var :)")
}

func newRoute(c echo.Context) error {

	db, _ := connectToSQL()
	defer db.Close()

	for {
		randomString, err := gostrgen.RandGen(10, gostrgen.Upper|gostrgen.Digit, "", "")
		if err != nil {
			fmt.Println(err)
		}

		insert, err := db.Query("INSERT INTO urls (url) VALUES (?)", randomString)

		if err == nil {
			defer insert.Close()
			return c.Redirect(http.StatusTemporaryRedirect, "/"+randomString)
		}

		if sqlerr, ok := err.(*mysql.MySQLError); ok {
			if sqlerr.Number == 1062 { // Duplicate entry
				continue
			}
		}

		panic(err.Error())
	}
}

func endpointRoute(c echo.Context) error {

	db, _ := connectToSQL()
	defer db.Close()

	headers, _ := httputil.DumpRequest(c.Request(), false)

	// bolB, _ := json.Marshal(true)
	// fmt.Println(string(bolB))

	_, queryError := db.Query("INSERT INTO requests (time, url, method, ip, post, headers, body) VALUES (?, ?, ?, ?, ?, ?, ?)",
		time.Now().Unix(), c.Param("url"), c.Scheme(), c.RealIP(), "c.FormParams", string(headers), "body")

	if queryError != nil {
		fmt.Println(queryError)
	}

	return c.String(http.StatusOK, "OK")
}

func requestsRoute(c echo.Context) error {
	// todo, check for 404s

	// url := c.Path
	// fmt.Println(url.)

	db, _ := connectToSQL()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM requests ORDER BY id ASC")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	results := []request{}
	request := request{}

	for rows.Next() {
		var id, time int
		var url, method, ip, post, headers, body string

		rows.Scan(&id, &url, &time, &method, &ip, &post, &headers, &body)

		request.id = id
		request.url = url
		request.time = time
		request.method = method
		request.ip = ip
		request.post = post
		request.headers = headers
		request.body = body

		results = append(results, request)
	}

	return c.Render(http.StatusOK, "requests", results)
}

// Template xxx
type Template struct {
	templates *template.Template
}

// Render xxx
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
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
