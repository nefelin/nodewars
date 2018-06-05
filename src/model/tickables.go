package model

import (
	"fmt"
	"nwmessage"
	"time"
)

// this is a naive approach, would be more performant to deal in deltas and only use tick to increment total, not recalculate rate
// TODO approach should be that on any module gain or loss that teams procPow is recalculated
// this entails making a pool of all nodes connected to POE and running the below logic
// DANGER this breaks our concurrency protections, if anything else touches coincoin we are no long safe and already there is a small risk of collisions around game status...
func (gm *GameModel) scoreTick(e time.Duration) {
	if gm.running {
		winners := make([]string, 0)

		for _, team := range gm.Teams {
			// gm.updateCoinPerTick(team) Don't think we need this if we update coinPerTick on any relevant change, IE machine gain/loss/reset
			team.CoinCoin += team.coinPerTick
			if team.CoinCoin >= gm.PointGoal {
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
	if gm.clock < 0 {
		gm.clock++

		switch {
		case gm.clock == 0:
			gm.running = true
			gm.psBroadcast(nwmessage.PsAlert(fmt.Sprintf("Game has started!")))

		case gm.clock == -5:
			gm.psBroadcast(nwmessage.PsAlert(fmt.Sprintf("Game starting in %d\n", gm.clock*-1)))

		case gm.clock > -5:
			gm.psBroadcast(nwmessage.PsAlert(fmt.Sprintf("%d\n", gm.clock*-1)))
		}

	} else if gm.running {
		gm.clock++
		// do alerts if game has a timelimit
	}
}
