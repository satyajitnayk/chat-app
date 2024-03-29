package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Maximum time to wait for client response.
	pongWait = 10 * time.Second

	// Interval to send ping messages.
	// Must be shorter than pongWait to maintain connection.
	// Otherwise, connections will be terminated.
	pingInterval = (pongWait * 9) / 10
)

// ClientList is a map used to help manage a map of clients
type ClientList map[*Client]bool

// Client is a websocket client, basically a frontend visitor
type Client struct {
	connection *websocket.Conn
	manager    *Manager
	chatroom   string
	// egress is used to avoid concurrent writes on the WS connection
	egress chan Event
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
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

	// set the wait time for pong message from client
	if err := client.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
		return
	}

	// set read limit to message
	// if client send more than limit server will close the connection
	client.connection.SetReadLimit(512) // set byte size as per requirement

	client.connection.SetPongHandler(client.pongHandler)

	for {
		// ReadMessage is used to read the next message in queue
		// in the connection
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

// writeMessages is a process that listens for new messages to output to the Client
func (client *Client) writeMessages() {
	// clean unused connection
	defer func() {
		client.manager.removeClient(client)
	}()

	ticker := time.NewTicker(pingInterval)

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

		case <-ticker.C:
			log.Println("ping")
			//send a ping to client
			if err := client.connection.WriteMessage(websocket.PingMessage, []byte(``)); err != nil {
				log.Println("write message error: ", err)
				return
			}
		}

	}
}

// pongHandler is used to handle PongMessages for the Client
func (client *Client) pongHandler(pongMsg string) error {
	log.Println("pong")
	// reset the wait timer for pong msg from client
	return client.connection.SetReadDeadline(time.Now().Add(pongWait))
}
