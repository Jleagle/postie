package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

func endpointRoute(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	// Gather data
	url := chi.URLParam(r, "url")
	headers, _ := json.Marshal(r.Header)
	body, _ := ioutil.ReadAll(r.Body)
	form, _ := json.Marshal(r.Form)

	// Save request to MySQL
	db, _ := connectToSQL()
	defer db.Close()

	request := request{}
	request.Time = time.Now().Unix()
	request.URL = url
	request.Method = r.Method
	request.IP = r.RemoteAddr
	request.Post = string(form)
	request.Headers = string(headers)
	request.Body = string(body)

	_, queryError := db.Query("INSERT INTO requests (time, url, method, ip, post, headers, body) VALUES (?, ?, ?, ?, ?, ?, ?)",
		request.Time, request.URL, request.Method, request.IP, request.Post, request.Headers, request.Body)

	if queryError != nil {
		fmt.Println(queryError)
	}

	// Check if there are websockets to send to
	for _, webSocket := range webSockets {
		if webSocket.key == url {
			err := webSocket.connection.WriteJSON(request)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	// Return
	fmt.Println(url)
	w.Write([]byte("OK"))
}

func webSocketRoute(w http.ResponseWriter, r *http.Request) {

	// Validation
	if r.Header.Get("Origin") != "http://"+r.Host {
		http.Error(w, "Origin not allowed", 403)
	}

	// Initialized slice
	if webSockets == nil {
		webSockets = []webSocket{}
	}

	// Make a connection
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		fmt.Printf("%s\n", err.Error())
	}

	// Save the connection
	newSocket := webSocket{}
	newSocket.connection = conn
	newSocket.time = time.Now().Unix()
	newSocket.key = chi.URLParam(r, "url")

	webSockets = append(webSockets, newSocket)
}
