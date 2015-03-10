package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func makeHandler(fn func(http.ResponseWriter, *http.Request, *sql.DB), db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, db)
	}
}

func AuthHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	var UserName, Password string

	userName := r.FormValue("UserName")

	err := db.QueryRow("SELECT UserName, Password FROM Users WHERE UserName = ?", userName).Scan(&UserName, &Password)
	switch {
	case err == sql.ErrNoRows:
		fmt.Fprintf(w, "No user with that UserName.")
	case err != nil:
		log.Fatal(err)
	case bcrypt.CompareHashAndPassword([]byte(Password), []byte(r.FormValue("Password"))) != nil:
		fmt.Fprintf(w, "Incorrect password")
	default:
		fmt.Fprintf(w, "Welcome "+UserName+"!")
	}
}

func RegistrationHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("Password")), bcrypt.MinCost)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("INSERT INTO Users (UserName, Password) VALUES(?,?)", r.FormValue("UserName"), passwordHash)
	if err != nil {
		fmt.Fprintf(w, "UserName Alredy taken")
	} else {
		fmt.Fprintf(w, "You are Registerd")
	}
}

func main() {
	// Open database connection
	db, err := sql.Open("mysql", "root:asmallpig21@/Tanks")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "login.html")
	})
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "register.html")
	})
	http.HandleFunc("/auth", makeHandler(AuthHandler, db))
	http.HandleFunc("/reg", makeHandler(RegistrationHandler, db))
	http.ListenAndServe(":80", nil)
}
