package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

func requestsRoute(w http.ResponseWriter, r *http.Request) {

	url := chi.URLParam(r, "url")

	db, err := connectToSQL()
	if err != nil {
		returnErrorTemplate(w, err)
		return
	}

	rows, err := db.Query("SELECT * FROM requests WHERE url = ? ORDER BY time ASC", url)
	if err != nil {
		returnErrorTemplate(w, err)
		return
	}

	var results []request
	request := request{}

	// Make an array of requests for the template
	for rows.Next() {
		var time int64
		var url, method, ip, post, headers, body, referer string

		err = rows.Scan(&url, &time, &method, &ip, &post, &headers, &body, &referer)
		if err != nil {
			returnErrorTemplate(w, err)
			return
		}

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
	vars.URL = url
	vars.Domain = os.Getenv("POSTIE_DOMAIN")

	if r.Header.Get("X-Forwarded-Proto") == "https" {
		vars.Protocol = "wss"
	} else {
		vars.Protocol = "ws"
	}

	returnTemplate(w, "requests", vars)
}

type requestTemplateVars struct {
	Requests string
	Protocol string
	Domain   string
	URL      string
}

func clearRoute(w http.ResponseWriter, r *http.Request) {

	url := chi.URLParam(r, "url")

	db, err := connectToSQL()
	if err != nil {
		returnErrorTemplate(w, err)
		return
	}

	_, err = db.Query("DELETE FROM requests WHERE url = ?", url)
	if err != nil {
		returnErrorTemplate(w, err)
		return
	}

	http.Redirect(w, r, "/"+url+"/list", http.StatusTemporaryRedirect)
	return
}
