package main

import (
	"net/http"

	"github.com/Jleagle/go-helpers/rollbar"
	"github.com/elgs/gostrgen"
	"github.com/go-sql-driver/mysql"
)

func homeRoute(w http.ResponseWriter, r *http.Request) {
	returnTemplate(w, "home", nil)
}

func infoRoute(w http.ResponseWriter, r *http.Request) {
	returnTemplate(w, "info", nil)
}

func newRoute(w http.ResponseWriter, r *http.Request) {

	db, err := connectToSQL()
	if err != nil {
		rollbar.ErrorCritical(err)
	}

	defer db.Close()

	for {
		randomString, err := gostrgen.RandGen(10, gostrgen.Upper|gostrgen.Digit, "", "")
		if err != nil {
			rollbar.ErrorCritical(err)
		}

		insert, err := db.Query("INSERT INTO urls (url) VALUES (?)", randomString)
		if err == nil {
			defer insert.Close()
			http.Redirect(w, r, "/"+randomString+"/list", http.StatusTemporaryRedirect)
			return
		}

		if sqlerr, ok := err.(*mysql.MySQLError); ok {
			if sqlerr.Number == 1062 { // Duplicate entry
				continue
			}
		}
		rollbar.ErrorCritical(err)
	}
}
