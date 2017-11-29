package nwmodel

import (
	"github.com/gorilla/websocket"
)

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

type nodeMap struct {
	Nodes []*node `json:"nodes"`
}

type eventMessage struct {
	Who   Player `json:"who"`
	What  string `json:"what"`
	Where node   `json:"where"`
}

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
}

// module ...
type module struct {
	ID         modID   `json:"id"`
	TestID     int     `json:"testId"`
	LanguageID int     `json:"languageId"`
	Builder    *Player `json:"builder"`
}

type teamName = string

type team struct {
	Name    teamName `json:"name"` // Names are only colors for now
	players map[*Player]bool
	MaxSize int `json:"maxSize"`
}

// Player ...
type Player struct {
	ID   playerID `json:"id"`
	Name string   `json:"name"`
	Team *team    `json:"team"`
	// PointOfEntry   nodeID `json:"pointOfEntry"`
	// NodeConnection nodeID `json:"nodeConnection"`
	// route          []*node
	socket   *websocket.Conn
	outgoing chan Message
}
