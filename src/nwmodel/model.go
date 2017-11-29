package nwmodel

import (
	"github.com/gorilla/websocket"
)

type playerName string

// GameModel holds all state information
type GameModel struct {
	Map     *nodeMap             `json:"map"`
	Teams   map[teamName]*team   `json:"teams"`
	Players map[playerID]*Player `json:"players"`
	Routes  map[playerID]*route  `json:"routes"`
	POEs    map[playerID]*node   `json:"poes"`
	// CurrentEvents []*gameEvent      `json:"currentEvents"`

}

type route struct {
	Endpoint *node   `json:"endpoint"`
	Nodes    []*node `json:"nodes"`
}

// simplifies gameModel for export, exclusively used for sending updates to the player
// type gameState struct {
// 	Map           nodeMap     `json:"nodeMap"`
// 	Teams         []string    `json:"teams"`
// 	Players       []*Player   `json:"players"`
// 	CurrentEvents []gameEvent `json:"currentEvents"`
// }

// using lists is more friendly to graphing library. Change? TODO
// type nodeMap struct {
// 	Nodes map[nodeID]*node `json:"nodes"`
// 	Edges map[edgeID]*edge `json:"edges"`
// }

type nodeMap struct {
	Nodes []*node `json:"nodes"`
}

type eventMessage struct {
	Who   Player `json:"who"`
	What  string `json:"what"`
	Where node   `json:"where"`
}

// type gameEvent struct {
// 	Who   Player `json:"who"`
// 	What  string `json:"what"`
// 	Where node   `json:"where"`
// }

type nodeID = int
type modID = int
type playerID = int

var playerIDCount playerID
var nodeIDCount nodeID
var moduleIDCount modID

// var edgeCount edgeID

// node ...
type node struct {
	ID          nodeID           `json:"id"` // keys and ids is redundant TODO
	Connections []nodeID         `json:"connections"`
	Size        int              `json:"size"`
	Modules     map[modID]module `json:"modules"`
	// Traffic          []*Player        `json:"traffic"`
	// POE              []*Player        `json:"poe"`
	// ConnectedPlayers []*Player        `json:"connectedPlayers"`
}

// edge ...
// type edge struct {
// 	ID      edgeID    `json:"id"` // keys and ids is redundant TODO
// 	Source  nodeID    `json:"source"`
// 	Target  nodeID    `json:"target"`
// 	Traffic []*Player `json:"traffic"`
// }

// module ...
type module struct {
	ID         modID   `json:"id"`
	TestID     int     `json:"testId"`
	LanguageID int     `json:"languageId"`
	Builder    *Player `json:"builder"`
}

type teamName string

type team struct {
	Name    teamName `json:"name"` // Names are only colors for now
	players map[*Player]bool
	MaxSize int `json:"maxSize"`
}

// Player ...
type Player struct {
	ID   playerID   `json:"id"`
	Name playerName `json:"name"`
	Team *team      `json:"team"`
	// PointOfEntry   nodeID `json:"pointOfEntry"`
	// NodeConnection nodeID `json:"nodeConnection"`
	// route          []*node
	socket   *websocket.Conn
	outgoing chan Message
}
