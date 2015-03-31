package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type User struct {
	UserName string
	Id       int
}

func mainPage(w http.ResponseWriter, r *http.Request, ctx *Context, g *Global) {
	g.t.Execute(w, ctx.User)
}

func login(w http.ResponseWriter, r *http.Request, ctx *Context, g *Global) {
	var Password string
	var id int

	UserName := r.FormValue("UserName")

	err := g.db.QueryRow("SELECT id, Password FROM Users WHERE UserName = ?", UserName).Scan(&id, &Password)
	switch {
	case err == sql.ErrNoRows:
		fmt.Fprintf(w, "No user with that UserName.")
	case err != nil:
		log.Fatal(err)
	case bcrypt.CompareHashAndPassword([]byte(Password), []byte(r.FormValue("Password"))) != nil:
		fmt.Fprintf(w, "Incorrect password")
	default:
		ctx.session.Values["UserName"] = UserName
		ctx.session.Values["id"] = id
		err := ctx.session.Save(r, w)
		if err != nil {
			log.Fatal(err)
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
