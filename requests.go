package main

import (
	"fmt"
	"net/http"

	"encoding/json"

	"github.com/go-chi/chi"
)

func requestsRoute(w http.ResponseWriter, r *http.Request) {

	url := chi.URLParam(r, "url")

	db, _ := connectToSQL()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM requests WHERE url = ? ORDER BY time ASC", url)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	results := []request{}
	request := request{}

	// Make an array of requests for the template
	for rows.Next() {
		var time int64
		var url, method, ip, post, headers, body, referer string

		rows.Scan(&url, &time, &method, &ip, &post, &headers, &body, &referer)

		request.URL = url
		request.Time = time
		request.Method = method
		request.IP = ip
		request.Post = post
		request.Headers = headers
		request.Body = body
		request.Referer = referer

		results = append(results, request)
	}

	resultsByteArray, err := json.Marshal(results)

	vars := requestTemplateVars{}
	vars.Requests = string(resultsByteArray)
	vars.Domain = r.Host
	vars.URL = url

	returnTemplate(w, "requests", vars)
}

func clearRoute(w http.ResponseWriter, r *http.Request) {

	url := chi.URLParam(r, "url")

	db, _ := connectToSQL()
	defer db.Close()

	_, err := db.Query("DELETE FROM requests WHERE url = ?", url)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	http.Redirect(w, r, "/"+url+"/list", http.StatusTemporaryRedirect)
	return
}

type requestTemplateVars struct {
	Requests string
	Domain   string
	URL      string
}
