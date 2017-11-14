package main

type NodeMap = map[int]*Node

type GameState struct {
	Map        NodeMap     `json:"node_map"`
	SeenEvents []GameEvent `json:"events"`
}

type GameEvent struct {
	Who   Player `json:"who"`
	What  string `json:"what"`
	Where Node   `json:"where"`
}

func NewNodeMap() NodeMap {
	return make(map[int]*Node)
}

func (s *GameState) addEvent(e GameEvent) {
	s.SeenEvents = append(s.SeenEvents, e)
}

func NewGameState(n NodeMap) GameState {
	return GameState{n, []GameEvent{}}
}

func addNodeTo(m NodeMap, node *Node) {
	m[node.ID] = node
}

func addNodesTo(m NodeMap, ns []*Node) {
	for _, node := range ns {
		m[node.ID] = node
	}
}

func NewDefaultMap() *NodeMap {
	blueHome := NewNode()
	redHome := NewNode()

	blueHome.addConnection(redHome)
	redHome.addConnection(blueHome)

	newMap := NewNodeMap()
	addNodesTo(newMap, []*Node{blueHome, redHome})

	return &newMap
}
