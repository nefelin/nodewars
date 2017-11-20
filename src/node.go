package main

import (
	"fmt"
	"log"
)

type nodeID = int
type edgeID = int

var nodeCount nodeID
var edgeCount edgeID

// node ...
type node struct {
	ID          nodeID   `json:"id"`
	Connections []nodeID `json:"connections"`
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
	Owner      *Team   `json:"owner"`
	Builder    *Player `json:"builder"`
}

// func (e *edge) teamTraffic() []colorName {
// 	colors := make([]colorName, 0)
// 	for _, player := range e.Traffic {
// 		colors = append(colors, player.Team.Name)
// 	}
// 	return colors
// }

// NewNode ...
func NewNode() *node {
	id := nodeCount
	nodeCount++

	connections := make([]int, 0)
	modules := make([]module, 0)

	return &node{id, connections, modules}
}

func newEdge(s, t nodeID) *edge {
	id := edgeCount
	edgeCount++

	return &edge{id, s, t}
}

// addConnection is reciprical
func (n *node) addConnection(m *node) {
	n.Connections = append(n.Connections, m.ID)
	m.Connections = append(m.Connections, n.ID)
	// log.Printf("node %v added connection to node %v. All connections: %v", n, m, n.Connections)
}

func (n *node) String() string {
	return fmt.Sprintf("<(node) ID: %v, Connections:%v, Modules:%v>", n.ID, n.Connections, n.Modules)
}

// check if a node has a connection to node n
func (n *node) connectsTo(check nodeID) bool {
	log.Printf("connectTo is checking connection list: %v for target: %v", n.Connections, check)
	log.Printf("node overview: \n%v", n)
	for _, id := range n.Connections {
		if id == check {
			return true
		}
	}
	return false
}

// returns true if nodes have a valid path
func (n node) canConnect(m node) bool {

	return false
}
