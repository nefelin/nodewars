package nwmodel

import (
	"nwmessage"
)

type teamName = string

type team struct {
	Name      string  `json:"name"` // Names are only colors for now
	ProcPow   float32 `json:"procPow"`
	VicPoints float32 `json:"vicPoints"`
	players   map[*Player]bool
	maxSize   int            //`json:"maxSize"`
	poe       *node          // point of entry, the place where all team.players connect to the map through
	powered   map[*node]bool // list of nodes connected ot the poe, optimization to minimize re-calculating which nodes are feeding processing power
}

// initializer:
// NewTeam creates a new team with color/name color
func NewTeam(n teamName) *team {
	return &team{
		Name:    n,
		players: make(map[*Player]bool),
		maxSize: 100,
		powered: make(map[*node]bool),
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
		player.Outgoing <- msg
	}
}

func (t *team) addPlayer(p *Player) {
	t.players[p] = true
	p.TeamName = t.Name
}

func (t *team) removePlayer(p *Player) {
	delete(t.players, p)
	p.TeamName = ""
	p.Outgoing <- nwmessage.TeamState("")
}
