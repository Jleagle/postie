package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	e.Renderer = t

	e.GET("/", home)
	e.GET("/new", new)
	e.GET("/:id", requests)
	e.GET("/:id/x", requests)

	e.Logger.Fatal(e.Start(":1323"))
}

func home(c echo.Context) error {
	return c.Render(http.StatusOK, "home", "James")
}

func new(c echo.Context) error {
	return c.String(http.StatusOK, "New")
}

func requests(c echo.Context) error {
	return c.String(http.StatusOK, "Requests")
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
