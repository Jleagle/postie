package main

import (
	"fmt"
	"html/template"
	"net/http"

	curl "github.com/andelf/go-curl"
)

func sendRoute(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("templates/header.html", "templates/footer.html", "templates/send.html")
	if err != nil {
		panic(err)
	}

	err = t.ExecuteTemplate(w, "send", nil)
	if err != nil {
		panic(err)
	}
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
		fmt.Printf("ERROR: %v\n", err)
	}

}
