package main

// import "strconv"

var nodeID int

// Node ...
type Node struct {
	ID          int   `json:"id"`
	Connections []int `json:"connections"`
	// Owner       *Team   `json:"owner"`
	// Contested   bool    `json:"contested"`
	// Ice []Ice `json:"ice"`
}

// Ice ...
type Ice struct {
	TestID string `json:"test_id"`
	Owner  Team   `json:"owner"`
}

// NewNode ...
func NewNode() *Node {
	id := nodeID
	nodeID++

	connections := make([]int, 0)
	// ice := make([]Ice, 0)

	return &Node{id, connections}
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
