package main

import (
	"flag"
	"github.com/apcera/nats"
	"log"
	"net/http"
	"text/template"
)

var homeTempl *template.Template

func homeHandler(c http.ResponseWriter, req *http.Request) {
	homeTempl.Execute(c, req.Host)
}

func main() {
	flag.Parse()
	homeTempl = template.Must(template.ParseFiles("home.html"))

	nc, _ := nats.Connect(nats.DefaultURL)
	ec, _ := nats.NewEncodedConn(nc, "json")
	defer ec.Close()

	http.HandleFunc("/", homeHandler)
	http.Handle("/ws", wsHandler{ec: ec})
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
