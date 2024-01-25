package main

import "github.com/gorilla/websocket"

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	manager    *Manger
}

func NewClient(conn *websocket.Conn, manager *Manger) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
	}
}
