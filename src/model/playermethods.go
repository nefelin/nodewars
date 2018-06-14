package model

import (
	"errors"
	"fmt"
	"model/machines"
	"model/node"
	"model/player"
	"nwmessage"
)

func (gm *GameModel) PlayerLocation(p *player.Player) *node.Node {
	r := gm.routes[p]
	if r == nil {
		return nil
	}
	return r.Endpoint()
}

func (gm *GameModel) CanSubmit(p *player.Player) error {

	mac := gm.CurrentMachine(p)
	switch {
	case p.Editor() == "":
		return errors.New("No code to submit")
	case mac == nil:
		return errors.New("Not attached to a machine")
	case !mac.AcceptsLanguageFrom(p, p.Language()):
		return fmt.Errorf("Machine solutions is written in %s, your solution must also be written in %[1]s", mac.Solution.Language)
	}
	return nil

}

func (gm *GameModel) BreakConnection(p *player.Player, forced bool) {
	if gm.routes[p] == nil {
		return
	}

	gm.detachPlayer(p)
	gm.routes[p] = nil

	if forced {
		p.Outgoing(nwmessage.PsError(errors.New("Connection interrupted!")))

	}
}

// TODO refactor this, modify how slots are tracked, probably with IDs
func (gm *GameModel) CurrentMachine(p *player.Player) *machines.Machine {
	r := gm.routes[p]
	if r == nil || p.MacAddress() == "" {
		return nil
	}

	return r.Endpoint().MacAt(p.MacAddress())
}
