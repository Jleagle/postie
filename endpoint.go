package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pressly/chi"
)

func endpointRoute(w http.ResponseWriter, r *http.Request) {

	db, _ := connectToSQL()
	defer db.Close()

	// headers, _ := httputil.DumpRequest(c.Request(), false)

	// bolB, _ := json.Marshal(true)
	// fmt.Println(string(bolB))

	_, queryError := db.Query("INSERT INTO requests (time, url, method, ip, post, headers, body) VALUES (?, ?, ?, ?, ?, ?, ?)",
		time.Now().Unix(), chi.URLParam(r, "url"), r.Method, r.RemoteAddr, "c.FormParams", "string(headers)", "body")

	if queryError != nil {
		fmt.Println(queryError)
	}

	w.Write([]byte("OK"))
}

// formatRequest generates ascii representation of a request
func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}
