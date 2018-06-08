package room

import (
	"model/player"
)

type Room interface {
	Name() string
	Type() Type // TODO switch this to a roomtype definition
	//	Recv(msg nwmessage.ClientMessage) error
	AddPlayer(p *player.Player) error
	RemovePlayer(p *player.Player) error
	GetPlayers() []*player.Player
}

type RoomManager interface {
	Rooms() []Room
	RemovePlayer(p *player.Player, r Room)
	PlacePlayer(p *player.Player, r Room)
}

type Type interface {
	implementsType()
}

type roomType string

func (r roomType) implementsType() {}

const (
	Game  roomType = "Game"
	Lobby roomType = "Lobby"
)
