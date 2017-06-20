package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func requestsRoute(w http.ResponseWriter, r *http.Request) {
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

	t, err := template.ParseGlob("templates/*.html")
	if err != nil {
		panic(err)
	}

	err = t.ExecuteTemplate(w, "requests", results)
	if err != nil {
		panic(err)
	}
}
