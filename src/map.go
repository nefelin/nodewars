package main

type NodeMap struct {
	Nodes map[nodeID]*Node
}

type GameState struct {
	Map        NodeMap     `json:"node_map"`
	SeenEvents []GameEvent `json:"events"`
}

type GameEvent struct {
	Who   Player `json:"who"`
	What  string `json:"what"`
	Where Node   `json:"where"`
}

// NodeMap functions

func NewNodeMap() NodeMap {
	return NodeMap{make(map[int]*Node)}
}

func (m *NodeMap) addNode(node *Node) {
	m.Nodes[node.ID] = node
}

func (m *NodeMap) addNodes(ns []*Node) {
	for _, node := range ns {
		m.Nodes[node.ID] = node
	}
}

// Deal with redundant edges on the front end TODO
func (m NodeMap) edges() []Edge {
	edges := make([]Edge, 0)
	for _, node := range m.Nodes {
		for _, connection := range node.Connections {
			edges = append(edges, Edge{node.ID, connection})
		}
	}
	return edges
}

func NewGameState(n NodeMap) GameState {
	return GameState{n, []GameEvent{}}
}

func (s *GameState) addEvent(e GameEvent) {
	s.SeenEvents = append(s.SeenEvents, e)
}

func NewDefaultMap() *NodeMap {
	blueHome := NewNode()
	redHome := NewNode()

	blueHome.addConnection(redHome)
	redHome.addConnection(blueHome)

	newMap := NewNodeMap()
	newMap.addNodes([]*Node{blueHome, redHome})

	return &newMap
}
