package nwmodel

import (
	"errors"
	"log"
	"sort"
)

type nodeMap struct {
	Nodes []*node `json:"nodes"`
	// POEs        map[nodeID]bool `json:"poes"`
	diameter    float64
	radius      float64
	nodeIDCount nodeID
}

// initializer:
func newNodeMap() nodeMap {
	return nodeMap{
		Nodes: make([]*node, 0),
		// POEs:     make(map[nodeID]bool),
		diameter: 0,
		radius:   1000,
	}
}

// nodeMap methods -----------------------------------------------------------------------------

func (m *nodeMap) initAllNodes() {
	m.initAllSlots()
	m.initAllRemoteness()
}

func (m *nodeMap) initAllSlots() {
	// initialize each node's slots
	for _, node := range m.Nodes {
		node.initMachines()
	}
}

func (m *nodeMap) initAllRemoteness() {
	// initialize each nodes remoteness.
	for _, node := range m.Nodes {
		node.Remoteness = float64(m.findNodeEccentricity(node))
		if node.Remoteness > m.diameter {
			m.diameter = node.Remoteness
		}
		if node.Remoteness < m.radius {
			m.radius = node.Remoteness
		}
		// log.Printf("Node %d, eccentricity: %d", node.ID, node.Remoteness)
	}

	for _, node := range m.Nodes {
		node.Remoteness = node.Remoteness / m.diameter
		// log.Printf("Node %d, remoteness: %d", node.ID, node.Remoteness)
	}
}

func (m *nodeMap) findNodeEccentricity(n *node) int {
	// for every node, count distance to other nodes, pick the largest
	var maxDist int
	// var farthesNode nodeID
	for _, node := range m.Nodes {

		// don't check our starting point
		if n != node {
			nodePath := m.routeToNode(nil, n, node)
			if len(nodePath) > maxDist {
				maxDist = len(nodePath)
				// farthestNode = node.ID
			}
		}

	}
	// log.Printf("Farthest node from %d: %d", n.ID, farthesNode)
	return maxDist
}

func (m *nodeMap) addPoes(ns ...nodeID) {
	for _, id := range ns {
		// skip bad ids
		if !m.nodeExists(id) {
			continue
		}
		// make an available POE for each nodeID passed
		m.Nodes[id].Feature.Type = "poe"
		// m.POEs[id] = true
	}
}

func (m *nodeMap) collectEmptyPoes() []*node {
	poes := make([]*node, 0)
	for _, node := range m.Nodes {
		if node.Feature.Type == "poe" {
			poes = append(poes, node)
		}
	}
	return poes
}

// initPoes right now places poes at remotest locations, which is not idea if remoteness = value
func (m *nodeMap) initPoes(n int) {
	// make a map of remotesnesses to nodes
	remMap := make(map[float64][]*node)
	for _, node := range m.Nodes {
		remMap[node.Remoteness] = append(remMap[node.Remoteness], node)
	}

	// create sorted list of remotenesses
	ordRem := make([]float64, len(remMap))
	I := 0
	for k := range remMap {
		ordRem[I] = k
		I++
	}

	sort.Sort(sort.Reverse(sort.Float64Slice(ordRem)))

	// check each remoteness in ascending order,

	// if not, move up to the next remoteness tier
	for _, v := range ordRem {
		// if there're enough nodes of that remoteness, seen if we can place n poes in those remotenesses
		if len(remMap[v]) >= n {
			// add a poe at each node of that remoteness
			for _, node := range remMap[v] {
				m.addPoes(node.ID)
			}
			break
		}
		// log.Printf("We have %d nodes of remoteness %v", len(remMap[ordRem[i]]), ordRem[i])
	}

	// if all else fails, assign at random
	// only ensuring distance

}

// func (m *nodeMap) NewNode() *node {
// 	id := m.nodeIDCount
// 	m.nodeIDCount++

// 	connections := make([]int, 0)
// 	modules := make(map[modID]*module)

// 	return &node{
// 		ID:          id,
// 		Connections: connections,
// 		Modules:     modules,
// 	}
// }

func (m *nodeMap) newNode() *node {
	id := m.nodeIDCount
	m.nodeIDCount++

	return &node{
		ID:          id,
		Connections: make([]int, 0),
		Remoteness:  100,
		Machines:    []*machine{newMachine()},
		Feature:     newFeature(),
	}
}

func (m *nodeMap) addNodes(count int) []*node {
	enter := make([]*node, count)
	for i := 0; i < count; i++ {

		newNode := m.newNode()

		enter[i] = newNode
		m.Nodes = append(m.Nodes, newNode)
	}
	return enter
}

func (m *nodeMap) removeNodes(ns []int) {
	for _, id := range ns {
		// look at connections and remove any connections point to node
		for _, conn := range m.Nodes[id].Connections {
			m.Nodes[conn].remConnection(id)
		}

		// remove from the Map.Nodes list
		m.Nodes[id] = nil
	}

	// fix holes in the slice
	m.Nodes = fillNodeSliceHoles(m.Nodes)

	// fix node ID count
	m.nodeIDCount -= len(ns)
}

func fillNodeSliceHoles(ns []*node) []*node {
	// for every node i the node slice
	for i := 0; i < len(ns); i++ {
		// for i, n := range ns {
		n := ns[i]
		// when we find a nil slot
		if n == nil {
			j := len(ns) - 1
			// if the last element is also nil
			for ns[j] == nil {
				// keep backing up till we find a non nil
				j--
				// if we wind up where we started skip and just cut the nils out of the list
				if j == i {
					break
				}
			}
			// once we've found a non nil to swap, swap
			ns[i], ns[j] = ns[j], ns[i]

			if ns[i] != nil {
				// fix connections and ids
				for _, connI := range ns[i].Connections {
					// for each of our old connections, remove pointer to our old location
					ns[connI].remConnection(ns[i].ID)

					// add a fresh connection to our updated location
					ns[connI].Connections = append(ns[connI].Connections, i)
				}
				// fix our id to match our new index
				ns[i].ID = i
			}
			// cut the (should be all nils) tail off the slice
			ns = ns[:j]
		}
	}
	return ns
}

func (m *nodeMap) connectNodes(n1, n2 nodeID) error {
	// Check existence of both elements
	if m.nodeExists(n1) && m.nodeExists(n2) {

		// add connection value to each node,
		m.Nodes[n1].addConnection(m.Nodes[n2])
		return nil

	}

	log.Println("connectNodes error")
	return errors.New("One or both nodes out of range")
}

func (m *nodeMap) nodeExists(n nodeID) bool {
	if n > -1 && n < len(m.Nodes) {
		return true
	}
	return false
}

// nodesConnections takes one of the maps nodes and converts its connections (in the form of nodeIDs) into pointers to actual node objects
// TODO ask about this, feels hacky
func (m *nodeMap) nodesConnections(n *node) []*node {
	res := make([]*node, 0)
	for _, nodeID := range n.Connections {
		res = append(res, m.Nodes[nodeID])
	}
	return res
}

func (m *nodeMap) nodesTouch(n1, n2 *node) bool {
	// for every one of n1's connections
	for _, connectedNode := range m.nodesConnections(n1) {
		// if it is n2, return true
		if connectedNode == n2 {
			return true
		}
	}
	return false
}

type searchField struct {
	unchecked map[*node]bool
	dist      map[*node]int
	prev      map[*node]*node
}

func (m *nodeMap) newSearchField(t *team, source *node) searchField {
	retField := searchField{
		unchecked: make(map[*node]bool), // TODO this should be a priority queue for efficiency
		dist:      make(map[*node]int),
		prev:      make(map[*node]*node),
	}

	seen := make(map[*node]bool)
	tocheck := make([]*node, 1)
	tocheck[0] = source

	for len(tocheck) > 0 {
		thisNode := tocheck[0]
		tocheck = tocheck[1:]
		// log.Printf("this: %v", thisNode)
		// t == nil signifies that we don't care about routability and we want a field containing the whole (contiguous) map
		if t == nil || thisNode.allowsRoutingFor(t) {
			retField.unchecked[thisNode] = true
			retField.dist[thisNode] = 1000
			seen[thisNode] = true
			for _, nodeID := range thisNode.Connections {
				// log.Printf("nodeid: %v", nodeID)
				if !seen[m.Nodes[nodeID]] {
					tocheck = append(tocheck, m.Nodes[nodeID])

				}
				// log.Printf("tocheck %v", tocheck)
			}
		}
	}

	return retField
}

// routeToNode uses vanilla dijkstra's (vanilla for now) algorithm to find node path
// TODO get code review on this. I think I'm maybe not getting optimal route
func (m *nodeMap) routeToNode(t *team, source, target *node) []*node {

	if source.allowsRoutingFor(t) {
		// if we're connecting to our POE, return a route which is only our POE
		if source == target {
			route := make([]*node, 1)
			route[0] = source
			return route
		}

		nodePool := m.newSearchField(t, source)

		nodePool.dist[source] = 0

		for len(nodePool.unchecked) > 0 {
			thisNode := getBestNode(nodePool.unchecked, nodePool.dist)

			delete(nodePool.unchecked, thisNode)

			if m.nodesTouch(thisNode, target) {
				nodePool.prev[target] = thisNode
				route := constructPath(nodePool.prev, target)
				// log.Println("Found target!")
				return route
			}

			for _, cNode := range m.nodesConnections(thisNode) {
				// TODO refactor to take least risky routes by weighing against vulnerability to enemy connection
				alt := nodePool.dist[thisNode] + 1
				if alt < nodePool.dist[cNode] {
					nodePool.dist[cNode] = alt
					nodePool.prev[cNode] = thisNode
				}
			}
		}
	}
	//else {
	// log.Println("POE Blocked")
	//}
	// log.Println("No possible route")
	return nil
}

// helper functions for routeToNode ------------------------------------------------------------
// constructPath takes the routes discovered via routeToNode and the endpoint (target) and creates a slice of the correct path, note order is still reversed and path contains source but not target node
func constructPath(prevMap map[*node]*node, t *node) []*node {
	// log.Printf("constructPath working from prev: %v", prevMap)

	route := make([]*node, 0)

	for step, ok := prevMap[t]; ok; step, ok = prevMap[step] {
		route = append(route, step)
	}

	return route
}

// getBestNode TODO extract the node with shortes path from pool, it is a substitute for using a priority queue
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
