package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/Jleagle/go-helpers/rollbar"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

func endpointRoute(w http.ResponseWriter, r *http.Request) {

	url := chi.URLParam(r, "url")

	match, err := regexp.MatchString("^[A-Z0-9]{10}$", url)
	if err != nil {
		rollbar.ErrorCritical(err)
	}
	if !match {
		http.NotFound(w, r)
		return
	}

	r.ParseForm()

	// Gather data
	headers, err := json.Marshal(r.Header)
	if err != nil {
		rollbar.ErrorCritical(err)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rollbar.ErrorCritical(err)
	}

	form, err := json.Marshal(r.Form)
	if err != nil {
		rollbar.ErrorCritical(err)
	}

	// Save request to MySQL
	db, err := connectToSQL()
	if err != nil {
		rollbar.ErrorCritical(err)
	}

	defer db.Close()

	request := request{}
	request.Time = time.Now().UnixNano()
	request.URL = url
	request.Method = r.Method
	request.IP = r.Header.Get("X-Forwarded-For")
	request.Post = string(form)
	request.Headers = string(headers)
	request.Body = string(body)
	request.Referer = r.Referer()

	_, queryError := db.Query("INSERT INTO requests (time, url, method, ip, post, headers, body, referer) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		request.Time, request.URL, request.Method, request.IP, request.Post, request.Headers, request.Body, request.Referer)

	if queryError != nil {
		rollbar.ErrorCritical(err)
	}

	// Check if there are websockets to send to
	for _, webSocket := range webSockets {
		if webSocket.key == url {
			err := webSocket.connection.WriteJSON(request)
			if err != nil {
				rollbar.ErrorCritical(err)
			}
		}
	}

	// Return
	w.Write([]byte("OK"))
}

func webSocketRoute(w http.ResponseWriter, r *http.Request) {

	// Initialized slice
	if webSockets == nil {
		webSockets = []webSocket{}
	}

	// Make a connection
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		rollbar.ErrorCritical(err)
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	// Save the connection
	newSocket := webSocket{}
	newSocket.connection = conn
	newSocket.time = time.Now().Unix()
	newSocket.key = chi.URLParam(r, "url")

	webSockets = append(webSockets, newSocket)
}
