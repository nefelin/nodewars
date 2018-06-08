package context

import "receiver"

type Context interface {
	receiver.Receiver
	Type() Type
}

type Type interface {
	implementsType()
}

type contextType string

func (f contextType) implementsType() {}

const (
	Game  contextType = "Game"
	Lobby             = "Lobby"
)
