package main

import (
	"encoding/json"
	"time"
)

type Event struct {
	Type string `json:"type"`
	// this to allow user to put any type of data in payload
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, client *Client) error

const (
	EventSendMessage    = "send_message"
	EventReceiveMessage = "receive_message"
	EventChangeChatRoom = "change_chat_room"
)

type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
}

type NewMessageEvent struct {
	SendMessageEvent
	Sent time.Time `json:"sent"`
}

type ChangeChatRoomEvent struct {
	Name string `json:"name"`
}
