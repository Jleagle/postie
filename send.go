package main

import (
	"net/http"

	"github.com/andelf/go-curl"
)

func sendRoute(w http.ResponseWriter, r *http.Request) {
	returnTemplate(w, "send", nil)
}

func postSendRoute(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	easy := curl.EasyInit()
	defer easy.Cleanup()

	easy.Setopt(curl.OPT_URL, "http://google.com")

	// make a callback function
	fooTest := func(buf []byte, userdata interface{}) bool {
		println("DEBUG: size=>", len(buf))
		println("DEBUG: content=>", string(buf))
		return true
	}

	easy.Setopt(curl.OPT_WRITEFUNCTION, fooTest)

	if err := easy.Perform(); err != nil {
		Error(err)
	}
}
