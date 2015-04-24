// TODO: generate uuid for all users (probably store in database). To allow for
// one to one rather than one to many messages.

package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

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

	recvCh := make(chan interface{})
	defer close(recvCh)
	g.ec.BindRecvChan("hello", recvCh)
	sendCh := make(chan interface{})
	defer close(sendCh)
	g.ec.BindSendChan("hello", sendCh)

	c := &connection{sendCh: sendCh, recvCh: recvCh, ws: ws, user: user}
	go c.writer()
	c.reader()
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

	// listen to messages from all friends
	recvCh := make(chan interface{})
	defer close(recvCh)
	for _, friend := range friends {
		sub, err := g.ec.Subscribe(string(friend.Id), func(v interface{}) {
			recvCh <- v
		})
		if err != nil {
			log.Fatal(err)
		}
		defer sub.Unsubscribe()
	}

	// create bind send channel
	sendCh := make(chan interface{})
	defer close(sendCh)
	g.ec.BindSendChan(string(user.Id), sendCh)

	sendCh <- Message{Origin: user, Content: "PING"}

	c := &connection{sendCh: sendCh, recvCh: recvCh, ws: ws, user: user}

	go c.writer()
	c.reader()
}
