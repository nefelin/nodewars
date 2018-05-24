package nwmodel

import (
	"errors"
	"feature"
	"fmt"
	"nwmessage"
	"nwmodel/player"
)

type teamName = string

type team struct {
	Name        string  `json:"name"` // Names are only colors for now
	VicPoints   float32 `json:"vicPoints"`
	players     map[*player.Player]bool
	maxSize     int            //`json:"maxSize"`
	poes        map[*node]bool // point of entry, the place where all team.players connect to the map through
	powered     []*node        // list of nodes connected ot the poe, optimization to minimize re-calculating which nodes are feeding processing power
	coinPerTick float32        // stored current coint production so we don't need to recalculate every tick
}

// initializer:
// NewTeam creates a new team with color/name color
func NewTeam(n teamName) *team {
	return &team{
		Name:    n,
		players: make(map[*player.Player]bool),
		maxSize: 100,
		powered: make([]*node, 0),
		poes:    make(map[*node]bool),
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

func (t *team) addPoe(n *node) error {
	if n.Feature.Type != feature.POE {
		return fmt.Errorf("No Point of Entry feature at Node, '%d'", n.ID)
	}

	if !n.Feature.belongsTo(t.Name) {
		return errors.New("Team can only route to poes where it controls the feature")
	}

	fmt.Printf("%s team adding poe\n", t.Name)

	// set the teams poe
	t.poes[n] = true
	fmt.Printf("%s's poes: %v\n", t.Name, t.poes)
	return nil
}

func (t *team) remPoe(n *node) error {
	if _, ok := t.poes[n]; !ok {
		return fmt.Errorf("%s team's poes do not include node %d", t.Name, n.ID)
	}

	fmt.Printf("%s team removing poe\n", t.Name)

	delete(t.poes, n)
	fmt.Printf("%s's poes: %v\n", t.Name, t.poes)
	return nil
}
