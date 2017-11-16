package main

type nodeID = int

var nodeIDCounter int

// Node ...
type Node struct {
	ID          int      `json:"id"`
	Connections []nodeID `json:"connections"`
	Modules     []Module `json:"modules"`
}

// Edge ...
type Edge struct {
	Source nodeID
	Target nodeID
}

// Module ...
type Module struct {
	TestID     string  `json:"testId"`
	LanguageID string  `json:"languageId"`
	Owner      *Team   `json:"owner"`
	Builder    *Player `json:"builder"`
}

// NewNode ...
func NewNode() *Node {
	id := nodeIDCounter
	nodeIDCounter++

	connections := make([]int, 0)
	modules := make([]Module, 0)

	return &Node{id, connections, modules}
}

func (n *Node) addConnection(m *Node) {
	n.Connections = append(n.Connections, m.ID)
}

// check for contested status and set flag (and owner?)
func (n *Node) isContested() {

}

// returns true if nodes have a valid path
func (n Node) canConnect(m Node) bool {

	return false
}
