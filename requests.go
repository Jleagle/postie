package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/websocket"
)

func requestsRoute(w http.ResponseWriter, r *http.Request) {
	// todo, check for 404s
	// http.NotFound(w, r)

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
		request.IP = ip
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

func wsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Origin") != "http://"+r.Host {
		http.Error(w, "Origin not allowed", 403)
		return
	}

	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	go echo(conn)
}

func echo(conn *websocket.Conn) {
	for {
		m := request{}

		err := conn.ReadJSON(&m)
		if err != nil {
			fmt.Println("Error reading json.", err)
		}

		fmt.Printf("Got message: %#v\n", m)

		err = conn.WriteJSON(m)
		if err != nil {
			fmt.Println(err)
		}
	}
}

type requestTemplateVars struct {
	requests []request
}
