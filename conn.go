package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

func ErrNoPubTo(m Message) error {
	return errors.New(fmt.Sprintf("Don't know where to send this %+v", m))
}

type Message struct {
	Origin  User
	PubTo   string
	Content string
}

type connection struct {
	// The websocket connection.
	ws     *websocket.Conn
	toPage <-chan interface{}
	user   User
	g      *Global
}

func (c *connection) reader() {
	for {
		_, data, err := c.ws.ReadMessage()
		if err != nil {
			break
		}

		var message Message

		err = json.Unmarshal(data, &message)
		if err != nil {
			log.Fatal(err)
		}
		if message.PubTo == "" {
			log.Fatal(ErrNoPubTo(message))
		}

		fmt.Printf("Sent to hub %+v\n", message)

		// publish to Nats message hub
		c.g.ec.Publish(message.PubTo, message)
	}
	c.g.ec.Publish(c.user.Twitter, Message{Origin: c.user, Content: "CLOSE"})
	c.ws.Close()
}

func (c *connection) writer() {
	for message := range c.toPage {

		fmt.Printf("Received from hub %+v\n", message)

		err := c.ws.WriteJSON(message)
		if err != nil {
			break
		}
	}
	c.ws.Close()
}
