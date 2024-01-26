package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	manager    *Manger

	// egress is used to avoid concurrent writes  on the WS connection
	egress chan Event
}

func NewClient(conn *websocket.Conn, manager *Manger) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
	}
}

func (client *Client) readMessages() {
	// clean unused connection
	defer func() {
		client.manager.removeClient(client)
	}()

	for {
		// refer to Message types in websocket section of readme for messageTypes
		_, payload, err := client.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
				// TODO: Handle this type of error
			}
			break
		}

		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("error marshalling event: %v", err)
			break
		}

		if err := client.manager.routeEvent(request, client); err != nil {
			log.Printf("error handeling message: %v", err)
		}
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

			data, err := json.Marshal(message)
			if err != nil {
				// TODO: Handle error
				log.Println(err)
				return
			}

			if err := client.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("failed to send message: %v", err)
			}
			log.Println("message sent")
		}
	}
}
