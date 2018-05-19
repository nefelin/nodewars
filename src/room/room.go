package room

import "nwmodel"

type Room interface {
	Name() string
	Type() string // TODO switch this to a roomtype definition
	Recv(msg nwmodel.ClientMessage) error
	AddPlayer(p *nwmodel.Player) error
	RemovePlayer(p *nwmodel.Player) error
	GetPlayers() []*nwmodel.Player
}
