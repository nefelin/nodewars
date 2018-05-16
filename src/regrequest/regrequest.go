package regrequest

import (
	"nwmodel"

	"github.com/gorilla/websocket"
)

type Request struct {
	action  Action
	ws      *websocket.Conn
	retChan chan *nwmodel.Player
}

// Type ...
type Action interface {
	implementsType()
}

type regAction int

func (f regAction) implementsType() {}

type regAction string

const (
	Register regAction = iota
	Deregister
)
