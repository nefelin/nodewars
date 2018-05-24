package regrequest

import (
	"model/player"
)

type Request struct {
	Action  Action
	Player  *player.Player
	ResChan chan bool
}

func Reg(p *player.Player, c chan bool) Request {
	return Request{Register, p, c}
}

func Dereg(p *player.Player) Request {
	return Request{Deregister, p, nil}
}

// Type ...
type Action interface {
	implementsType()
}

type regAction int

func (f regAction) implementsType() {}

const (
	Register regAction = iota
	Deregister
)
