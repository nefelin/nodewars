package model

import (
	"challenges"
	"errors"
	"fmt"
	"model/player"
	"nwmessage"
)

func (gm *GameModel) PlayerLocation(p *player.Player) *node {
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
	case !mac.isNeutral() && !mac.belongsTo(p.TeamName) && mac.language != p.Language():
		return fmt.Errorf("This machine is written in %s, your code must also be written in %s", mac.language, mac.language)
	}
	return nil

}

func (gm *GameModel) MacDetach(p *player.Player) {

	mac := gm.CurrentMachine(p)
	p.SetChallenge(challenges.Challenge{})

	if mac != nil {
		mac.remPlayer(p)
		p.SetMacAddress("")
	}
}

func (gm *GameModel) BreakConnection(p *player.Player, forced bool) {
	if gm.routes[p] == nil {
		return
	}

	gm.MacDetach(p)
	gm.routes[p] = nil

	if forced {
		p.Outgoing(nwmessage.PsError(errors.New("Connection interrupted!")))

	}
}

// TODO refactor this, modify how slots are tracked, probably with IDs
func (gm *GameModel) CurrentMachine(p *player.Player) *machine {
	r := gm.routes[p]
	if r == nil || p.MacAddress() == "" {
		return nil
	}

	return r.Endpoint().addressMap[p.MacAddress()]
}
