package nwmessage

import "github.com/gorilla/websocket"

type Client interface {
	Outgoing(Message)
	ChatMode() bool
	Socket() *websocket.Conn
}
