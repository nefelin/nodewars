package nwmodel

import (
	"nwmessage"
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
	*nodeMap          //`json:"map"`
	Alerts    []alert `json:"alerts"`
	PlayerLoc nodeID  `json:"player_location"`
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

type trafficMap struct {
	Traffic map[string][]trafficPacket `json:"traffic"`
}

type trafficPacket struct {
	owner     string `json:"owner"` // team generating this traffic
	direction string `json:"dir"`   // up/down whether this traffic is moving from a higher id node to lower or vice-versa
}

func (t *trafficMap) addRoute(r *route, color string) {
	for i := range r.Nodes {
		var n1, n2 nodeID
		var dir, edgeID string
		n1 := r.Nodes[i].ID

		if i == len(r.Nodes)-1 {
			n2 = r.Endpoint.ID
		} else {
			n2 = r.Nodes[i+1].ID
		}

		if n1 > n2 {
			dir = "down"
			edgeID = string(n1) + "e" + string(n2)
		} else {
			dir = "up"
			edgeID = string(n2) + "e" + string(n1)
		}
		t.Traffic[edgeID] = append(t.Traffic[edgeID], trafficPacket{color, dir})
	}
}

func (t *trafficMap) removeRoute(r *route, color string) {

}

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
	Feature     feature    `json:"feature`
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

type feature struct {
	Type string `json:"type"` // type of feature
	machine
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
