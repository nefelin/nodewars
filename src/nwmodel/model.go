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
type modID = int
type edgeID = string

var nodeCount nodeID
var moduleCount modID

// var edgeCount edgeID

// node ...
// TODO make a decision about slices vs maps
type node struct {
	ID               nodeID           `json:"id"` // keys and ids is redundant TODO
	Connections      []nodeID         `json:"connections"`
	Size             int              `json:"size"`
	Modules          map[modID]module `json:"modules"`
	Traffic          []*Player        `json:"traffic"`
	POE              []*Player        `json:"poe"`
	ConnectedPlayers []*Player        `json:"connectedPlayers"`
}

// edge ...
type edge struct {
	ID      edgeID    `json:"id"` // keys and ids is redundant TODO
	Source  nodeID    `json:"source"`
	Target  nodeID    `json:"target"`
	Traffic []*Player `json:"traffic"`
}

// module ...
type module struct {
	ID         modID `json:"id"`
	TestID     int   `json:"testId"`
	LanguageID int   `json:"languageId"`
	// Owner      *team   `json:"owner"`
	Builder *Player `json:"builder"`
}

type teamName = string

// team ...
type team struct {
	Name    teamName `json:"name"` // Names are only colors for now
	players map[*Player]bool
	MaxSize int `json:"maxSize"`
}

// Player ...
type Player struct {
	Name           string `json:"name"`
	Team           *team  `json:"team"`
	PointOfEntry   nodeID `json:"pointOfEntry"`
	NodeConnection nodeID `json:"nodeConnection"`
	route          []*node
	socket         *websocket.Conn
	outgoing       chan Message
}
