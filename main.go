package main

import (
	"database/sql"
	"fmt"
	"github.com/apcera/nats"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"html/template"
	"log"
	"net/http"
)

type Global struct {
	store *sessions.CookieStore
	db    *sql.DB
	ec    *nats.EncodedConn
}

var templates = map[string]*template.Template{
	"index": template.Must(template.ParseFiles("templates/main.tmpl",
		"templates/header.tmpl", "templates/footer.tmpl")),
}

type User struct {
	UserName string
	Id       int
}

func main() {
	// create session store
	store := sessions.NewCookieStore([]byte("nRrHLlHcHH0u7fUz25Hje9m7uJ5SnJzP"))

	// store options
	store.Options = &sessions.Options{
		Path:   "/",
		Domain: "aaronstgeorge.co",
		MaxAge: 0,
	}

	// open database connection
	db, err := sql.Open("mysql", "root:@/Tanks")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// test connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// create connection to nats server
	nc, _ := nats.Connect(nats.DefaultURL)
	ec, _ := nats.NewEncodedConn(nc, "json")
	defer ec.Close()

	global := &Global{db: db, store: store, ec: ec}

	stdChain := alice.New(global.loadSession)
	loadUser := stdChain.Append(global.loadUser)
	loadUserData := loadUser.Append(global.loadFriends)

	router := httprouter.New()

	// mainPage
	router.Handler("GET", "/", loadUserData.Then(http.HandlerFunc(mainPage)))

	// login
	router.HandlerFunc("GET", "/login", loginGET)
	router.Handler("POST", "/login", stdChain.Then(
		http.HandlerFunc(global.loginPOST)))

	// logout
	router.Handler("GET", "/logout", stdChain.Then(http.HandlerFunc(logout)))

	// play
	router.HandlerFunc("GET", "/play", http.HandlerFunc(play))

	// websocket
	router.Handler("GET", "/ws", loadUser.Then(http.HandlerFunc(
		global.wsHandler)))

	// register
	router.HandlerFunc("GET", "/register", registerGET)
	router.Handler("POST", "/register", stdChain.Then(
		http.HandlerFunc(global.registerPOST)))

	// add friend
	router.HandlerFunc("GET", "/friend", friendGET)
	router.Handler("POST", "/friend", loadUser.Then(
		http.HandlerFunc(global.friendPOST)))

	// Serve static files from the ./public directory
	router.ServeFiles("/static/*filepath", http.Dir("./public/"))

	fmt.Println("serving...")
	log.Fatal(http.ListenAndServe(":80", context.ClearHandler(router)))
}
