package main

import (
	"database/sql"
	//"fmt"
	"github.com/apcera/nats"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/context"
	//"github.com/gorilla/mux"
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

	stdChain := alice.New(context.ClearHandler, global.loadUser)
	withFreinds := stdChain.Append(global.loadFriends)

	router := httprouter.New()

	// Serve static files from the ./static directory
	router.Handler("GET", "/", withFreinds.Then(http.HandlerFunc(mainPage)))
	router.ServeFiles("/public/*filepath", http.Dir("/public"))
	log.Fatal(http.ListenAndServe(":80", router))

	/*

		r := mux.NewRouter()

		// Subrouter for POSTed requests
		s := r.Methods("POST").Subrouter()
		s.Handle("/login", &Auth{global: global, fn: login})
		s.Handle("/register", &Auth{global: global, fn: registrationHandler})
		s.Handle("/friend", &Auth{global: global, fn: friendHandler})

		// serve static files
		r.PathPrefix("/static/stylesheets/").
			Handler(http.StripPrefix("/static/stylesheets/",
			http.FileServer(http.Dir("static/stylesheets"))))

		r.PathPrefix("/static/images/").
			Handler(http.StripPrefix("/static/images/",
			http.FileServer(http.Dir("static/images"))))

		r.PathPrefix("/static/js/").
			Handler(http.StripPrefix("/static/js/",
			http.FileServer(http.Dir("static/js"))))

		r.Handle("/", &Auth{global: global, fn: mainPage})
		r.Handle("/logout", &Auth{global: global, fn: logout})
		r.Handle("/ws", &Auth{global: global, fn: wsHandler})

		r.HandleFunc("/play", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "static/html/play.html")
		})

		r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "static/html/login.html")
		})

		r.HandleFunc("/friend", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "static/html/friend.html")
		})

		r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "static/html/register.html")
		})

		fmt.Println("serving...")
		http.ListenAndServe(":80", r)
	*/
}
