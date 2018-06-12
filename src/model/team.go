package model

import (
	"errors"
	"feature"
	"fmt"
	"model/node"
	"model/player"
	"nwmessage"
)

type teamName = string

type team struct {
	Name        string  `json:"name"` // Names are only colors for now
	CoinCoin    float32 `json:"cc"`
	players     map[*player.Player]bool
	maxSize     int                 //`json:"maxSize"`
	poes        map[*node.Node]bool // point of entry, the place where all team.players connect to the map through
	powered     []*node.Node        // list of nodes connected ot the poe, optimization to minimize re-calculating which nodes are feeding processing power
	coinPerTick float32             // stored current coint production so we don't need to recalculate every tick
}

// initializer:
// NewTeam creates a new team with color/name color
func NewTeam(n teamName) *team {
	return &team{
		Name:    n,
		players: make(map[*player.Player]bool),
		maxSize: 100,
		powered: make([]*node.Node, 0),
		poes:    make(map[*node.Node]bool),
	}
}

// team methods -------------------------------------------------------------------------------

func (t team) isFull() bool {
	if len(t.players) < t.maxSize {
		return false
	}
	return true
}

func (t *team) broadcast(msg nwmessage.Message) {
	msg.Sender = "pseudoServer"

	for player := range t.players {
		player.Outgoing(msg)

	}
}

func (t *team) addPlayer(p *player.Player) {
	t.players[p] = true
	p.TeamName = t.Name
}

func (t *team) removePlayer(p *player.Player) {
	delete(t.players, p)
	p.TeamName = ""
	p.Outgoing(nwmessage.TeamState(""))

}

func (t *team) addPoe(n *node.Node) error {
	if n.Feature.Type != feature.POE {
		return fmt.Errorf("No Point of Entry feature at Node, '%d'", n.ID)
	}

	if !n.Feature.BelongsTo(t.Name) {
		return errors.New("Team can only route to poes where it controls the feature")
	}

	// fmt.Printf("%s team adding poe\n", t.Name)

	// set the teams poe
	t.poes[n] = true
	// fmt.Printf("%s's poes: %v\n", t.Name, t.poes)
	return nil
}

func (t *team) remPoe(n *node.Node) error {
	if _, ok := t.poes[n]; !ok {
		return fmt.Errorf("%s team's poes do not include node %d", t.Name, n.ID)
	}
	//
	// fmt.Printf("%s team removing poe\n", t.Name)
	//
	delete(t.poes, n)
	// fmt.Printf("%s's poes: %v\n", t.Name, t.poes)
	return nil
}
