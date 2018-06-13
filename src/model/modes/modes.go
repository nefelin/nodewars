package model

type Mode interface {
	implementsType()
}

type gameMode int

func (f gameMode) implementsType() {}

const (
	// pre-start modes
	Map gameMode = iota
	AwaitingPlayers
	Countdown

	// post-start modes
	Running
	Over
)
