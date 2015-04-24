package main

import (
	"github.com/gorilla/websocket"
)

type Message struct {
	Origin  User
	Content string
}

type connection struct {
	// The websocket connection.
	ws     *websocket.Conn
	sendCh chan<- interface{}
	recvCh <-chan interface{}
	user   User
}

func (c *connection) reader() {
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}

		c.sendCh <- Message{Origin: c.user, Content: string(message)}
	}
	c.sendCh <- Message{Origin: c.user, Content: "CLOSE"}
	c.ws.Close()
}

func (c *connection) writer() {
	for message := range c.recvCh {

		err := c.ws.WriteJSON(message)
		if err != nil {
			break
		}
	}
	c.ws.Close()
}
