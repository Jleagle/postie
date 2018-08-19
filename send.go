package main

import (
	"net/http"
)

func sendRoute(w http.ResponseWriter, r *http.Request) {
	returnTemplate(w, "send", nil)
}

func postSendRoute(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

}
