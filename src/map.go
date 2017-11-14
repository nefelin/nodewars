package main

type NodeMap struct {
	Nodes []*Node `json:"nodes"`
}

type WorldState struct {
	Map        NodeMap
	SeenEvents []GameEvent
}

type GameEvent struct {
	Who   Player
	What  string
	Where Node
}

func (s *WorldState) addEvent(e GameEvent) {
	s.SeenEvents = append(s.SeenEvents, e)
}

func NewWorldState(n NodeMap) WorldState {
	return WorldState{n, []GameEvent{}}
}

func NewDefaultMap() *NodeMap {
	blueHome := NewNode("blue home")
	redHome := NewNode("red home")

	blueHome.addConnection(redHome)
	redHome.addConnection(blueHome)

	return &NodeMap{[]*Node{blueHome, redHome}}
}
