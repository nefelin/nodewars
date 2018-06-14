package node

import (
	"errors"
	"feature"
	"log"
	"model/machines"
	"sort"
)

type Map struct {
	Nodes       []*Node `json:"nodes"`
	diameter    float64
	radius      float64
	NodeIDCount NodeID
}

// initializer:
func NewMap() *Map {
	return &Map{
		Nodes:    make([]*Node, 0),
		diameter: 0,
		radius:   1000,
	}
}

// Map methods -----------------------------------------------------------------------------

func (m *Map) initAllNodes() {
	m.initAllMachines()
	m.initAllRemoteness()
}

func (m *Map) initAllMachines() {
	// initialize each Node's machines
	for _, Node := range m.Nodes {
		Node.initMachines()
	}
}

func (m *Map) initAllRemoteness() {
	// initialize each Nodes remoteness.
	for _, Node := range m.Nodes {
		Node.Remoteness = float64(m.findNodeEccentricity(Node))
		if Node.Remoteness > m.diameter {
			m.diameter = Node.Remoteness
		}
		if Node.Remoteness < m.radius {
			m.radius = Node.Remoteness
		}
		// log.Printf("Node %d, eccentricity: %d", Node.ID, Node.Remoteness)
	}

	for _, Node := range m.Nodes {
		Node.Remoteness = Node.Remoteness / m.diameter
		// log.Printf("Node %d, remoteness: %d", Node.ID, Node.Remoteness)
	}
}

func (m *Map) findNodeEccentricity(n *Node) int {
	// for every Node, count distance to other Nodes, pick the largest
	var maxDist int
	// var farthesNode NodeID
	for _, Node := range m.Nodes {

		// don't check our starting point
		if n != Node {
			NodePath := m.RouteToNode("", n, Node)
			if NodePath.Length() > maxDist {
				maxDist = NodePath.Length()
				// farthestNode = Node.ID
			}
		}

	}
	// log.Printf("Farthest Node from %d: %d", n.ID, farthesNode)
	return maxDist
}

func (m *Map) addPoes(ns ...NodeID) {
	for _, id := range ns {
		// skip bad ids
		Node := m.GetNode(id)
		if Node == nil {
			continue
		}
		// make an available POE for each NodeID passed
		Node.Feature.Type = feature.POE
		// m.POEs[id] = true
	}
}

func (m *Map) CollectEmptyPoes() []*Node {
	poes := make([]*Node, 0)
	for _, Node := range m.Nodes {
		if Node.Feature.Type == feature.POE {
			poes = append(poes, Node)
		}
	}
	return poes
}

// initPoes right now places poes at remotest locations, which is not idea if remoteness = value
func (m *Map) initPoes(n int) {
	// make a map of remotesnesses to Nodes
	remMap := make(map[float64][]*Node)
	for _, Node := range m.Nodes {
		remMap[Node.Remoteness] = append(remMap[Node.Remoteness], Node)
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
		// if there're enough Nodes of that remoteness, seen if we can place n poes in those remotenesses
		if len(remMap[v]) >= n {
			// add a poe at each Node of that remoteness
			for _, Node := range remMap[v] {
				m.addPoes(Node.ID)
			}
			break
		}
		// log.Printf("We have %d Nodes of remoteness %v", len(remMap[ordRem[i]]), ordRem[i])
	}

	// if all else fails, assign at random
	// only ensuring distance

}

// func (m *Map) NewNode() *Node {
// 	id := m.NodeIDCount
// 	m.NodeIDCount++

// 	connections := make([]int, 0)
// 	modules := make(map[modID]*module)

// 	return &Node{
// 		ID:          id,
// 		Connections: connections,
// 		Modules:     modules,
// 	}
// }

func (m *Map) newNode() *Node {
	id := m.NodeIDCount
	m.NodeIDCount++

	return &Node{
		ID:          id,
		Connections: make([]int, 0),
		Remoteness:  100,
		Machines:    []*machines.Machine{machines.NewMachine()},
		Feature:     machines.NewFeature(),
		addressMap:  make(map[string]*machines.Machine),
	}
}

func (m *Map) addNodes(count int) []*Node {
	enter := make([]*Node, count)
	for i := 0; i < count; i++ {

		newNode := m.newNode()

		enter[i] = newNode
		m.Nodes = append(m.Nodes, newNode)
	}
	return enter
}

func (m *Map) removeNodes(ns []int) {
	for _, id := range ns {
		// look at connections and remove any connections point to Node
		for _, conn := range m.Nodes[id].Connections {
			m.Nodes[conn].remConnection(id)
		}

		// remove from the Map.Nodes list
		m.Nodes[id] = nil
	}

	// fix holes in the slice
	m.Nodes = fillNodeSliceHoles(m.Nodes)

	// fix Node ID count
	m.NodeIDCount -= len(ns)
}

func fillNodeSliceHoles(ns []*Node) []*Node {
	// for every Node i the Node slice
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

func (m *Map) connectNodes(n1, n2 NodeID) error {
	// Check existence of both elements
	Node1 := m.GetNode(n1)
	Node2 := m.GetNode(n2)

	if Node1 == nil && Node2 == nil {
		log.Println("connectNodes error")
		return errors.New("One or both Nodes out of range")
	}

	m.Nodes[n1].addConnection(m.Nodes[n2])
	return nil
}

func (m *Map) GetNode(n NodeID) *Node {
	if n < 0 || n > len(m.Nodes)-1 {
		return nil
	}
	return m.Nodes[n]
}

// NodesConnections takes one of the maps Nodes and converts its connections (in the form of NodeIDs) into pointers to actual Node objects
// TODO ask about this, feels hacky
func (m *Map) NodesConnections(n *Node) []*Node {
	res := make([]*Node, 0)
	for _, NodeID := range n.Connections {
		res = append(res, m.Nodes[NodeID])
	}
	return res
}

func (m *Map) NodesTouch(n1, n2 *Node) bool {
	// for every one of n1's connections
	for _, connectedNode := range m.NodesConnections(n1) {
		// if it is n2, return true
		if connectedNode == n2 {
			return true
		}
	}
	return false
}

type searchField struct {
	unchecked map[*Node]bool
	dist      map[*Node]int
	prev      map[*Node]*Node
}

func (m *Map) newSearchField(t teamName, source *Node) searchField {
	retField := searchField{
		unchecked: make(map[*Node]bool), // TODO this should be a priority queue for efficiency
		dist:      make(map[*Node]int),
		prev:      make(map[*Node]*Node),
	}

	seen := make(map[*Node]bool)
	tocheck := make([]*Node, 1)
	tocheck[0] = source

	for len(tocheck) > 0 {
		thisNode := tocheck[0]
		tocheck = tocheck[1:]
		// log.Printf("this: %v", thisNode)
		// t == nil signifies that we don't care about routability and we want a field containing the whole (contiguous) map
		if t == "" || thisNode.supportsRouting(t) {
			retField.unchecked[thisNode] = true
			retField.dist[thisNode] = 1000
			seen[thisNode] = true
			for _, NodeID := range thisNode.Connections {
				// log.Printf("Nodeid: %v", NodeID)
				if !seen[m.Nodes[NodeID]] {
					tocheck = append(tocheck, m.Nodes[NodeID])

				}
				// log.Printf("tocheck %v", tocheck)
			}
		}
	}

	// fmt.Printf("searchField: %+v\n", retField)
	return retField
}

// routeToNode uses vanilla dijkstra's (vanilla for now) algorithm to find Node path
// TODO get code review on this. I think I'm maybe not getting optimal route
func (m *Map) RouteToNode(t teamName, source, target *Node) *Route {

	if source.HasMachineFor(t) {
		// if we're connecting to our POE, return a route which is only our POE
		if source == target {
			route := make([]*Node, 1)
			route[0] = source
			return &Route{route}
		}

		NodePool := m.newSearchField(t, source)

		NodePool.dist[source] = 0

		for len(NodePool.unchecked) > 0 {
			thisNode := getBestNode(NodePool.unchecked, NodePool.dist)

			delete(NodePool.unchecked, thisNode)

			if m.NodesTouch(thisNode, target) {
				NodePool.prev[target] = thisNode
				route := constructRoute(NodePool.prev, target)
				// log.Println("Found target!")
				return &route
			}

			for _, cNode := range m.NodesConnections(thisNode) {
				// TODO refactor to take least risky routes by weighing against vulnerability to enemy connection
				alt := NodePool.dist[thisNode] + 1
				if alt < NodePool.dist[cNode] {
					NodePool.dist[cNode] = alt
					NodePool.prev[cNode] = thisNode
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
// constructRoute takes the routes discovered via routeToNode and the endpoint (target) and creates a slice of the correct path, note order is still reversed and path contains source but not target Node
func constructRoute(prevMap map[*Node]*Node, t *Node) Route {
	// log.Printf("constructRoute working from prev: %v", prevMap)

	route := make([]*Node, 1)
	route[0] = t

	for step, ok := prevMap[t]; ok; step, ok = prevMap[step] {
		route = append(route, step)
	}

	return Route{route}
}

// getBestNode TODO extract the Node with shortes path from pool, it is a substitute for using a priority queue
func getBestNode(pool map[*Node]bool, distMap map[*Node]int) *Node {
	bestDist := 100000
	var bestNode *Node
	for Node := range pool {
		if distMap[Node] < bestDist {
			bestNode = Node
			bestDist = distMap[Node]
		}
	}
	return bestNode
}
