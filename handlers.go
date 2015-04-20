package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func saveUserToSessionAndSendHome(w http.ResponseWriter, r *http.Request,
	session *sessions.Session, userName string, id int) {

	session.Values["UserName"] = userName
	session.Values["Id"] = id
	err := session.Save(r, w)
	if err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func mainPage(w http.ResponseWriter, r *http.Request) {

	user, err := GetUser(r)
	if err != nil {
		log.Fatal(err)
	}

	friends, err := GetFrends(r)
	if err != nil {
		log.Fatal(err)
	}

	templates["index"].Execute(w, struct {
		User_c      User
		Num_freinds int
		Friends     []User
	}{User_c: user, Num_freinds: len(friends), Friends: friends})
}

func logout(w http.ResponseWriter, r *http.Request) {

	session, err := GetSession(r)
	if err != nil {
		log.Fatal(err)
	}

	saveUserToSessionAndSendHome(w, r, session, "", -1)
}

func loginGET(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/login.html")
}

func (g *Global) loginPOST(w http.ResponseWriter, r *http.Request) {
	// TODO: make errors flash messages
	// TODO: make all communication with server AJAX or similar

	session, err := GetSession(r)
	if err != nil {
		log.Fatal(err)
	}

	userName := r.FormValue("UserName")

	var password string
	var id int

	err = g.db.QueryRow("SELECT id, Password FROM Users WHERE UserName = ?",
		userName).Scan(&id, &password)
	switch {
	case err == sql.ErrNoRows:
		fmt.Fprintf(w, "No user with that UserName.")
	case err != nil:
		log.Fatal(err)
	case bcrypt.CompareHashAndPassword([]byte(password),
		[]byte(r.FormValue("Password"))) != nil:
		fmt.Fprintf(w, "Incorrect password")
	default:
		saveUserToSessionAndSendHome(w, r, session, userName, id)
	}
}

func friendGET(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/friend.html")
}

func (g *Global) friendPOST(w http.ResponseWriter, r *http.Request) {
	// TODO: make sure that you don't friend someone multiple times
	// TODO: make sure you don't freind yourself

	friendUserName := r.FormValue("UserName")

	user, err := GetUser(r)
	if err != nil {
		log.Fatal(err)
	}

	var friendId int

	// begin transaction
	tx, err := g.db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// find friends id
	err = tx.QueryRow("SELECT id FROM Users WHERE UserName = ?",
		friendUserName).Scan(&friendId)
	switch {
	case err == sql.ErrNoRows:
		fmt.Fprintf(w, "No user with that UserName.")
		err = tx.Rollback()
		if err != nil {
			log.Fatal(err)
		}
		return
	case err != nil:
		log.Fatal(err)
	}

	// create friendship in database
	_, err = tx.Exec("INSERT INTO Friends (user_id_1, user_id_2) VALUES(?,?)",
		user.Id, friendId)
	if err != nil {
		log.Fatal(err)
	}

	_, err = tx.Exec("INSERT INTO Friends (user_id_1, user_id_2) VALUES(?,?)",
		friendId, user.Id)
	if err != nil {
		log.Fatal(err)
	}

	// commit transaction
	tx.Commit()

	http.Redirect(w, r, "/", http.StatusFound)
}

func registerGET(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/register.html")
}

func (g *Global) registerPOST(w http.ResponseWriter, r *http.Request) {

	session, err := GetSession(r)
	if err != nil {
		log.Fatal(err)
	}

	userName := r.FormValue("UserName")
	password := r.FormValue("Password")

	// generate hashed password
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Fatal(err)
	}

	res, err := g.db.Exec("INSERT INTO Users (UserName, Password) VALUES(?,?)",
		userName, passwordHash)
	if err != nil {
		fmt.Fprintf(w, "UserName Alredy taken")
	} else {
		id, err := res.LastInsertId()
		if err != nil {
			log.Fatal(err)
		} else {
			saveUserToSessionAndSendHome(w, r, session, userName, int(id))
		}
	}
}

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

func (g *Global) wsHandler(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	user, err := GetUser(r)
	if err != nil {
		log.Fatal(err)
	}

	recvCh := make(chan string)
	g.ec.BindRecvChan("hello", recvCh)
	sendCh := make(chan string)
	g.ec.BindSendChan("hello", sendCh)

	c := &connection{sendCh: sendCh, recvCh: recvCh, ws: ws, user: user}
	go c.writer()
	c.reader()
}

func play(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/play.html")
}
