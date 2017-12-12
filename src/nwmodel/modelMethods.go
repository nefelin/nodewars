package nwmodel

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

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
		TeamName:  p.TeamName,
		Health:    response.passed(),
		MaxHealth: len(response.PassFail),
	}
}

// NewTeam creates a new team with color/name color
func NewTeam(n teamName) *team {
	return &team{
		Name:    n,
		players: make(map[*Player]bool),
		maxSize: 10,
	}
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
		Modules:     modules,
	}
}

func newNodeMap() nodeMap {
	return nodeMap{
		Nodes:    make([]*node, 0),
		POEs:     make(map[nodeID]bool),
		diameter: 0,
		radius:   1000,
	}
}

// Instantiation with values -----------------------------------------------------------------

// NewDefaultModel Generic game model
func NewDefaultModel() *GameModel {
	m := newRandMap(10)
	t := makeDummyTeams()
	p := make(map[playerID]*Player)
	poes := make(map[playerID]*node)

	aChan := make(chan Message, 100)

	gm := &GameModel{
		Map:     m,
		Teams:   t,
		Players: p,
		// Routes:  r,
		POEs:      poes,
		languages: getLanguages(),
		aChan:     aChan,
	}

	go actionConsumer(gm)

	return gm
}

func makeDummyTeams() map[teamName]*team {
	teams := make(map[teamName]*team)
	teams["red"] = NewTeam("red")
	teams["blue"] = NewTeam("blue")
	return teams
}

func newRandMap(n int) *nodeMap {
	rand.Seed(time.Now().UTC().UnixNano())
	nodeCount := n
	newMap := newNodeMap()

	for i := 0; i < nodeCount; i++ {
		newMap.addNodes(NewNode())
	}

	for i := 0; i < nodeCount; i++ {
		if i < nodeCount-1 {
			newMap.connectNodes(i, i+1)
		}

		for j := 0; j < rand.Intn(2); j++ {
			newMap.connectNodes(i, rand.Intn(nodeCount))
		}

	}

	newMap.initAllNodes()

	newMap.initPoes(2)
	// for len(newMap.POEs) < 2 {
	// 	newMap.addPoes(rand.Intn(nodeCount))
	// }
	log.Printf("newRandMap nodes: %v", newMap.Nodes)
	return &newMap
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

	newMap.addPoes(1, 10)

	return &newMap
}

// GameModel methods --------------------------------------------------------------------------

func (gm *GameModel) resetMap(m *nodeMap) {
	gm.Map = m

	// Tell everyone to clear their maps
	for _, p := range gm.Players {
		p.Outgoing <- Message{
			Type:   "graphReset",
			Sender: serverStr,
			Data:   "",
		}
	}

	// Clear map specific data:
	gm.POEs = make(map[playerID]*node)
	for _, t := range gm.Teams {
		t.poe = nil
	}
	// log.Println(gm.POEs)

	// send our new state
	gm.broadcastState()
}

func (gm *GameModel) broadcastState() {
	for _, player := range gm.Players {
		state := gm.calcState(player)
		player.Outgoing <- graphStateMsg(state)
	}
}

func (gm *GameModel) calcState(p *Player) string {
	stateMsg, err := json.Marshal(gm)

	if err != nil {
		log.Println(err)
	}
	return string(stateMsg)
}

func (gm *GameModel) broadcastAlertFlash(color string) {
	// TODO abstract this to messages
	for _, player := range gm.Players {
		player.Outgoing <- Message{
			Type:   "alertFlash",
			Sender: serverStr,
			Data:   color,
		}
	}
}

// send a pseudoServer message to all players
func (gm *GameModel) psBroadcast(msg Message) {
	msg.Sender = "pseudoServer"

	for _, player := range gm.Players {
		player.Outgoing <- msg
	}
}

// broadcast to all but one player
func (gm *GameModel) psBroadcastExcept(p *Player, msg Message) {
	msg.Sender = "pseudoServer"

	for _, player := range gm.Players {
		//skip if it's our player
		if player == p {
			continue
		}
		player.Outgoing <- msg
	}
}

func (gm *GameModel) setTeamPoe(t *team, ni nodeID) error {
	if t.poe != nil {
		return fmt.Errorf("Team %s already has a point of entry at node '%d'", t.Name, t.poe.ID)
	}

	if !gm.Map.nodeExists(ni) {
		return fmt.Errorf("Node '%v' does not exist", ni)
	}
	avail, ok := gm.Map.POEs[ni]
	if !ok {
		return fmt.Errorf("Node '%v' is not a valid point of entry", ni)
	}

	if !avail {
		return errors.New("That point of entry is already taken")
	}

	node := gm.Map.Nodes[ni]
	// set the teams poe
	t.poe = node

	// set all teams players poes
	for player := range t.players {
		gm.setPlayerPOE(player, ni)
	}

	// mark the spot as taken

	gm.Map.POEs[ni] = false

	return nil
}

// func (gm *GameModel) setPlayerName(p *Player, n string) error {

// 	// check to see if name is in use
// 	for _, player := range gm.Players {
// 		if player.Name == n {
// 			return errors.New("Name '" + n + "' already in use")
// 		}
// 	}
// 	// if not, set it and return no error
// 	p.Name = n
// 	return nil
// }

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
	if p.TeamName != "" {
		gm.Teams[p.TeamName].removePlayer(p)
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
	} else if p.TeamName != "" {
		return errors.New("Already on the " + p.TeamName + " team")
	} else if team.isFull() {
		return errors.New("team: " + tn + " is full")
	}

	t := gm.Teams[tn]
	t.addPlayer(p)
	if t.poe != nil {
		gm.setPlayerPOE(p, t.poe.ID)
	}

	return nil
}

func (gm *GameModel) tryConnectPlayerToNode(p *Player, n nodeID) (*route, error) {

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

	routeNodes := gm.Map.routeToNode(gm.Teams[p.TeamName], source, target)
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
	n.addPlayer(p)
	return p.Route
	// return gm.Routes[p.ID]
}

func (gm *GameModel) breakConnection(p *Player) {
	// if _, exists := gm.Routes[p.ID]; exists {
	if p.Route != nil {
		p.Route.Endpoint.removePlayer(p)
		// delete(gm.Routes, p.ID)
		p.Route = nil

		// TODO this should only happen if connection was severed non-voluntarily...
		p.Outgoing <- psError(errors.New("Connection interrupted!"))
	}

	// detach from any slots
	if p.slotNum != -1 {
		p.slotNum = -1
	}
}

func (gm *GameModel) evalTrafficForTeam(n *node, t *team) {
	// if the module no longer supports routing for this modules team
	if !n.allowsRoutingFor(t) {
		for _, player := range gm.Players {
			// check each player who is on team's route
			if player.TeamName == t.Name {
				// and if it contained that node, break the players connection
				if _, ok := player.Route.containsNode(n); ok {
					gm.breakConnection(player)
				}
			}
		}
		// if this is a POE, announce that teams elimination
		if t.poe == n {
			gm.psBroadcast(psAlert(fmt.Sprintf("(%s) has been ELIMINATED!", t.Name)))

		}

	}
}

func (gm *GameModel) setLanguage(p *Player, l string) error {

	for lang := range gm.languages {
		if l == lang {
			p.language = strings.ToLower(l)

			p.Outgoing <- Message{
				Type:   "languageState",
				Sender: "server",
				Data:   p.language,
			}
			return nil
		}
	}
	return fmt.Errorf("'%v' is not a supported in this match. Use 'langs' to list available languages")
}

// module methods -------------------------------------------------------------------------

func (m module) isFriendlyTo(t *team) bool {
	if m.TeamName == t.Name {
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

// addConnection is reciprocol
func (n *node) addConnection(m *node) {
	// if the connection already exists, ignore
	for _, nID := range n.Connections {
		if m.ID == nID {
			return
		}
	}

	if m.ID == n.ID {
		return
	}

	log.Printf("Adding: %d to %d's connections", m.ID, n.ID)
	n.Connections = append(n.Connections, m.ID)
	m.Connections = append(m.Connections, n.ID)
}

func (n *node) allowsRoutingFor(t *team) bool {
	// t == nil means we don't care...
	if t == nil {
		return true
	}

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

	//remove module from node and empty slot
	delete(n.Modules, slot.module.id)
	slot.module = nil

	// assign new task to slot
	slot.challenge = getRandomTest()
	return nil

}

// TODO this is broken:
func (n *node) addPlayer(p *Player) {
	n.playersHere = append(n.playersHere, p)
}

func (n *node) removePlayer(p *Player) {
	n.playersHere = cutPFromSlice(n.playersHere, p)
}

// helper function for removing player from slice of players
func cutPFromSlice(s []*Player, p *Player) []*Player {
	for i, thisP := range s {
		if p == thisP {
			// swaps the last element with the found element and returns with the last element cut
			s[len(s)-1], s[i] = s[i], s[len(s)-1]
			return s[:len(s)-1]
		}
	}
	// log.Printf("CutPlayer returning: %v", s)
	log.Println("Player not found in slice")
	return s
}

// nodeMap methods -----------------------------------------------------------------------------

func (m *nodeMap) initAllNodes() {

	// initialize each node's slots
	for _, node := range m.Nodes {
		node.initSlots()
	}

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

// func (m *nodeMap) openPoes() int {
// 	var count int
// 	for _, open := range m.POEs {
// 		if open {
// 			count++
// 		}
// 	}
// 	return count
// }

func (m *nodeMap) addPoes(ns ...nodeID) {
	for _, id := range ns {
		// skip bad ids
		if !m.nodeExists(id) {
			continue
		}
		// make an available POE for each nodeID passed
		m.POEs[id] = true
	}
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
	for i, v := range ordRem {
		// if there're enough nodes of that remoteness, seen if we can place n poes in those remotenesses
		if len(remMap[v]) >= n {
			// add a poe at each node of that remoteness
			for _, node := range remMap[v] {
				m.addPoes(node.ID)
			}
			break
		}
		log.Printf("We have %d nodes of remoteness %v", len(remMap[ordRem[i]]), ordRem[i])
	}

	// if all else fails, assign at random
	// only ensuring distance

}

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

// player methods -------------------------------------------------------------------------------
// TODO this is in the wrong place
func NewPlayer(ws *websocket.Conn) *Player {
	ret := &Player{
		ID:       playerIDCount,
		Name:     "",
		Socket:   ws,
		Outgoing: make(chan Message),
		slotNum:  -1,
	}

	// log.Println("New player created, setting language...")
	playerIDCount++
	return ret
}

func (p *Player) prompt() string {
	promptEndChar := ">"
	prompt := fmt.Sprintf("(%s)", p.name())
	if p.TeamName != "" {
		prompt += fmt.Sprintf(":%s:", p.TeamName)
	}
	if p.Route != nil {
		prompt += fmt.Sprintf("@n%d", p.Route.Endpoint.ID)
	}
	if p.slotNum != -1 {
		prompt += fmt.Sprintf(":s%d", p.slotNum)
	}
	prompt += fmt.Sprintf("[%s]", p.language)

	prompt += promptEndChar

	return prompt
}

// TODO refactor this, modify how slots are tracked, probably with IDs
func (p *Player) slot() *modSlot {
	if p.Route == nil || p.slotNum < 0 || p.slotNum > len(p.Route.Endpoint.slots) {
		return nil
	}

	return p.Route.Endpoint.slots[p.slotNum]
}

func (p *Player) name() string {
	for p.Name == "" {
		p.Name = "player_" + strconv.Itoa(p.ID)
	}

	return p.Name
}

func (p Player) hasTeam() bool {
	if p.TeamName == "" {
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

func (r route) length() int {
	return len(r.Nodes)
}

// team methods -------------------------------------------------------------------------------
func (t team) isFull() bool {
	if len(t.players) < t.maxSize {
		return false
	}
	return true
}

func (t *team) broadcast(msg Message) {
	msg.Sender = "pseudoServer"

	for player := range t.players {
		player.Outgoing <- msg
	}
}

func (t *team) addPlayer(p *Player) {
	t.players[p] = true
	p.TeamName = t.Name
}

func (t *team) removePlayer(p *Player) {
	delete(t.players, p)
	p.TeamName = ""
}

// func (t *team) setPoe(n *node) {

// }

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
	return fmt.Sprintf("( <player> {Name: %v, team: %v} )", p.name(), p.TeamName)
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
	return fmt.Sprintf("[%s] [%s] [%s] [%d/%d]", m.TeamName, m.builder, m.language, m.Health, m.MaxHealth)
}

func (n node) forMsg() string {

	slotList := ""
	for i, slot := range n.slots {
		slotList += "\n" + strconv.Itoa(i) + ":" + slot.forMsg()
	}

	return fmt.Sprintf("NodeID: %v\nMemory Slots:%v", n.ID, slotList)
}

func (m modSlot) forMsg() string {
	switch {
	case m.module != nil:
		return "(" + m.module.forMsg() + ")"
	default:
		return "( -empty- )"
	}
}

// func (m modSlot) forProbe() string {
// 	var header string
// 	switch {
// 	case m.module != nil:
// 		header = "( " + m.module.forMsg() + " )\n"
// 	default:
// 		header = "( -empty- )\n"
// 	}
// 	// task := "Task:\n" + m.challenge.Description
// 	return header //+ task

// }

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
