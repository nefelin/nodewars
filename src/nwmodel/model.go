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
	Map     *nodeMap             `json:"map"`
	Teams   map[teamName]*team   `json:"teams"`
	Players map[playerID]*Player `json:"players"`
	POEs    map[playerID]*node   `json:"poes"`
}

type route struct {
	Endpoint *node   `json:"endpoint"`
	Nodes    []*node `json:"nodes"`
}

type nodeMap struct {
	Nodes []*node `json:"nodes"`
}

type node struct {
	ID          nodeID            `json:"id"` // keys and ids is redundant TODO
	Connections []nodeID          `json:"connections"`
	Modules     map[modID]*module `json:"modules"`
	// using map enables more interesting slot names
	Slot map[string]*modSlot `json:"slot"`
}

type modSlot struct {
	ChallengeID int     `json:"challengeID"`
	Module      *module `json:"module"`
}

type module struct {
	ID         modID `json:"id"`
	TestID     int   `json:"testId"`
	LanguageID int   `json:"languageId"`
	// Builder    *Player `json:"builder"`
	Team *team `json:"team"`
}

type team struct {
	Name    string `json:"name"` // Names are only colors for now
	players map[*Player]bool
	MaxSize int `json:"maxSize"`
}

// Player ...
type Player struct {
	ID       playerID `json:"id"`
	Name     string   `json:"name"`
	Team     *team    `json:"team"`
	Route    *route   `json:"route"`
	socket   *websocket.Conn
	outgoing chan Message
}
