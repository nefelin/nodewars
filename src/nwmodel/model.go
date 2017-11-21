package nwmodel

import (
	"github.com/gorilla/websocket"
)

// GameModel holds all state information
type GameModel struct {
	Map           *nodeMap         `json:"map"`
	Teams         map[string]*team `json:"teams"`
	Players       map[*Player]bool `json:"players"`
	CurrentEvents []*gameEvent     `json:"currentEvents"`
}

// simplifies gameModel for export, exclusively used for sending updates to the player
type gameState struct {
	Map           nodeMap     `json:"nodeMap"`
	Teams         []string    `json:"teams"`
	Players       []*Player   `json:"players"`
	CurrentEvents []gameEvent `json:"currentEvents"`
}

// using lists is more friendly to graphing library. Change? TODO
type nodeMap struct {
	Nodes map[nodeID]*node `json:"nodes"`
	Edges map[edgeID]*edge `json:"edges"`
}

type gameEvent struct {
	Who   Player `json:"who"`
	What  string `json:"what"`
	Where node   `json:"where"`
}

type nodeID = int
type edgeID = int

var nodeCount nodeID
var edgeCount edgeID

// node ...
type node struct {
	ID          nodeID   `json:"id"`
	Connections []nodeID `json:"connections"`
	Size        int      `json:"size"`
	Modules     []module `json:"modules"`
}

// edge ...
type edge struct {
	ID     edgeID `json:"id"`
	Source nodeID `json:"source"`
	Target nodeID `json:"target"`
	// Traffic []*Player `json:"traffic"`
}

// module ...
type module struct {
	TestID     string  `json:"testId"`
	LanguageID string  `json:"languageId"`
	Owner      *team   `json:"owner"`
	Builder    *Player `json:"builder"`
}

type teamName = string

// team ...
type team struct {
	Name    teamName         // Names are only colors for now
	Players map[*Player]bool //THESE create circular JSON problem, decide on org scheme and fix, TODO
	MaxSize int
}

// Player ...
type Player struct {
	Name         string
	team         *team
	PointOfEntry nodeID
	socket       *websocket.Conn
	outgoing     chan Message
}
