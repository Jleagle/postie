package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi"
)

func requestsRoute(w http.ResponseWriter, r *http.Request) {
	// todo, check for 404s
	// http.NotFound(w, r)

	url := chi.URLParam(r, "url")

	db, _ := connectToSQL()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM requests WHERE url = ? ORDER BY id ASC", url)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	results := []request{}
	request := request{}

	// Make an array of requests for the template
	for rows.Next() {
		var id, time int
		var url, method, ip, post, headers, body string

		rows.Scan(&id, &url, &time, &method, &ip, &post, &headers, &body)

		request.ID = id
		request.URL = url
		request.Time = time
		request.Method = method
		request.IP = ip
		request.Post = post
		request.Headers = headers
		request.Body = body

		results = append(results, request)
	}

	t, err := template.ParseGlob("templates/*.html")
	if err != nil {
		panic(err)
	}

	vars := requestTemplateVars{}
	vars.Requests = results

	err = t.ExecuteTemplate(w, "requests", vars)
	if err != nil {
		panic(err)
	}
}

type requestTemplateVars struct {
	Requests []request
}
