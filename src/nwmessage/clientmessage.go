package nwmessage

// ClientMessage is used to hold incoming messages from Players attached to a pointer to the player object
type ClientMessage struct {
	Type   string
	Sender Client
	Data   string
}

func MsgFromClient(c Client) (ClientMessage, error) {
	var msg ClientMessage
	msg.Sender = c

	err := c.Socket().ReadJSON(&msg)

	return msg, err
}
