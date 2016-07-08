package main

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type stopChan chan bool

type FindHandler func(string) (Handler, bool)

type Client struct {
	send        chan Message
	socket      *websocket.Conn
	findHandler FindHandler
	stopChannel stopChan
}

func (client *Client) Read() {
	var msg Message
	for {
		fmt.Println("reading")
		if err := client.socket.ReadJSON(&msg); err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(msg)
		if handler, found := client.findHandler(msg.Name); found {
			handler(client, msg.Data)
		}
	}
	client.socket.Close()
}

func (client *Client) Write() {
	for msg := range client.send {
		fmt.Println(msg)
		if err := client.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	client.socket.Close()
}

func NewClient(socket *websocket.Conn, findHandler FindHandler) *Client {
	return &Client{
		send:        make(chan Message),
		socket:      socket,
		findHandler: findHandler,
		stopChannel: make(stopChan),
	}
}
