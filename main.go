package main

import (
	"database/sql"
	"fmt"
	"github.com/apcera/nats"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	//"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
)

type Global struct {
	store *sessions.CookieStore
	db    *sql.DB
	t     *template.Template
	ec    *nats.EncodedConn
}

type Context struct {
	User    User
	session *sessions.Session
}

type Auth struct {
	global *Global
	fn     func(w http.ResponseWriter, r *http.Request, ctx *Context, g *Global)
}

func (a *Auth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, err := a.global.store.Get(r, "tanks-app")
	if err != nil {
		log.Fatal(err)
	}

	user := User{}

	if session.IsNew {
		user.Id = -1
	} else {
		user.Id = session.Values["id"].(int)
		user.UserName = session.Values["UserName"].(string)
	}

	ctx := &Context{User: user, session: session}

	a.fn(w, r, ctx, a.global)
}

func main() {
	// create session store
	store := sessions.NewCookieStore([]byte("nRrHLlHcHH0u7fUz25Hje9m7uJ5SnJzP"))

	// store optons
	store.Options = &sessions.Options{
		Path:   "/",
		Domain: "aaronstgeorge.co",
		MaxAge: 0,
	}

	// open database connection
	db, err := sql.Open("mysql", "root:asmallpig21@/Tanks")
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

	t := template.Must(template.ParseFiles("templates/main.tmpl",
		"templates/header.tmpl", "templates/footer.tmpl",
		"templates/home.html"))

	global := &Global{db: db, store: store, ec: ec, t: t}

	r := mux.NewRouter()

	// Subrouter for POSTed requests
	s := r.Methods("POST").Subrouter()
	s.Handle("/login", &Auth{global: global, fn: login})

	// serve static files
	r.PathPrefix("/static/stylesheets/").
		Handler(http.StripPrefix("/static/stylesheets/",
		http.FileServer(http.Dir("static/stylesheets"))))

	r.PathPrefix("/static/images/").
		Handler(http.StripPrefix("/static/images/",
		http.FileServer(http.Dir("static/images"))))

	r.Handle("/", &Auth{global: global, fn: mainPage})

	r.Handle("/logout", &Auth{global: global, fn: logout})
	r.Handle("/play", &Auth{global: global, fn: play})
	r.Handle("/ws", &Auth{global: global, fn: wsHandler})

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/html/login.html")
	})

	fmt.Println("serving...")
	http.ListenAndServe(":80", r)
}
