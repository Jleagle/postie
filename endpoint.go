package main

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

var webSockets map[string]map[int]*websocket.Conn

func endpointRoute(w http.ResponseWriter, r *http.Request) {

	url := chi.URLParam(r, "url")

	match, err := regexp.MatchString("^[A-Z0-9]{10}$", url)
	if err != nil {
		returnErrorTemplate(w, err)
		return
	}
	if !match {
		http.NotFound(w, r)
		return
	}

	r.ParseForm()

	// Gather data
	headers, err := json.Marshal(r.Header)
	if err != nil {
		returnErrorTemplate(w, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		returnErrorTemplate(w, err)
		return
	}

	form, err := json.Marshal(r.Form)
	if err != nil {
		returnErrorTemplate(w, err)
		return
	}

	// Save request to MySQL
	db, err := connectToSQL()
	if err != nil {
		returnErrorTemplate(w, err)
		return
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
		returnErrorTemplate(w, err)
		return
	}

	// Check if there are websockets to send to
	val, ok := webSockets[url]
	if ok {
		for k, webSocket := range val {
			err := webSocket.WriteJSON(request)
			if err != nil {
				if strings.Contains(err.Error(), "broken pipe") {
					webSocket.Close()
					delete(webSockets[url], k)
					if len(webSockets[url]) < 1 {
						delete(webSockets, url)
					}
				} else {
					Error(err)
				}
			}
		}
	}

	// Return
	w.Write([]byte("OK"))
}

func webSocketRoute(w http.ResponseWriter, r *http.Request) {

	key := chi.URLParam(r, "url")
	ran := rand.Int()

	// Initialize maps
	if webSockets == nil {
		webSockets = map[string]map[int]*websocket.Conn{}
	}

	_, ok := webSockets[key]
	if !ok {
		webSockets[key] = map[int]*websocket.Conn{}
	}

	// Make a connection
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		Error(err)
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	webSockets[key][ran] = conn
}
