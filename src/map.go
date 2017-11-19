package main

import (
	"errors"
	"log"
)

// using lists is more friendly to graphing library. Change? TODO
type nodeMap struct {
	Nodes map[nodeID]*Node `json:"nodes"`
	Edges map[edgeID]*Edge `json:"edges"`
}

type gameState struct {
	Map        nodeMap     `json:"nodeMap"`
	SeenEvents []gameEvent `json:"events"`
}

type gameEvent struct {
	Who   Player `json:"who"`
	What  string `json:"what"`
	Where Node   `json:"where"`
}

// nodeMap functions

func newNodeMap() nodeMap {
	return nodeMap{make(map[nodeID]*Node), make(map[edgeID]*Edge)}
}

func (m *nodeMap) addNodes(ns ...*Node) {
	for _, node := range ns {
		m.Nodes[node.ID] = node
	}
}

func (m *nodeMap) addEdges(es ...*Edge) {
	// check for redundancy TODO
	for _, edge := range es {
		m.Edges[edge.ID] = edge
	}
}

func (m *nodeMap) connectNodes(n1, n2 nodeID) error {
	log.Println("connectNodes called")
	// Check existence of both elements
	_, ok1 := m.Nodes[1]
	if _, ok2 := m.Nodes[n2]; !ok1 || !ok2 {
		log.Println("connectNodes error")
		return errors.New("One or both nodes out of range")
	}

	// add connection value to each node,
	m.Nodes[n1].addConnection(m.Nodes[n2])

	// add edge to self
	m.addEdges(newEdge(n1, n2))
	return nil
}

func newGameState(n nodeMap) gameState {
	return gameState{n, []gameEvent{}}
}

func (s *gameState) addEvent(e gameEvent) {
	s.SeenEvents = append(s.SeenEvents, e)
}

func newDefaultMap() *nodeMap {
	newMap := newNodeMap()
	NODECOUNT := 12

	for i := 0; i < NODECOUNT; i++ {
		//Make new nodes
		newMap.addNodes(NewNode())
	}

	for i := 0; i < NODECOUNT; i++ {
		//Make new edges
		targ1, targ2 := -1, -1

		if i < NODECOUNT-3 {
			targ1 = i + 3
		}

		if i%3 < 2 {
			targ2 = i + 1
		}

		if targ1 != -1 {
			newMap.connectNodes(i, targ1)
		}

		if targ2 != -1 {
			newMap.connectNodes(i, targ2)
		}
	}

	return &newMap
}
