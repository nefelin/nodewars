package nwmodel

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/gorilla/websocket"
)

// Initialization methods ------------------------------------------------------------------

func newModuleBy(p *Player) module {
	id := moduleCount
	moduleCount++

	return module{
		ID:         id,
		TestID:     0,
		LanguageID: 0,
		Builder:    p,
	}
}

func newGameState() *gameState {
	// make a list of all team names
	teams := make([]string, 0)
	for name := range gm.Teams {
		teams = append(teams, name)
	}

	players := make([]*Player, 0)
	for player := range gm.Players {
		// Only add players to the stat object that have names and teams
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
	modules := make(map[modID]module)

	return &node{
		ID:               id,
		Connections:      connections,
		Size:             3,
		Modules:          modules,
		Traffic:          make([]*Player, 0),
		POE:              make([]*Player, 0),
		ConnectedPlayers: make([]*Player, 0),
	}
}

// hiLo is a helper function that lets new edge sort node pairs for its ID scheme
func hiLo(a, b nodeID) (nodeID, nodeID) {
	if a > b {
		return a, b
	}
	return b, a
}

func newEdge(s, t nodeID) *edge {
	hi, lo := hiLo(s, t)

	hiStr := strconv.Itoa(hi)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	loStr := strconv.Itoa(lo)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	id := hiStr + "e" + loStr
	// id := edgeCount
	// edgeCount++

	return &edge{id, s, t, make([]*Player, 0)}
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
		Name:           "",
		Team:           nil,
		PointOfEntry:   -1,
		socket:         ws,
		outgoing:       make(chan Message),
		NodeConnection: -1,
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

func (gm *GameModel) setPlayerPOE(p *Player, n nodeID) bool {
	// TODO move this node validity check to a nodeMap method
	// if nodeID is valid and player has no other POE
	if n > -1 && n < nodeCount && p.PointOfEntry == -1 {

		// // if player already has POE, clear old POE
		// if p.PointOfEntry > -1 {
		// 	gm.Map.Nodes[p.PointOfEntry].removePOE(p)
		// }

		log.Printf("setting %v's poe to %v", p.Name, n)

		p.PointOfEntry = n
		gm.Map.Nodes[n].addPOE(p)
		gm.Map.Nodes[n].addModule(newModuleBy(p))
		return true
	}

	return false
}

// RemovePlayer ...
func (gm *GameModel) RemovePlayer(p *Player) error {
	if _, ok := gm.Players[p]; !ok {
		return errors.New("player '" + p.Name + "' is not registered")
	}

	if p.Team != nil {
		p.Team.removePlayer(p)
	}

	if p.PointOfEntry > -1 {
		log.Printf("removing poe for: %v", p.Name)
		gm.Map.Nodes[p.PointOfEntry].removePOE(p)
	}

	if p.NodeConnection > -1 {
		gm.Map.Nodes[p.NodeConnection].removePlayerConnection(p)
	}

	for _, node := range p.route {
		node.removeTraffic(p)
	}

	delete(gm.Players, p)
	return nil
}

func (gm *GameModel) assignPlayerToTeam(p *Player, tn teamName) error {
	if team, ok := gm.Teams[tn]; !ok {
		return errors.New("The team: " + tn + " does not exist")
	} else if p.Team != nil {
		return errors.New(p.Name + " is alread a member of team: " + tn)
	} else if team.isFull() {
		return errors.New("team: " + tn + " is full")
	}

	gm.Teams[tn].addPlayer(p)
	return nil
}

func (gm *GameModel) tryConnectPlayerToNode(p *Player, n nodeID) bool {
	log.Printf("player %v attempting to connect to node %v from POE %v", p.Name, n, p.PointOfEntry)

	// if player is connected elsewhere, break that first, regardless of success of this attempt
	if p.NodeConnection != -1 {
		gm.breakConnection(p)
	}

	// TODO handle player connecting to own POE

	source := gm.Map.Nodes[p.PointOfEntry]
	target := gm.Map.Nodes[n]

	route := gm.Map.routeToNode(p, source, target)
	if route != nil {
		log.Println("Successful Connect")
		gm.establishConnection(p, route, target) // This should add player traffic to each intermediary and establish a connection on n
		return true
	}
	log.Println("Cannot Connect")
	return false
	// return errors.New("Player: " + p.Name + " cannot reach node: " + strconv.Itoa(n))
}

// TODO should this have gm as receiver? there's no need but makes sense syntactically
func (gm *GameModel) establishConnection(p *Player, route []*node, n *node) {
	// set's players route to the route generated via routeToNode
	p.route = route

	// sets players nodeConnection to target's id

	// adds player to target nodes playerConnections
	n.addPlayerConnection(p)

	// adds traffic to every intermediary node
	for _, node := range route {
		node.addTraffic(p)
	}
}

func (gm *GameModel) breakConnection(p *Player) {
	if p.NodeConnection > -1 {
		for _, node := range p.route {
			node.removeTraffic(p)
		}

		p.route = nil // how to clear :?

		gm.Map.Nodes[p.NodeConnection].removePlayerConnection(p)

		p.NodeConnection = -1
	}
}

// module methods -------------------------------------------------------------------------

func (m module) isFriendlyTo(p *Player) bool {
	if m.Builder.Team == p.Team {
		return true
	}
	return false
}

// node methods -------------------------------------------------------------------------------

func (n *node) addConnection(m *node) {
	n.Connections = append(n.Connections, m.ID)
	m.Connections = append(m.Connections, n.ID)
	// log.Printf("node %v added connection to node %v. All connections: %v", n, m, n.Connections)
}

func (n *node) allowsRoutingFor(p *Player) bool {
	for _, module := range n.Modules {
		if module.isFriendlyTo(p) {
			return true
		}
	}
	return false
}

func (n *node) addModule(m module) {
	// TODO QUESTION do I need to return bool since failure is possible?
	if len(n.Modules) < n.Size {
		n.Modules[m.ID] = m
	}
}

func (n *node) removeModule(m module) {
	delete(n.Modules, m.ID)
}

func (n *node) addTraffic(p *Player) {
	n.Traffic = append(n.Traffic, p)
}

func (n *node) removeTraffic(p *Player) {
	cutPlayer(n.Traffic, p)
}

func (n *node) addPlayerConnection(p *Player) {
	n.ConnectedPlayers = append(n.ConnectedPlayers, p)
}

func (n *node) removePlayerConnection(p *Player) {
	n.Traffic = cutPlayer(n.Traffic, p)
}

func (n *node) addPOE(p *Player) {

	n.POE = append(n.POE, p)
}

func (n *node) removePOE(p *Player) {
	log.Printf("Should remove poe for: %v", p.Name)
	n.POE = cutPlayer(n.POE, p)
}

// helper function for removing item from slice
func cutPlayer(s []*Player, p *Player) []*Player {
	for i, thisP := range s {
		if p == thisP {
			// swaps the last element with the found element and returns with the last element cut
			s[len(s)-1], s[i] = s[i], s[len(s)-1]
			return s[:len(s)-1]
		}
	}
	log.Printf("CutPlayer returning: %v", s)
	return s
}

// edge methods

func (e *edge) addTraffic(p *Player) {
	e.Traffic = append(e.Traffic, p)
}

func (e *edge) removeTraffic(p *Player) {
	cutPlayer(e.Traffic, p)
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

// nodesConnections takes one of the maps nodes and converts its connections (in the form of nodeIDs) into pointers to actual node objects
func (m *nodeMap) nodesConnections(n *node) []*node {
	res := make([]*node, 0)
	for _, nodeID := range n.Connections {
		res = append(res, m.Nodes[nodeID])
	}

	return res
}

func (m *nodeMap) nodesAreTouching(n1, n2 *node) bool {
	// for every one of n1's connections
	for _, connectedNode := range m.nodesConnections(n1) {
		// if it is n2, return true
		if connectedNode == n2 {
			return true
		}
	}
	return false
}

// routeToNode uses vanilla dijkstra's (vanilla for now) algorithm to find node path
func (m *nodeMap) routeToNode(p *Player, source, target *node) []*node {
	unchecked := make(map[*node]bool) // this should be a priority queue for efficiency
	dist := make(map[*node]int)
	prev := make(map[*node]*node)

	for _, node := range m.Nodes {
		// Only do these if node is friendly to player
		if node.allowsRoutingFor(p) {
			dist[node] = 10000
			unchecked[node] = true
		}
	}

	dist[source] = 0

	for len(unchecked) > 0 {
		thisNode := getBestNode(unchecked, dist)

		delete(unchecked, thisNode)

		if m.nodesAreTouching(thisNode, target) {
			route := constructPath(prev, target)
			log.Println("Found target!")
			log.Printf("%v", route)
			return route
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
// 	if p.Team == nil {
// 		if !t.isFull() {
// 			t.addPlayer(p)
// 		} else {
// 			// tell player team is full, TODO centralize control messages
// 			p.outgoing <- Message{"teamFull", "server", t.Name}
// 		}
// 	} else {
// 		p.outgoing <- Message{"error", "server", "you are already a member of " + p.Team.Name}
// 	}
// }

func (p Player) hasTeam() bool {
	if p.Team == nil {
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
	if len(t.players) < t.MaxSize {
		return false
	}
	return true
}

func (t *team) broadcast(msg Message) {
	for player := range t.players {
		player.outgoing <- msg
	}
}

func (t *team) addPlayer(p *Player) {
	t.players[p] = true
	p.Team = t

	// Tell client they've joined model shouldn't handle messaging, fix
	p.outgoing <- Message{
		Type:   "teamAssign",
		Sender: "server",
		Data:   t.Name,
	}
}

func (t *team) removePlayer(p *Player) {
	delete(t.players, p)
	p.Team = nil

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
	for player := range t.players {
		playerList = append(playerList, player.Name)
	}
	return fmt.Sprintf("<team> (Name: %v, Players:%v)", t.Name, playerList)
}

func (p Player) String() string {
	return fmt.Sprintf("<player> Name: %v, team: %v", p.Name, p.Team)
}
