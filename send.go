package main

import (
	"html/template"
	"net/http"
)

func sendRoute(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("templates/header.html", "templates/footer.html", "templates/send.html")
	if err != nil {
		panic(err)
	}

	err = t.ExecuteTemplate(w, "send", nil)
	if err != nil {
		panic(err)
	}
}
