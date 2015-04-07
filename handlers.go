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

func login(w http.ResponseWriter, r *http.Request, ctx *Context, g *Global) {
	// TODO: make errors flash messages
	// TODO: make all communication with server AJAX or similar
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

func friendHandler(w http.ResponseWriter, r *http.Request,
	// TODO: make sure that you don't friend someone multiple times and make
	// sure you don't friend yourself.
	// TODO: store friends in session
	ctx *Context, g *Global) {

	UserName := r.FormValue("UserName")

	var friendId int

	// begin transaction
	tx, err := g.db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// find friends id
	err = tx.QueryRow("SELECT id FROM Users WHERE UserName = ?",
		UserName).Scan(&friendId)
	switch {
	case err == sql.ErrNoRows:
		fmt.Fprintf(w, "No user with that UserName.")
	case err != nil:
		log.Fatal(err)
	}

	// create friendship in database
	_, err = tx.Exec("INSERT INTO Friends (user_id_1, user_id_2) VALUES(?,?)",
		ctx.User.Id, friendId)
	if err != nil {
		log.Fatal(err)
	}

	_, err = tx.Exec("INSERT INTO Friends (user_id_1, user_id_2) VALUES(?,?)",
		friendId, ctx.User.Id)
	if err != nil {
		log.Fatal(err)
	}

	// commit transaction
	tx.Commit()

	fmt.Fprintf(w, "Friend added")
}

func registrationHandler(w http.ResponseWriter, r *http.Request,
	// TODO: Redirect to main page and automatic login
	ctx *Context, g *Global) {

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(r.FormValue("Password")), bcrypt.MinCost)
	if err != nil {
		log.Fatal(err)
	}
	_, err = g.db.Exec("INSERT INTO Users (UserName, Password) VALUES(?,?)",
		r.FormValue("UserName"), passwordHash)
	if err != nil {
		fmt.Fprintf(w, "UserName Alredy taken")
	} else {
		fmt.Fprintf(w, "You are Registerd")
	}
}

func mainPage(w http.ResponseWriter, r *http.Request, ctx *Context, g *Global) {

	friends := make([]User, 0, 10)

	if ctx.User.Id != -1 {

		rows, err := g.db.Query("SELECT user_id_2, Username FROM "+
			"Friends INNER JOIN Users "+
			"ON Friends.user_id_2=Users.id "+
			"WHERE user_id_1 = ?",
			ctx.User.Id)
		if err != nil {
			log.Fatal(err)
		}
		for rows.Next() {
			var id int
			var UserName string
			err = rows.Scan(&id, &UserName)
			if err != nil {
				log.Fatal(err)
			}
			friends = append(friends, User{UserName: UserName, Id: id})
		}
		rows.Close()
	}

	g.t.ExecuteTemplate(w, "main.tmpl",
		struct {
			User_c      User
			Num_freinds int
			Friends     []User
		}{User_c: ctx.User, Num_freinds: len(friends), Friends: friends})
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
