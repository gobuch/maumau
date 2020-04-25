package main

import (
	"log"

	"github.com/gorilla/websocket"
)

// client ist used for handling the websocket
// connection to a browser instance.
type client struct {
	socket   *websocket.Conn
	messages chan []byte
	playerID string
}

// write receives messages over a channel and write that message into the
// websocke connection. That method runs as long the message-channel is open
// and the connection exists. That methos blocks the program-flow, so it
// must run inside a own goroutine.
func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.messages {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("error websocket client.write:", err)
			return
		}
	}
}
