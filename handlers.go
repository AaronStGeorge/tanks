package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type User struct {
	UserName string
	Id       int
}

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

func wsHandler(w http.ResponseWriter, r *http.Request, ctx *Context, g *Global) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	recvCh := make(chan string)
	g.ec.BindRecvChan("hello", recvCh)
	sendCh := make(chan string)
	g.ec.BindSendChan("hello", sendCh)

	c := &connection{sendCh: sendCh, recvCh: recvCh, ws: ws, ctx: ctx}
	go c.writer()
	c.reader()
}

func mainPage(w http.ResponseWriter, r *http.Request, ctx *Context, g *Global) {
	g.t.ExecuteTemplate(w, "main.tmpl", ctx.User)
}

func logout(w http.ResponseWriter, r *http.Request, ctx *Context, g *Global) {
	ctx.session.Values["id"] = -1
	err := ctx.session.Save(r, w)
	if err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func play(w http.ResponseWriter, r *http.Request, ctx *Context, g *Global) {
	g.t.ExecuteTemplate(w, "home.html", r.Host)
}

func login(w http.ResponseWriter, r *http.Request, ctx *Context, g *Global) {
	var Password string
	var id int

	UserName := r.FormValue("UserName")

	err := g.db.QueryRow("SELECT id, Password FROM Users WHERE UserName = ?",
		UserName).Scan(&id, &Password)
	switch {
	case err == sql.ErrNoRows:
		fmt.Fprintf(w, "No user with that UserName.")
	case err != nil:
		log.Fatal(err)
	case bcrypt.CompareHashAndPassword([]byte(Password),
		[]byte(r.FormValue("Password"))) != nil:
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
