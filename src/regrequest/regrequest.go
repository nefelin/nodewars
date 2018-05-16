package regrequest

import (
	"nwmodel"

	"github.com/gorilla/websocket"
)

type Request struct {
	Action  Action
	Ws      *websocket.Conn
	ResChan chan *nwmodel.Player
}

func Reg(ws *websocket.Conn, c chan *nwmodel.Player) Request {
	return Request{Register, ws, c}
}

func Dereg(ws *websocket.Conn, c chan *nwmodel.Player) Request {
	return Request{Deregister, ws, c}
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
