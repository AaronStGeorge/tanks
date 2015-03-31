package main

import (
	"github.com/apcera/nats"
	"github.com/gorilla/websocket"
	"net/http"
)

type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	sendCh chan string
	recvCh chan string
}

func (c *connection) reader() {
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}

		c.sendCh <- string(message)
	}
	c.ws.Close()
}

func (c *connection) writer() {
	for message := range c.recvCh {
		err := c.ws.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			break
		}
	}
	c.ws.Close()
}

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

type wsHandler struct {
	ec *nats.EncodedConn
}

func (wsh wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	recvCh := make(chan string)
	wsh.ec.BindRecvChan("hello", recvCh)
	sendCh := make(chan string)
	wsh.ec.BindSendChan("hello", sendCh)

	c := &connection{sendCh: sendCh, recvCh: recvCh, ws: ws}
	go c.writer()
	c.reader()
}
