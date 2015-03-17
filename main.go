package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

type User struct {
	userName string
	id       int
}

func main() {
	r := mux.NewRouter()

	// Subrouter for POSTed requests
	s := r.Methods("POST").Subrouter()

	s.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "login POST "+r.FormValue("UserName")+" "+r.FormValue("Password"))
	})

	t := template.Must(template.ParseFiles("templates/main.tmpl", "templates/header.tmpl", "templates/footer.tmpl"))

	r.PathPrefix("/static/stylesheets/").Handler(http.StripPrefix("/static/stylesheets/", http.FileServer(http.Dir("static/stylesheets"))))
	r.PathPrefix("/static/images/").Handler(http.StripPrefix("/static/images/", http.FileServer(http.Dir("static/images"))))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t.Execute(w, nil)
	})

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/html/login.html")
	})

	fmt.Println("serving...")
	http.ListenAndServe(":80", r)
}
