// TODO: generate uuid for all users (probably store in database). To allow for
// one to one rather than one to many messages.

package main

import (
	"github.com/apcera/nats"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

func (g *Global) subFunc(s string, toPage chan Message) *nats.Subscription {
	sub, err := g.ec.Subscribe(s, func(v Message) {
		toPage <- v
	})
	if err != nil {
		log.Fatal(err)
	}
	return sub
}

func (g *Global) mainPageWs(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	defer ws.Close()
	if err != nil {
		log.Fatal(err)
	}

	user, err := GetUser(r)
	if err != nil {
		log.Fatal(err)
	}

	friends, err := GetFrends(r)
	if err != nil {
		log.Fatal(err)
	}

	toPage := make(chan Message)
	defer close(toPage)
	sub := g.subFunc(user.PhoneNumber, toPage)
	defer sub.Unsubscribe()
	for _, friend := range friends {
		sub := g.subFunc(friend.Twitter, toPage)
		defer sub.Unsubscribe()
	}

	g.ec.Publish(user.PhoneNumber, Message{Origin: user,
		PubTo: user.PhoneNumber, Content: "INIT"})

	c := &connection{g: g, toPage: toPage, ws: ws, user: user}

	go c.writer()
	c.reader()
}

func (g *Global) gamePageWs(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	defer ws.Close()
	if err != nil {
		log.Fatal(err)
	}

	user, err := GetUser(r)
	if err != nil {
		log.Fatal(err)
	}

	toPage := make(chan Message)
	sub := g.subFunc(user.PhoneNumber, toPage)
	defer sub.Unsubscribe()

	c := &connection{g: g, toPage: toPage, ws: ws, user: user}

	go c.writer()
	c.reader()
}
