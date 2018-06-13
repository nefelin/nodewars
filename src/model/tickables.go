package model

import (
	"fmt"
	"model/modes"
	"nwmessage"
	"time"
)

// type tickFunc func(elapsed time.Duration)

// this is a naive approach, would be more performant to deal in deltas and only use tick to increment total, not recalculate rate
// TODO approach should be that on any module gain or loss that teams procPow is recalculated
// this entails making a pool of all nodes connected to POE and running the below logic
// DANGER this breaks our concurrency protections, if anything else touches coincoin we are no long safe and already there is a small risk of collisions around game status...
func (gm *GameModel) scoreTick(e time.Duration) {
	if gm.mode == modes.Running {
		winners := make([]string, 0)

		for _, team := range gm.Teams {
			// gm.updateCoinPerTick(team) Don't think we need this if we update coinPerTick on any relevant change, IE machine gain/loss/reset
			team.CoinCoin += team.coinPerTick
			if team.CoinCoin >= gm.options.coinGoal {
				winners = append(winners, team.Name)
			}
			// gm.psBroadcast(nwmessage.PsAlert(fmt.Sprintf("Team %s has completed %d calculations", team.Name, team.CoinCoin)))
		}

		if len(winners) > 0 {
			for _, name := range winners {
				gm.psBroadcast(nwmessage.PsAlert(fmt.Sprintf("Team %s wins!", name)))
			}
			gm.stopGame()
		}

		gm.broadcastScore()
	}
}

func (gm *GameModel) gameClock(e time.Duration) {
	switch gm.mode {
	case modes.Countdown:
		gm.clock++

		switch {
		case gm.clock == -5:
			gm.psBroadcast(nwmessage.PsAlert(fmt.Sprintf("Game starting in %d\n", gm.clock*-1)))

		case gm.clock > -5:
			gm.psBroadcast(nwmessage.PsAlert(fmt.Sprintf("%d\n", gm.clock*-1)))

		case gm.clock >= 0:
			gm.mode = modes.Running
			gm.psBroadcast(nwmessage.PsAlert(fmt.Sprintf("Game has started!")))
		}

	case modes.Running:
		gm.clock++
		// do alerts if there's a timelimit
	}

}
