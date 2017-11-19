package main

import (
	"fmt"
	"log"
)

type nodeID = int
type edgeID = int

var nodeCount nodeID
var edgeCount edgeID

// Node ...
type Node struct {
	ID          nodeID   `json:"id"`
	Connections []nodeID `json:"connections"`
	Modules     []Module `json:"modules"`
}

// Edge ...
type Edge struct {
	ID     edgeID `json:"id"`
	Source nodeID `json:"source"`
	Target nodeID `json:"target"`
	// Traffic []*Player `json:"traffic"`
}

// Module ...
type Module struct {
	TestID     string  `json:"testId"`
	LanguageID string  `json:"languageId"`
	Owner      *Team   `json:"owner"`
	Builder    *Player `json:"builder"`
}

// func (e *Edge) teamTraffic() []colorName {
// 	colors := make([]colorName, 0)
// 	for _, player := range e.Traffic {
// 		colors = append(colors, player.Team.Name)
// 	}
// 	return colors
// }

// NewNode ...
func NewNode() *Node {
	id := nodeCount
	nodeCount++

	connections := make([]int, 0)
	modules := make([]Module, 0)

	return &Node{id, connections, modules}
}

func newEdge(s, t nodeID) *Edge {
	id := edgeCount
	edgeCount++

	return &Edge{id, s, t}
}

// addConnection is reciprical
func (n *Node) addConnection(m *Node) {
	n.Connections = append(n.Connections, m.ID)
	m.Connections = append(m.Connections, n.ID)
	log.Printf("Node %v added connection to node %v. All connections: %v", n, m, n.Connections)
}

func (n *Node) String() string {
	return fmt.Sprintf("(node)\nID: %v\nConnections:%v\nModules:%v\n", n.ID, n.Connections, n.Modules)
}

// check if a node has a connection to node n
func (n *Node) connectsTo(check nodeID) bool {
	log.Printf("connectTo is checking connection list: %v for target: %v", n.Connections, check)
	log.Printf("Node overview: \n%v", n)
	for _, id := range n.Connections {
		if id == check {
			return true
		}
	}
	return false
}

// returns true if nodes have a valid path
func (n Node) canConnect(m Node) bool {

	return false
}
