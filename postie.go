package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"net/http"

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
	e.Any("/:id", requestsRoute)
	e.GET("/:id/x", requestsRoute)

	e.Logger.Fatal(e.Start(":1323"))
}

func homeRoute(c echo.Context) error {
	return c.Render(http.StatusOK, "home", "xx")
}

func newRoute(c echo.Context) error {

	randomString, err := gostrgen.RandGen(10, gostrgen.Upper|gostrgen.Digit, "", "")
	if err != nil {
		fmt.Println(err)
	}

	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/postie")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	for {
		fmt.Println("loop")
		insert, err := db.Query("INSERT INTO urls (url) VALUES (?)", "test")
		defer insert.Close()

		if err == nil {
			break
		} else {
			fmt.Println("err")
			if sqlerr, ok := err.(*mysql.MySQLError); ok {
				fmt.Println("mysql err")
				if sqlerr.Number == 1062 { // Duplicate entry
					fmt.Println("1062 err")
					continue
				}
			}

			panic(err.Error())
		}
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/"+randomString)
}

func requestsRoute(c echo.Context) error {
	return c.String(http.StatusOK, "Requests")
}

// Template xxx
type Template struct {
	templates *template.Template
}

// Render xxx
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
