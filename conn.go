package main

import (
	"github.com/gorilla/websocket"
)

type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	sendCh chan string
	recvCh chan string
	user   User
}

func (c *connection) reader() {
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}

		c.sendCh <- c.user.UserName + ": " + string(message)
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
