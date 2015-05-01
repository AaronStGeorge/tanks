package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 10 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

func ErrNoPubTo(m Message) error {
	return errors.New(fmt.Sprintf("Don't know where to send this %+v", m))
}

type Message struct {
	Origin  User
	PubTo   string
	Content interface{}
}

type connection struct {
	// The websocket connection.
	ws     *websocket.Conn
	toPage <-chan Message
	user   User
	g      *Global
}

func (c *connection) reader() {
	defer func() {
		c.g.ec.Publish(c.user.Twitter, Message{Origin: c.user,
			Content: "CLOSE"})
		c.ws.Close()
	}()

	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, data, err := c.ws.ReadMessage()
		if err != nil {
			return
		}

		var message Message
		err = json.Unmarshal(data, &message)
		if err != nil {
			log.Fatal(err)
		}
		if message.PubTo == "" {
			log.Fatal(ErrNoPubTo(message))
		}

		// publish to Nats message hub
		c.g.ec.Publish(message.PubTo, message)
	}
}

func (c *connection) writer() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.toPage:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.writeJSON(message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (c *connection) write(mt int, message []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, message)
}

func (c *connection) writeJSON(message Message) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteJSON(message)
}
