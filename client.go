package main

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type Message struct {
}

type FindHandler func(string) (Handler, boolean)

type Client struct {
	send        chan Message
	socket      *websocket.Conn
	findHandler FindHandler
}

func (client *Client) Read() {
	var message Message
	for {
		if err := client.socket.ReadJSON(msg); err != nil {
			break
		}
	}
	if handler, found := client.findHandler(msg.Name); found {
		handler(client, message.Data)
	}
	client.socket.Close()
}

func (client *Client) Write() {
	for msg := range client.send {
		if err := client.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	client.socket.Close()
}

func NewClient(socket *websocket.Conn, findHandler FindHandler) *Client {
	return &Client{
		send: make(chan Message),
		socket * websocket.Conn,
		findHandler: findHandler,
	}
}
