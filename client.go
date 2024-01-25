package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	manager    *Manger

	// egress is used to avoid concurrent writes  on the WS connection
	egress chan []byte
}

func NewClient(conn *websocket.Conn, manager *Manger) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan []byte),
	}
}

func (client *Client) readMessages() {
	// clean unused connection
	defer func() {
		client.manager.removeClient(client)
	}()

	for {
		// refer to Message types in websocket section of readme for messageTypes
		messageType, payload, err := client.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
				// TODO: Handle this type of error
			}
			break
		}

		// heck to check if egress working or not
		for wsclient := range client.manager.clients {
			wsclient.egress <- payload
		}

		log.Println(messageType)
		log.Println(string(payload))
	}
}

func (client *Client) writeMessages() {
	// clean unused connection
	defer func() {
		client.manager.removeClient(client)
	}()

	for {
		select {
		case message, ok := <-client.egress:
			// check if egress channel is closed
			if !ok {
				if err := client.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed: ", err)
				}
				return // break for loop
			}

			if err := client.connection.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("failed to send message: %v", err)
			}
			log.Println("message sent")
		}
	}
}
