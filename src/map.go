package main

import (
	"errors"
	"log"
)

// using lists is more friendly to graphing library. Change? TODO
type nodeMap struct {
	Nodes map[nodeID]*node `json:"nodes"`
	Edges map[edgeID]*edge `json:"edges"`
}

type gameState struct {
	Map        nodeMap     `json:"nodeMap"`
	SeenEvents []gameEvent `json:"events"`
}

type gameEvent struct {
	Who   Player `json:"who"`
	What  string `json:"what"`
	Where node   `json:"where"`
}

// nodeMap functions

func newNodeMap() nodeMap {
	return nodeMap{make(map[nodeID]*node), make(map[edgeID]*edge)}
}

func (m *nodeMap) addNodes(ns ...*node) {
	for _, node := range ns {
		m.Nodes[node.ID] = node
	}
}

func (m *nodeMap) addEdges(es ...*edge) {
	// check for redundancy TODO
	for _, edge := range es {
		m.Edges[edge.ID] = edge
	}
}

func (m *nodeMap) connectNodes(n1, n2 nodeID) error {
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

// nodesConnections takes one of the maps nodes and converts it connections (in from of nodeIDs) into pointers to actual node objects
func (m *nodeMap) nodesConnections(n *node) []*node {
	res := make([]*node, 0)
	for _, nodeID := range n.Connections {
		res = append(res, m.Nodes[nodeID])
	}

	return res
}

// constructPath takes the routes discovered via routeToNode and the endpoint (target) and creates a slice of the correct path, note order is still reversed and path contains source but not target node
func constructPath(prevMap map[*node]*node, t *node) []*node {
	route := make([]*node, 0)

	for step, ok := prevMap[t]; ok; step, ok = prevMap[step] {
		route = append(route, step)
	}

	return route
}

// getBestNode extract the node with shortes path from pool, it is a poor substitute for using a priority queue
func getBestNode(pool map[*node]bool, distMap map[*node]int) *node {
	bestDist := 100000
	var bestNode *node
	for node := range pool {
		if distMap[node] < bestDist {
			bestNode = node
			bestDist = distMap[node]
		}
	}
	return bestNode
}

// routeToNode uses vanilla dijkstra's (vanilla for now) algorithm to find node path
func (m *nodeMap) routeToNode(source, target *node) []*node {
	unchecked := make(map[*node]bool) // this should be a priority queue for efficiency
	dist := make(map[*node]int)
	prev := make(map[*node]*node)

	for _, node := range m.Nodes {
		dist[node] = 10000
		unchecked[node] = true
	}

	dist[source] = 0

	for len(unchecked) > 0 {
		thisNode := getBestNode(unchecked, dist)

		delete(unchecked, thisNode)

		if thisNode == target {
			log.Println("Found target!")
			log.Printf("%v", constructPath(prev, target))
			return make([]*node, 0)
		}

		for _, cNode := range m.nodesConnections(thisNode) {
			alt := dist[thisNode] + 1
			if alt < dist[cNode] {
				dist[cNode] = alt
				prev[cNode] = thisNode
			}
		}
	}
	log.Println("No possible route")
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
