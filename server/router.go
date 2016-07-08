package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Router struct {
	rules map[string]Handler
}

type Handler func(*Client, interface{})

func NewRouter() *Router {
	return &Router{
		rules: make(map[string]Handler),
	}
}

func (r *Router) Handle(msgName string, handler Handler) {
	r.rules[msgName] = handler
}

func (e *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	client := NewClient(socket, e.FindHandler)
	go client.Write()
	client.Read()
}

func (r *Router) FindHandler(msg string) (Handler, bool) {
	handler, found := r.rules[msg]
	return handler, found
}
