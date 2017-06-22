package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/elgs/gostrgen"
	"github.com/go-sql-driver/mysql"
)

func homeRoute(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseGlob("templates/*.html")
	if err != nil {
		panic(err)
	}

	err = t.ExecuteTemplate(w, "home", nil)
	if err != nil {
		panic(err)
	}
}

func newRoute(w http.ResponseWriter, r *http.Request) {

	db, _ := connectToSQL()
	defer db.Close()

	for {
		randomString, err := gostrgen.RandGen(10, gostrgen.Upper|gostrgen.Digit, "", "")
		if err != nil {
			fmt.Println(err)
		}

		insert, err := db.Query("INSERT INTO urls (url) VALUES (?)", randomString)

		if err == nil {
			defer insert.Close()
			http.Redirect(w, r, "/"+randomString+"/x", http.StatusTemporaryRedirect)
			return
		}

		if sqlerr, ok := err.(*mysql.MySQLError); ok {
			//fmt.Println(sqlerr.Duplicate) // todo, use this?
			if sqlerr.Number == 1062 { // Duplicate entry
				continue
			}
		}

		panic(err.Error())
	}
}
