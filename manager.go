package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Manger struct {
	clients ClientList
	sync.RWMutex
}

func NewManger() *Manger {
	return &Manger{
		clients: make(ClientList), // to avoid null ptr exception
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

	// defer conn.Close()
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
