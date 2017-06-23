package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

func endpointRoute(w http.ResponseWriter, r *http.Request) {

	// Save request to MySQL
	db, _ := connectToSQL()
	defer db.Close()

	url := chi.URLParam(r, "url")

	_, queryError := db.Query("INSERT INTO requests (time, url, method, ip, post, headers, body) VALUES (?, ?, ?, ?, ?, ?, ?)",
		time.Now().Unix(), url, r.Method, r.RemoteAddr, "c.FormParams", "string(headers)", "body")

	if queryError != nil {
		fmt.Println(queryError)
	}

	// Check if there are websockets to send to
	for _, webSocket := range webSockets {
		if webSocket.key == url {
			m := request{}
			m.IP = "66"

			// fmt.Printf("TEST: %v", conn)

			err := webSocket.connection.WriteJSON(m)
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

// https://medium.com/doing-things-right/pretty-printing-http-requests-in-golang-a918d5aaa000
// // formatRequest generates ascii representation of a request
// func formatRequest(r *http.Request) string {
// 	// Create return string
// 	var request []string
// 	// Add the request string
// 	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
// 	request = append(request, url)
// 	// Add the host
// 	request = append(request, fmt.Sprintf("Host: %v", r.Host))
// 	// Loop through headers
// 	for name, headers := range r.Header {
// 		name = strings.ToLower(name)
// 		for _, h := range headers {
// 			request = append(request, fmt.Sprintf("%v: %v", name, h))
// 		}
// 	}

// 	// If this is a POST, add post data
// 	if r.Method == "POST" {
// 		r.ParseForm()
// 		request = append(request, "\n")
// 		request = append(request, r.Form.Encode())
// 	}
// 	// Return the request as a string
// 	return strings.Join(request, "\n")
// }
