package main

import "net/http"
import "github.com/labstack/echo"

func main() {
	e := echo.New()

	e.GET("/", home)
	e.GET("/new", new)
	e.GET("/:id", requests)
	e.GET("/:id/x", requests)

	e.Logger.Fatal(e.Start(":1323"))
}

func home(c echo.Context) error {
	return c.String(http.StatusOK, "Home")
}

func new(c echo.Context) error {
	return c.String(http.StatusOK, "New")
}

func requests(c echo.Context) error {
	return c.String(http.StatusOK, "Requests")
}
