package regrequest

import (
	"nwmodel"
)

type Request struct {
	Action  Action
	Player  *nwmodel.Player
	ResChan chan bool
}

func Reg(p *nwmodel.Player, c chan bool) Request {
	return Request{Register, p, c}
}

func Dereg(p *nwmodel.Player) Request {
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
