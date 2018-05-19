package room

import "nwmodel"
import "nwmessage"

type Room interface {
	Name() string
	Type() string // TODO switch this to a roomtype definition
	Recv(msg nwmessage.ClientMessage) error
	AddPlayer(p *nwmodel.Player) error
	RemovePlayer(p *nwmodel.Player) error
	GetPlayers() []*nwmodel.Player
}
