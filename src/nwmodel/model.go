package nwmodel

import (
	"github.com/gorilla/websocket"
)

// type aliases
type nodeID = int
type modID = int
type playerID = int
type teamName = string

// incrementing ID counters
var playerIDCount playerID
var nodeIDCount nodeID
var moduleIDCount modID

// GameModel holds all state information
type GameModel struct {
	Map       *nodeMap             `json:"map"`
	Teams     map[teamName]*team   `json:"teams"`
	Players   map[playerID]*Player `json:"players"`
	POEs      map[playerID]*node   `json:"poes"`
	languages map[string]LanguageDetails
}

type route struct {
	Endpoint *node   `json:"endpoint"`
	Nodes    []*node `json:"nodes"`
}

type nodeMap struct {
	Nodes []*node         `json:"nodes"`
	POEs  map[nodeID]bool `json:"poes"`
}

type node struct {
	ID          nodeID   `json:"id"` // keys and ids is redundant TODO
	Connections []nodeID `json:"connections"`
	// address map concurrency TODO
	Modules map[modID]*module `json:"modules"`
	slots   []*modSlot
}

// TODO rethink I don't like that this setup exposes test information
// throughout the map to the client.
type modSlot struct {
	challenge Challenge // Should we just send this and let front end handle displaying on comannd? seems more in-paradigm to send only on request.
	module    *module
}

type module struct {
	id        modID  // `json:"id"`
	language  string // `json:"languageId"`
	builder   string // `json:"creator"`
	Health    int    `json:"health"`
	MaxHealth int    `json:"maxHealth"`
	Team      *team  `json:"team"`
}

type team struct {
	Name    string `json:"name"` // Names are only colors for now
	players map[*Player]bool
	maxSize int `json:"maxSize"`
	poe     *node
}

// TODO un export all but route

// Player ...
type Player struct {
	ID       playerID `json:"id"`
	Name     string   `json:"name"`
	Team     *team    `json:"team"`
	Route    *route   `json:"route"`
	socket   *websocket.Conn
	outgoing chan Message
	language string // current working language
	stdin    string // stdin buffer for testing
	slotNum  int    // currently attached to slotNum of current node
}
