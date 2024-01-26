package main

import "encoding/json"

type Event struct {
	Type string `json:"type"`
	// this to allow user to put any type of data in payload
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, client *Client) error

const (
	EventSendMessage = "send_message"
)

type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
}
