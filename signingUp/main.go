package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var tpl *template.Template

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 PAGE NOT", http.StatusNotFound)
		return
	}

	if r.Method == "GET" {
		http.ServeFile(w, r, "index.html")
	}

}

func login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {

		http.Error(w, "404 PAGE NOT FOUND", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "login.html")
	default:
		fmt.Fprintf(w, "Only GET and POST")

	}

}
func registration(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "registration.html", nil)
}

func registrationAuth(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fullname := r.FormValue("fullname")
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	str := fmt.Sprintf("Thanks %s, for your registration your Account:: %s , Email:: %s , Password:: %s", fullname, username, email, password)
	tpl.ExecuteTemplate(w, "registration.html", str)

}

func main() {
	tpl, _ = template.ParseGlob("*.html")

	http.HandleFunc("/registration", registration)
	http.HandleFunc("/registrationauth", registrationAuth)
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)

	fmt.Printf("Starting server got testing\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
