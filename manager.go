package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		CheckOrigin:     checkOrigin,
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Manager struct {
	clients ClientList
	sync.RWMutex

	otps     RetentionMap
	handlers map[string]EventHandler
}

func NewManager(ctx context.Context) *Manager {
	m := &Manager{
		clients:  make(ClientList), // to avoid null ptr exception
		handlers: make(map[string]EventHandler),
		otps:     NewRetentionMap(ctx, 5*time.Second),
	}

	m.setupEventHandlers()
	return m
}

func (m *Manager) setupEventHandlers() {
	m.handlers[EventSendMessage] = SendMessageHandler
	m.handlers[EventChangeChatRoom] = ChangeChatRoomHandler
}

func ChangeChatRoomHandler(event Event, client *Client) error {
	var changeChatRoomEvent ChangeChatRoomEvent

	if err := json.Unmarshal(event.Payload, &changeChatRoomEvent); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	client.chatroom = changeChatRoomEvent.Name
	return nil
}

// SendMessageHandler will send out a message to all other participants in the chat
func SendMessageHandler(event Event, client *Client) error {
	var chatevent SendMessageEvent
	if err := json.Unmarshal(event.Payload, &chatevent); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	// broadcast message to all clients
	var broadcastMessage NewMessageEvent
	broadcastMessage.Message = chatevent.Message
	broadcastMessage.From = chatevent.From

	data, err := json.Marshal(broadcastMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	outgoingEvent := Event{
		Payload: data,
		Type:    EventReceiveMessage,
	}

	for currentClient := range client.manager.clients {
		// broadcast to all clients in the same chat room
		if client.chatroom == currentClient.chatroom {
			currentClient.egress <- outgoingEvent
		}
	}

	return nil
}

func (m *Manager) routeEvent(event Event, client *Client) error {
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, client); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("there is no such event type")
	}
}

func (m *Manager) serverWS(w http.ResponseWriter, r *http.Request) {
	otp := r.URL.Query().Get("otp")

	// no otp in url or invalid otp
	if otp == "" || !m.otps.VerifyOTP(otp) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

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

func (m *Manager) loginHandler(w http.ResponseWriter, r *http.Request) {
	type userLoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req userLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// For simplicity I have hard coded them(As my focus is Websocket in this project)
	if req.Username == "satya" && req.Password == "1234" {
		type response struct {
			OTP string `json:"otp"`
		}

		otp := m.otps.NewOTP()

		resp := response{
			OTP: otp.Key,
		}

		data, err := json.Marshal(resp)
		if err != nil {
			log.Println(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	m.clients[client] = true
}

func (m *Manager) removeClient(client *Client) {
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
	case "https://localhost:8080":
		return true
	default:
		return false
	}
}
