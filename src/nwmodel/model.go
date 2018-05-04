package nwmodel

import (
	"nwmessage"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

// type aliases
type nodeID = int
type modID = int
type playerID = int
type teamName = string

// incrementing ID counters
var playerIDCount playerID

var moduleIDCount modID

// GameModel holds all state information
type GameModel struct {
	Map           *nodeMap             `json:"map"`
	Teams         map[teamName]*team   `json:"teams"`
	Players       map[playerID]*Player `json:"players"`
	POEs          map[playerID]*node   `json:"poes"`
	PointGoal     float32              `json:"pointGoal"`
	languages     map[string]Language
	aChan         chan nwmessage.Message
	running       bool //running should replace mapLocked
	pendingAlerts map[playerID][]alert

	// timelimit should be able to set a timelimit and count points at the end
}

type stateMessage struct {
	*nodeMap
	Alerts    []alert `json:"alerts"`
	PlayerLoc nodeID  `json:"player_location"`
	*trafficMap
	// Traffic map[string]trafficPacket `json:"traffic"` // maps edgeIds to lists of packets on that edge
}

type alert struct {
	Actor    string `json:"team"`
	Location nodeID `json:"location"`
}

type route struct {
	Endpoint *node   `json:"endpoint"`
	Nodes    []*node `json:"nodes"`
}

// easily rendered routes traffic TODO breaks current separation of interests
type packet struct {
	Owner     string `json:"owner"`
	Direction string `json:"dir"`
}

type trafficMap struct {
	Traffic map[string][]packet `json:"traffic"`
}

// TODO route is stored reversed?
func (t *trafficMap) addRoute(r *route, color string) {
	for i, n := range r.Nodes {
		n1 := n.ID
		var n2 nodeID

		var dir, edge string

		if i == len(r.Nodes)-1 {
			n2 = r.Endpoint.ID
		} else {
			n2 = r.Nodes[i+1].ID
		}

		if n1 > n2 {
			dir = "down"
			edge = strconv.Itoa(n1) + "e" + strconv.Itoa(n2)
		} else {
			dir = "up"
			edge = strconv.Itoa(n2) + "e" + strconv.Itoa(n1)
		}
		t.appendPacket(packet{color, dir}, edge)
	}
}

func (t *trafficMap) appendPacket(p packet, edge string) {
	_, ok := t.Traffic[edge]
	if !ok {
		t.Traffic[edge] = []packet{p}
	} else {
		t.Traffic[edge] = append(t.Traffic[edge], p)
	}
}

func newTrafficMap() *trafficMap {
	return &trafficMap{
		Traffic: make(map[string][]packet),
	}
}

// func (r *route) endpoint() nodeID {
// 	return r.Nodes[len(r.Nodes)-1]
// }

type nodeMap struct {
	Nodes       []*node         `json:"nodes"`
	POEs        map[nodeID]bool `json:"poes"`
	diameter    float64
	radius      float64
	nodeIDCount nodeID
}

type node struct {
	ID          nodeID     `json:"id"` // keys and ids is redundant? TODO
	Connections []nodeID   `json:"connections"`
	Machines    []*machine `json:"machines"` // TODO why is this a list of pointerS?
	Feature     *feature   `json:"feature"`
	Remoteness  float64    //`json:"remoteness"`
	playersHere []playerID
}

type challengeCriteria struct {
	IDs        []int64  // list of acceptable challenge ids
	Tags       []string // acceptable categories of challenge
	Difficulty [][]int  // acceptable difficulties, [5] = level five, [3,5] = 3,4, or 5
}

type machine struct {
	sync.Mutex
	// accepts   challengeCriteria
	challenge Challenge
	// Type      string `json:"type"`
	Powered  bool   `json:"powered"`
	builder  string // `json:"creator"`
	TeamName string `json:"team"`
	// solution  string
	language  string // `json:"languageId"`
	Health    int    `json:"health"`
	MaxHealth int    `json:"maxHealth"`
}

func newMachine() *machine {
	return &machine{Powered: true}
}

type feature struct {
	Type string `json:"type"` // type of feature
	machine
}

func newFeature() *feature {
	return &feature{
		machine: machine{Powered: true},
	}
}

type team struct {
	Name      string  `json:"name"` // Names are only colors for now
	ProcPow   float32 `json:"procPow"`
	VicPoints float32 `json:"vicPoints"`
	players   map[*Player]bool
	maxSize   int            //`json:"maxSize"`
	poe       *node          // point of entry, the place where all team.players connect to the map through
	powered   map[*node]bool // list of nodes connected ot the poe, optimization to minimize re-calculating which nodes are feeding processing power
}

// TODO un export all but route

// Player ...
type Player struct {
	ID        playerID               `json:"id"`
	name      string                 `json:"name"`
	TeamName  string                 `json:"team"`
	Route     *route                 `json:"route"`
	Socket    *websocket.Conn        `json:"-"`
	Outgoing  chan nwmessage.Message `json:"-"`
	language  string                 // current working language
	stdin     string                 // stdin buffer for testing
	slotNum   int                    // currently attached to slotNum of current node
	dialogue  *nwmessage.Dialogue    // this holds any dialogue the players in the middle of
	compiling bool                   // this is used to block player action while submitted code is compiling
	ChatMode  bool                   // track whether player is in chatmode or not (for use in lobby)
	inGame    bool                   // is player in a game?
}
