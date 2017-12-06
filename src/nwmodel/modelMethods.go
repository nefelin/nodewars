package nwmodel

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

// Initialization methods ------------------------------------------------------------------

func newModSlot() *modSlot {
	// get random challenge,

	// assign id
	return &modSlot{
		challenge: getRandomTest(),
	}
}

// creates a new module by p based on the results from response in language l
func newModule(p *Player, response ChallengeResponse, lang string) *module {
	id := moduleIDCount
	moduleIDCount++

	return &module{
		id: id,
		// testID:    testID,
		language:  lang,
		builder:   p.name(),
		Team:      p.Team,
		Health:    response.passed(),
		MaxHealth: len(response.PassFail),
	}
}

// NewTeam creates a new team with color/name color
func NewTeam(n teamName) *team {
	return &team{n, make(map[*Player]bool), 2}
}

// NewNode ...
func NewNode() *node {
	id := nodeIDCount
	nodeIDCount++

	connections := make([]int, 0)
	modules := make(map[modID]*module)

	return &node{
		ID:          id,
		Connections: connections,
		// Capacity:    3,
		Modules: modules,
		// Traffic:          make([]*Player, 0),
		// POE:              make([]*Player, 0),
		// ConnectedPlayers: make([]*Player, 0),
	}
}

func newNodeMap() nodeMap {
	return nodeMap{make([]*node, 0)}
}

// Instantiation with values ------------------------------------------------------------------

// NewDefaultModel Generic game model
func NewDefaultModel() *GameModel {
	m := newDefaultMap()
	t := makeDummyTeams()
	p := make(map[playerID]*Player)
	// r := make(map[playerID]*route)
	poes := make(map[playerID]*node)

	return &GameModel{
		Map:     m,
		Teams:   t,
		Players: p,
		// Routes:  r,
		POEs: poes,
	}
}

func makeDummyTeams() map[teamName]*team {
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

	// create module slots based on connectivity of node
	for _, node := range newMap.Nodes {
		node.initSlots()
	}

	return &newMap
}

// GameModel methods --------------------------------------------------------------------------

// RegisterPlayer adds a new player to our world model
func (gm *GameModel) RegisterPlayer(ws *websocket.Conn) *Player {
	// create player with this websocket
	newP := newPlayer(ws)

	// add this player to our registry
	gm.Players[newP.ID] = newP
	return newP
}

func (gm *GameModel) broadcastState() {
	for _, player := range gm.Players {
		player.outgoing <- calcStateMsgForPlayer(player)
	}
}

// send a pseudoServer message to all players
func (gm *GameModel) psBroadcast(msg Message) {
	msg.Sender = "pseudoServer"

	for _, player := range gm.Players {
		player.outgoing <- msg
	}
}

func (gm *GameModel) setPlayerName(p *Player, n string) error {

	// check to see if name is in use
	for _, player := range gm.Players {

		if player.Name == n {
			return errors.New("Name '" + n + "' already in use")
		}
	}
	// if not, set it and return no error
	p.Name = n
	return nil
}

func (gm *GameModel) setPlayerPOE(p *Player, n nodeID) bool {
	// TODO move this node validity check to a nodeMap method
	// if nodeID is valid

	if gm.Map.nodeExists(n) {
		gm.POEs[p.ID] = gm.Map.Nodes[n]
		return true
	}

	return false
}

// RemovePlayer ...
func (gm *GameModel) RemovePlayer(p *Player) error {
	if _, ok := gm.Players[p.ID]; !ok {
		return errors.New("player '" + p.Name + "' is not registered")
	}
	if p.Team != nil {
		p.Team.removePlayer(p)
	}

	// clean up POE
	delete(gm.POEs, p.ID)

	// Clean up route
	// delete(gm.Routes, p.ID)

	// Clean up player
	delete(gm.Players, p.ID)

	return nil
}

func (gm *GameModel) assignPlayerToTeam(p *Player, tn teamName) error {
	// log.Printf("assignPlayerToTeam, player: %v", p)
	if team, ok := gm.Teams[tn]; !ok {
		return errors.New("The team '" + tn + "' does not exist")
	} else if p.Team != nil {
		return errors.New("Already on the " + p.Team.Name + " team")
	} else if team.isFull() {
		return errors.New("team: " + tn + " is full")
	}

	gm.Teams[tn].addPlayer(p)
	return nil
}

func (gm *GameModel) tryConnectPlayerToNode(p *Player, n nodeID) (*route, error) {

	// break any pre-existing connection before connecting elsewhere
	gm.breakConnection(p)

	// TODO report errors here
	source, poeOK := gm.POEs[p.ID]

	// log.Printf("source: %v, poeOK: %v, gm.POEs: %v", source, poeOK, gm.POEs)
	if !poeOK {
		return nil, errors.New("No point of entry")
	}

	if !gm.Map.nodeExists(n) {
		return nil, fmt.Errorf("%v is not a valid node", n)
	}

	// log.Printf("player %v attempting to connect to node %v from POE %v", p.Name, n, gm.POEs[p.ID].ID)

	target := gm.Map.Nodes[n]

	routeNodes := gm.Map.routeToNode(p, source, target)
	if routeNodes != nil {
		// log.Println("Successful Connect")
		// log.Printf("Route to target: %v", routeNodes)
		route := gm.establishConnection(p, routeNodes, target) // This should add player traffic to each intermediary and establish a connection on n
		return route, nil
	}
	// log.Println("Cannot Connect")
	return nil, errors.New("No route exists")

}

// TODO should this have gm as receiver? there's no need but makes sense syntactically
func (gm *GameModel) establishConnection(p *Player, routeNodes []*node, n *node) *route {
	// set's players route to the route generated via routeToNode
	// gm.Routes[p.ID] = &route{Endpoint: n, Nodes: routeNodes}
	p.Route = &route{Endpoint: n, Nodes: routeNodes}
	return p.Route
	// return gm.Routes[p.ID]
}

func (gm *GameModel) breakConnection(p *Player) {
	// if _, exists := gm.Routes[p.ID]; exists {
	if p.Route != nil {
		// delete(gm.Routes, p.ID)
		p.Route = nil
	}

	// detach from any slots
	if p.slotNum != -1 {
		p.slotNum = -1
	}
}

// module methods -------------------------------------------------------------------------

func (m module) isFriendlyTo(t *team) bool {
	if m.Team == t {
		return true
	}
	return false
}

// modSlot methods -------------------------------------------------------------------------

func (m modSlot) isFull() bool {
	if m.module == nil {
		return false
	}
	return true
}

// node methods -------------------------------------------------------------------------------

func (n *node) initSlots() {
	for range n.Connections {
		newSlot := newModSlot()
		n.slots = append(n.slots, newSlot)
	}
}

// TODO deprecate this for modslot approach
func (n node) capacity() int {
	return len(n.Connections)
}

func (n node) isFull() bool {
	if len(n.Modules) > n.capacity()-1 {
		return true
	}
	return false
}

// addConnection is reciprocol
func (n *node) addConnection(m *node) {
	n.Connections = append(n.Connections, m.ID)
	m.Connections = append(m.Connections, n.ID)
}

func (n *node) allowsRoutingFor(t *team) bool {
	for _, module := range n.Modules {
		if module.isFriendlyTo(t) {
			return true
		}
	}
	return false
}

// TODO fix awckward redundancy of modules and slots
func (n *node) addModule(m *module, slotIndex int) error {

	slot := n.slots[slotIndex]
	if slot.module == nil {
		n.Modules[m.id] = m
		slot.module = m
		return nil
	}
	return errors.New("Slot not empty")
}

func (n *node) removeModule(slotIndex int) error {
	if slotIndex < 0 || slotIndex > len(n.slots)-1 {
		return errors.New("No valid attachment")
	}

	slot := n.slots[slotIndex]
	log.Printf("removeModule slot: %v", slot)

	if slot.module == nil {
		return errors.New("Slot is empty")
	}

	// track old team so we can evaluate traffic after
	oldModsTeam := slot.module.Team

	//remove module from node and empty slot
	delete(n.Modules, slot.module.id)
	slot.module = nil

	// evalTrafficForTeam makes sure all players that were routing through this node are still able to do so
	n.evalTrafficForTeam(oldModsTeam)

	// assign new task to slot
	slot.challenge = getRandomTest()
	return nil

}

func (n *node) evalTrafficForTeam(t *team) {
	// if the module no longer supports routing for this modules team
	if !n.allowsRoutingFor(t) {
		for _, player := range gm.Players {
			// check each player who is on team's route
			if player.Team == t {
				// and if it contained that node, break the players connection
				if _, ok := player.Route.containsNode(n); ok {
					gm.breakConnection(player)
				}
			}
		}
	}
}

// helper function for removing item from slice
// func cutPlayer(s []*Player, p *Player) []*Player {
// 	for i, thisP := range s {
// 		if p == thisP {
// 			// swaps the last element with the found element and returns with the last element cut
// 			s[len(s)-1], s[i] = s[i], s[len(s)-1]
// 			return s[:len(s)-1]
// 		}
// 	}
// 	log.Printf("CutPlayer returning: %v", s)
// 	return s
// }

// nodeMap methods -----------------------------------------------------------------------------

func (m *nodeMap) addNodes(ns ...*node) {
	for _, node := range ns {
		m.Nodes = append(m.Nodes, node)
	}
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

// routeToNode uses vanilla dijkstra's (vanilla for now) algorithm to find node path
// TODO get code review on this. I think I'm maybe not getting optimal route
func (m *nodeMap) routeToNode(p *Player, source, target *node) []*node {

	if source.allowsRoutingFor(p.Team) {
		// if we're connecting to our POE, return a route which is only our POE
		if source == target {
			route := make([]*node, 1)
			route[0] = source
			return route
		}

		unchecked := make(map[*node]bool) // TODO this should be a priority queue for efficiency
		dist := make(map[*node]int)
		prev := make(map[*node]*node)

		seen := make(map[*node]bool)
		tocheck := make([]*node, 1)
		tocheck[0] = source
		for len(tocheck) > 0 {
			thisNode := tocheck[0]
			tocheck = tocheck[1:]

			// log.Printf("this: %v", thisNode)
			if thisNode.allowsRoutingFor(p.Team) {
				unchecked[thisNode] = true
				dist[thisNode] = 1000
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

		// log.Printf("unchecked %v", unchecked)

		dist[source] = 0

		for len(unchecked) > 0 {
			thisNode := getBestNode(unchecked, dist)

			delete(unchecked, thisNode)

			if m.nodesTouch(thisNode, target) {
				prev[target] = thisNode
				route := constructPath(prev, target)
				// log.Println("Found target!")
				return route
			}

			for _, cNode := range m.nodesConnections(thisNode) {
				// TODO refactor to take least risky routes by weighing against vulnerability to enemy connection
				alt := dist[thisNode] + 1
				if alt < dist[cNode] {
					dist[cNode] = alt
					prev[cNode] = thisNode
				}
			}
		}
	} else {
		// log.Println("POE Blocked")
	}
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

// player methods -------------------------------------------------------------------------------
// TODO this is in the wrong place
func newPlayer(ws *websocket.Conn) *Player {
	ret := &Player{
		ID:       playerIDCount,
		Name:     "",
		language: "python",
		socket:   ws,
		outgoing: make(chan Message),
		slotNum:  -1,
	}

	ret.setLanguage("python")
	playerIDCount++
	return ret
}

func (p *Player) setLanguage(l string) {
	p.language = strings.ToLower(l)
}

// TODO refactor this, modify how slots are tracked, probably with IDs
func (p *Player) slot() *modSlot {
	if p.Route == nil || p.slotNum < 0 || p.slotNum > len(p.Route.Endpoint.slots) {
		return nil
	}

	return p.Route.Endpoint.slots[p.slotNum]
}

func (p *Player) name() string {
	rand.Seed(int64(p.ID))
	for p.Name == "" {

		propName := "player_" + strconv.Itoa(rand.Intn(100))
		err := gm.setPlayerName(p, propName)

		if err != nil {
			log.Println(err)
		}
	}

	return p.Name
}

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

// route methods --------------------------------------------
func (r route) containsNode(n *node) (int, bool) {
	for i, node := range r.Nodes {
		if n == node {
			return i, true
		}
	}
	return 0, false
}

// team methods -------------------------------------------------------------------------------
func (t team) isFull() bool {
	if len(t.players) < t.MaxSize {
		return false
	}
	return true
}

func (t *team) broadcast(msg Message) {
	msg.Sender = "pseudoServer"

	for player := range t.players {
		player.outgoing <- msg
	}
}

func (t *team) addPlayer(p *Player) {
	t.players[p] = true
	p.Team = t
}

func (t *team) removePlayer(p *Player) {
	delete(t.players, p)
	p.Team = nil
}

// Stringers ----------------------------------------------------------------------------------

func (n node) String() string {
	return fmt.Sprintf("( <node> {ID: %v, Connections:%v, Modules:%v} )", n.ID, n.Connections, n.modIDs())
}

func (n node) modIDs() []modID {
	ids := make([]modID, len(n.Modules))
	i := 0
	for _, mod := range n.Modules {
		ids[i] = mod.id
		i++
	}
	return ids
}

func (t team) String() string {
	var playerList []string
	for player := range t.players {
		playerList = append(playerList, string(player.Name))
	}
	return fmt.Sprintf("( <team> {Name: %v, Players:%v} )", t.Name, playerList)
}

func (p Player) String() string {
	return fmt.Sprintf("( <player> {Name: %v, team: %v} )", p.name(), p.Team)
}

func (r route) String() string {
	nodeCount := len(r.Nodes)
	nodeList := make([]string, nodeCount)

	for i, node := range r.Nodes {
		// this loop is a little funny because we are reversing the order of the node list
		// it's reverse ordered in the data structure but to be human readable we'd like
		// the list to read from source to target
		nodeList[nodeCount-i-1] = strconv.Itoa(node.ID)
	}

	return fmt.Sprintf("( <route> {Endpoint: %v, Through: %v} )", r.Endpoint.ID, strings.Join(nodeList, ", "))
}

func (m module) forMsg() string {
	return fmt.Sprintf("[%v] [%v] [%v]", m.Team.Name, m.language, m.builder)
}

func (n node) forMsg() string {

	slotList := ""
	for i, slot := range n.slots {
		slotList += strconv.Itoa(i) + ":" + slot.forMsg()
	}

	return fmt.Sprintf("NodeID: %v\nMemory Slots:\n%v", n.ID, slotList)
}

func (m modSlot) forMsg() string {
	switch {
	case m.module != nil:
		return "( " + m.module.forMsg() + " )\n"
	default:
		return "( -empty- )\n"
	}
}

func (m modSlot) forProbe() string {
	var header string
	switch {
	case m.module != nil:
		header = "( " + m.module.forMsg() + " )\n"
	default:
		header = "( -empty- )\n"
	}
	task := "Task:\n" + m.challenge.Description
	return header + task

}

func (r route) forMsg() string {
	nodeCount := len(r.Nodes)
	nodeList := make([]string, nodeCount)

	for i, node := range r.Nodes {
		// this loop is a little funny because we are reversing the order of the node list
		// it's reverse ordered in the data structure but to be human readable we'd like
		// the list to read from source to target
		nodeList[nodeCount-i-1] = strconv.Itoa(node.ID)
	}
	return fmt.Sprintf("(Endpoint: %v, Through: %v)", r.Endpoint.ID, strings.Join(nodeList, ", "))
}

func (c ChallengeResponse) String() string {
	ret := ""
	for k, v := range c.PassFail {
		ret += fmt.Sprintf("(in: %v, out: %v)", k, v)
	}
	return ret
}

// func (m modSlot) String() string {

// }
