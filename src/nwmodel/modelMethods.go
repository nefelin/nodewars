package nwmodel

import (
	"errors"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

// Initialization methods ------------------------------------------------------------------

func newGameState() *gameState {
	// make a list of all team names
	teams := make([]string, 0)
	for name := range gm.Teams {
		teams = append(teams, name)
	}

	players := make([]*Player, 0)
	for player := range gm.Players {
		if player.hasName() && player.hasTeam() {
			players = append(players, player)
		}
	}

	return &gameState{
		Map:           *gm.Map,
		Teams:         teams,
		Players:       players,
		CurrentEvents: make([]gameEvent, 0),
	}
}

// NewTeam creates a new team with color/name color
func NewTeam(n teamName) *team {
	return &team{n, make(map[*Player]bool), 2}
}

// NewNode ...
func NewNode() *node {
	id := nodeCount
	nodeCount++

	connections := make([]int, 0)
	modules := make([]module, 0)

	return &node{
		ID:          id,
		Connections: connections,
		Size:        3,
		Modules:     modules}
}

func newEdge(s, t nodeID) *edge {
	id := edgeCount
	edgeCount++

	return &edge{id, s, t}
}

func newNodeMap() nodeMap {
	return nodeMap{make(map[nodeID]*node), make(map[edgeID]*edge)}
}

// Instantiation with values ------------------------------------------------------------------

// NewDefaultModel Generic game model
func NewDefaultModel() *GameModel {
	m := newDefaultMap()
	t := makeDummyTeams()
	p := make(map[*Player]bool)
	e := make([]*gameEvent, 0)

	return &GameModel{
		Map:           m,
		Teams:         t,
		Players:       p,
		CurrentEvents: e,
	}
}

func makeDummyTeams() map[string]*team {
	teams := make(map[teamName]*team)
	teams["red"] = NewTeam("red")
	teams["blue"] = NewTeam("blue")
	return teams
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

// GameModel methods --------------------------------------------------------------------------

// RegisterPlayer adds a new player to our world model
func (gm *GameModel) RegisterPlayer(ws *websocket.Conn) *Player {
	newPlayer := &Player{
		Name:         "",
		team:         nil,
		PointOfEntry: -1,
		socket:       ws,
		outgoing:     make(chan Message),
	}

	gm.Players[newPlayer] = true
	return newPlayer
}

func (gm *GameModel) setPlayerName(p *Player, n string) error {
	// check to see if name is in use
	for player := range gm.Players {
		if player.Name == n {
			return errors.New("Name '" + n + "' already in use")
		}
	}

	// if not set it and return no error
	p.Name = n
	return nil
}

// RemovePlayer ...
func (gm *GameModel) RemovePlayer(p *Player) error {
	if _, ok := gm.Players[p]; !ok {
		return errors.New("player '" + p.Name + "' is not registered")
	}

	if p.team != nil {
		p.team.removePlayer(p)
	}

	delete(gm.Players, p)
	return nil
}

func (gm *GameModel) assignPlayerToTeam(p *Player, tn teamName) error {
	if team, ok := gm.Teams[tn]; !ok {
		return errors.New("The team: " + tn + " does not exist")
	} else if p.team != nil {
		return errors.New(p.Name + " is alread a member of team: " + tn)
	} else if team.isFull() {
		return errors.New("team: " + tn + " is full")
	}

	gm.Teams[tn].addPlayer(p)
	return nil
}

func (gm *GameModel) connectPlayerToNode(p *Player, n nodeID) bool {
	log.Printf("player %v attempting to connect to node %v from POE %v", p.Name, n, p.PointOfEntry)
	route := gm.Map.routeToNode(gm.Map.Nodes[p.PointOfEntry], gm.Map.Nodes[n])
	if route != nil {
		log.Println("Successful Connect")
		return true
	}
	return false
	// return errors.New("Player: " + p.Name + " cannot reach node: " + strconv.Itoa(n))
}

// node methods -------------------------------------------------------------------------------
func (n *node) addConnection(m *node) {
	n.Connections = append(n.Connections, m.ID)
	m.Connections = append(m.Connections, n.ID)
	// log.Printf("node %v added connection to node %v. All connections: %v", n, m, n.Connections)
}

// nodeMap methods -----------------------------------------------------------------------------

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

// helper functions for routeToNode ------------------------------------------------------------
// constructPath takes the routes discovered via routeToNode and the endpoint (target) and creates a slice of the correct path, note order is still reversed and path contains source but not target node
func constructPath(prevMap map[*node]*node, t *node) []*node {
	// log.Println(prevMap)

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

// player methods -------------------------------------------------------------------------------

// func (p *Player) joinTeam(t *team) {
// 	if p.team == nil {
// 		if !t.isFull() {
// 			t.addPlayer(p)
// 		} else {
// 			// tell player team is full, TODO centralize control messages
// 			p.outgoing <- Message{"teamFull", "server", t.Name}
// 		}
// 	} else {
// 		p.outgoing <- Message{"error", "server", "you are already a member of " + p.team.Name}
// 	}
// }

func (p Player) hasTeam() bool {
	if p.team == nil {
		return false
	}
	return true
}

func (p Player) hasName() bool {
	if p.Name == "" {
		return false
	}
	return true
}

// player attempting to node, use routeToNode (dijkstra) to determine possible path (if any)
// func (p *player) connectToNode(n nodeID) bool {
// 	log.Printf("player %v attempting to connect to node %v from POE %v", p.Name, n, p.PointOfEntry)
// 	if gameMap.routeToNode(gameMap.Nodes[p.PointOfEntry], gameMap.Nodes[n]) != nil {
// 		return true
// 	}
// 	return false
// }

// team methods -------------------------------------------------------------------------------
func (t team) isFull() bool {
	if len(t.Players) < t.MaxSize {
		return false
	}
	return true
}

func (t *team) broadcast(msg Message) {
	for player := range t.Players {
		player.outgoing <- msg
	}
}

func (t *team) addPlayer(p *Player) {
	t.Players[p] = true
	p.team = t

	// Tell client they've joined model shouldn't handle messaging, fix
	p.outgoing <- Message{
		Type:   "teamAssign",
		Sender: "server",
		Data:   t.Name,
	}
}

func (t *team) removePlayer(p *Player) {
	delete(t.Players, p)
	p.team = nil

	// Notify client model shouldn't handle messaging, fix
	p.outgoing <- Message{
		Type:   "teamUnassign",
		Sender: "server",
		Data:   t.Name,
	}
}

// Stringers ----------------------------------------------------------------------------------

func (n node) String() string {
	return fmt.Sprintf("<(node) ID: %v, Connections:%v, Modules:%v>", n.ID, n.Connections, n.Modules)
}

func (t team) String() string {
	var playerList []string
	for player := range t.Players {
		playerList = append(playerList, player.Name)
	}
	return fmt.Sprintf("<team> (Name: %v, Players:%v)", t.Name, playerList)
}

func (p Player) String() string {
	return fmt.Sprintf("<player> Name: %v, team: %v", p.Name, p.team)
}
