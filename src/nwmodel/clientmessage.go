package nwmodel

// ClientMessage is used to hold incoming messages from Players attached to a pointer to the player object
type ClientMessage struct {
	Type   string
	Sender *Player
	Data   string
}

func MsgFromPlayer(p *Player) (ClientMessage, error) {
	var msg ClientMessage
	msg.Sender = p

	err := p.Socket.ReadJSON(&msg)

	return msg, err
}
