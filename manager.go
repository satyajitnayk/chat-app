package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		CheckOrigin:     checkOrigin,
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Manger struct {
	clients ClientList
	sync.RWMutex
	handlers map[string]EventHandler
}

func NewManger() *Manger {
	m := &Manger{
		clients:  make(ClientList), // to avoid null ptr exception
		handlers: make(map[string]EventHandler),
	}

	m.setupEventHandlers()
	return m
}

func (m *Manger) setupEventHandlers() {
	m.handlers[EventSendMessage] = SendMessage
}

func SendMessage(event Event, client *Client) error {
	fmt.Println(event)
	return nil
}

func (m *Manger) routeEvent(event Event, client *Client) error {
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, client); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("there is no such event type")
	}
}

func (m *Manger) serverWS(w http.ResponseWriter, r *http.Request) {
	log.Println("new connection")
	// upgrade regular HTTP connection to websocket
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(conn, m)

	m.addClient(client)

	// start client process
	go client.readMessages()
	go client.writeMessages()

}

func (m *Manger) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	m.clients[client] = true
}

func (m *Manger) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		client.connection.Close() // close client connection
		delete(m.clients, client)
	}
}

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")

	// TODO:make it configurable from env/config file
	switch origin {
	case "http://localhost:8080":
		return true
	default:
		return false
	}
}
