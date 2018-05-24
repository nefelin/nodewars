package room

import (
	"nwmessage"
	"nwmodel/player"
)

type Room interface {
	Name() string
	Type() string // TODO switch this to a roomtype definition
	Recv(msg nwmessage.ClientMessage) error
	AddPlayer(p *player.Player) error
	RemovePlayer(p *player.Player) error
	GetPlayers() []*player.Player
}

type RoomManager interface {
	Rooms() []Room
	RemovePlayer(p *player.Player, r Room)
	PlacePlayer(p *player.Player, r Room)
}
