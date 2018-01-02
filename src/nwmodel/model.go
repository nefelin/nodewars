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
	Map       *nodeMap             `json:"map"`
	Teams     map[teamName]*team   `json:"teams"`
	Players   map[playerID]*Player `json:"players"`
	POEs      map[playerID]*node   `json:"poes"`
	PointGoal float32              `json:"pointGoal"`
	languages map[string]LanguageDetails
	aChan     chan nwmessage.Message
	running   bool //running should replace mapLocked

	// timelimit should be able to set a timelimit and count points at the end
}

type route struct {
	Endpoint *node   `json:"endpoint"`
	Nodes    []*node `json:"nodes"`
}

type nodeMap struct {
	Nodes       []*node         `json:"nodes"`
	POEs        map[nodeID]bool `json:"poes"`
	diameter    float64
	radius      float64
	nodeIDCount nodeID
}

type node struct {
	ID          nodeID   `json:"id"` // keys and ids is redundant? TODO
	Connections []nodeID `json:"connections"`
	// Modules     map[modID]*module `json:"modules"`
	Slots       []*modSlot `json:"slots"`
	Remoteness  float64    `json:"remoteness"`
	playersHere []playerID
}

type modSlot struct {
	sync.Mutex
	challenge Challenge
	Type      string  `json:"type"`
	Module    *module `json:"module"`
	Powered   bool    `json:"powered"`
}

type module struct {
	id        modID  // `json:"id"`
	language  string // `json:"languageId"`
	builder   string // `json:"creator"`
	Health    int    `json:"health"`
	MaxHealth int    `json:"maxHealth"`
	TeamName  string `json:"team"`
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
}
