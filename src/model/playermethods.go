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
	case !mac.IsNeutral() && !mac.BelongsTo(p.TeamName) && mac.Language != p.Language():
		return fmt.Errorf("This machine is written in %s, your code must also be written in %s", mac.Language, mac.Language)
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
