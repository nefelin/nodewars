package nwmessage

import "client"

// ClientMessage is used to hold incoming messages from Players attached to a pointer to the player object
type ClientMessage struct {
	Type   string
	Sender client.Client
	Data   string
}

func MsgFromClient(c *client.Client) (ClientMessage, error) {
	var msg ClientMessage
	msg.Sender = c

	err := c.Socket.ReadJSON(&msg)

	return msg, err
}
