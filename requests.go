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

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn := makeAWebSocket(w, r)

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
	Requests []request
}
