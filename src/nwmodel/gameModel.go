package nwmodel

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"nwmessage"
	"time"

	"feature"
)

// GameModel holds all state information
type GameModel struct {
	Map     *nodeMap             `json:"map"`
	Teams   map[teamName]*team   `json:"teams"`
	Players map[playerID]*Player `json:"players"`
	// POEs          map[playerID]*node   `json:"poes"`
	PointGoal     float32 `json:"pointGoal"`
	languages     map[string]Language
	aChan         chan nwmessage.Message
	running       bool //running should replace mapLocked
	pendingAlerts map[playerID][]alert

	// timelimit should be able to set a timelimit and count points at the end
}

// init methods:

// NewDefaultModel Generic game model
func NewDefaultModel() *GameModel {
	m := newRandMap(10)
	p := make(map[playerID]*Player)
	// poes := make(map[playerID]*node)

	aChan := make(chan nwmessage.Message, 100)

	gm := &GameModel{
		Map:     m,
		Teams:   make(map[string]*team),
		Players: p,
		// Routes:  r,
		// POEs:          poes,
		languages:     getLanguages(),
		aChan:         aChan,
		PointGoal:     1000,
		pendingAlerts: make(map[playerID][]alert),
	}

	err := gm.addTeams(makeDummyTeams())
	if err != nil {
		fmt.Println(err)
	}

	go actionConsumer(gm)

	return gm
}

func makeDummyTeams() []*team {
	teams := make([]*team, 2)
	teams[0] = NewTeam("red")
	teams[1] = NewTeam("blue")
	return teams
}

// GameModel methods --------------------------------------------------------------------------

// fulfill Room interface
func (gm *GameModel) GetPlayers() map[playerID]*Player {
	return gm.Players
}

func (gm *GameModel) Recv(msg nwmessage.Message) {
	gm.aChan <- msg
}

// Addteams can only be called once. After addteams is called unused poes are removed
func (gm *GameModel) addTeams(teams []*team) error {
	nodes := gm.Map.collectEmptyPoes()

	// add each team to gm.Teams
	for _, t := range teams {
		gm.Teams[t.Name] = t
	}

	// find a poe for each team
	for i, t := range teams {
		if len(nodes) < 1 {
			// TODO alternative approach would be to add poes when we run out...
			return fmt.Errorf("Ran out of poes with %d teams left to place", len(teams)-i)
		}

		myPoe := rand.Intn(len(nodes))
		gm.setTeamPoe(t, nodes[myPoe].ID)
		nodes = append(nodes[:myPoe], nodes[myPoe+1:]...)
	}

	// remove any leftovers
	if len(nodes) > 0 {
		for _, node := range nodes {
			node.Feature.Type = feature.None
		}
	}

	return nil
}

// trailingTeam should hand back either smallest team or currently losing team, depending on game settings TODO
func (gm *GameModel) trailingTeam() string {
	var tt *team
	// fmt.Printf("<trailingTeam> gm.Teams %v\n", gm.Teams)
	for _, team := range gm.Teams {
		if tt == nil {
			tt = team
			continue
		}
		if len(team.players) < len(tt.players) {
			tt = team
		}
	}
	// fmt.Printf("<trailingTeam> returning %v\n", tt)
	return tt.Name
}

func (gm *GameModel) updateCoinPerTick(t *team) {
	// for each node in t.powered add power for each module
	// if a slot is producing, set slot.Processing = true so we can animate this

	// reset
	t.coinPerTick = 0
	// go through each node
	for _, node := range t.powered {
		// store the powerpermod of that node
		t.coinPerTick += node.coinProduction(t.Name)
	}
}

func (gm *GameModel) calcPoweredNodes(t *team) {
	for _, n := range gm.Map.Nodes {
		// clear previus powered status
		t.powered = nil
		// t.powered[n] = false
		if n.hasMachineFor(t) {
			if gm.Map.routeToNode(t, n, t.poe) != nil {
				// power machines
				n.powerMachines(t.Name, true)

				// add to our list for production calc
				t.powered = append(t.powered, n)
			} else {
				// depower machines
				n.powerMachines(t.Name, false)
			}
		}

	}
}

func (gm *GameModel) playersAt(n *node) []*Player {
	players := make([]*Player, len(n.playersHere))
	for i, pID := range n.playersHere {
		players[i] = gm.Players[pID]
	}
	return players
}

func (gm *GameModel) detachOtherPlayers(p *Player, msg string) {
	if p.currentMachine == nil {
		log.Panic("Player is not attached to a machine")
	}

	for _, player := range gm.playersAt(p.Route.Endpoint) {
		if player != p {
			comment := gm.languages[player.language].CommentPrefix

			editMsg := fmt.Sprintf("%s %s", comment, msg)

			player.Outgoing <- nwmessage.PsAlert(fmt.Sprintf("You have been detached from machine %d", player.slotNum))
			player.Outgoing <- nwmessage.EditState(editMsg)
			player.slotNum = -1
		}
	}
}

func (gm *GameModel) claimMachine(p *Player, r ExecutionResult) {
	mac := p.currentMachine()

	if mac == nil {
		log.Panic("Player is not attached to a machine")
		return
	}

	if mac.TeamName != "" {
		log.Panic(errors.New("Machine is not neutral"))
		return
	}

	n := p.Route.Endpoint
	t := gm.Teams[p.TeamName]

	// track whether node allowed routing for active player before building
	allowed := n.hasMachineFor(t)

	mac.claim(p, r)

	// if routing status has changed, recalculate powered nodes
	if p.Route.Endpoint.hasMachineFor(t) != allowed {
		gm.calcPoweredNodes(t)
	}

	// recalculate this teams processsing power
	gm.updateCoinPerTick(t)

	// kick out other players working at this mac
	gm.detachOtherPlayers(p, fmt.Sprintf("%s took control of the machine you were working on", p.name))
}

func (gm *GameModel) refactorMachine(m *machine, p *Player, newHealth int) {
	// track old owner to evaluate traffic after module loss
	oldTeam := gm.Teams[m.TeamName]

	pTeam := gm.Teams[p.TeamName]
	// track whether node allowed routing for active player before refactor
	allowed := p.Route.Endpoint.hasMachineFor(pTeam)

	// refactor module to new owner and health
	m.TeamName = p.TeamName
	m.Health = newHealth

	// evaluate routing of player trffic through node
	gm.evalTrafficForTeam(p.Route.Endpoint, oldTeam)

	// if routing status has changed, recalculate powered nodes
	if p.Route.Endpoint.hasMachineFor(pTeam) != allowed {
		gm.calcPoweredNodes(pTeam)
	}

	// recalculate old teams processsing power
	gm.updateCoinPerTick(oldTeam)

	// recalculate this teams processsing power
	gm.updateCoinPerTick(pTeam)
}

func (gm *GameModel) resetMachine(p *Player) {
	mac := p.currentMachine()

	if mac == nil {
		log.Panic("Player is not attached to a machine")
		return
	}

	if mac.TeamName == "" {
		log.Panic("Machine is already neutral")
		return
	}

	// track old owner to evaluate traffic after module loss
	oldTeam := gm.Teams[mac.TeamName]

	// track whether node allowed routing for active player before refactor
	allowed := p.Route.Endpoint.hasMachineFor(oldTeam)

	// remove the module
	err := p.Route.Endpoint.resetMachine(p.slotNum)
	if err != nil {
		log.Panic(err)
	}

	// evaluate routing of player trffic through node
	gm.evalTrafficForTeam(p.Route.Endpoint, oldTeam)

	// if routing status has changed, recalculate powered nodes
	if p.Route.Endpoint.hasMachineFor(oldTeam) != allowed {
		gm.calcPoweredNodes(oldTeam)
	}

	// recalculate teams processsing power
	gm.updateCoinPerTick(oldTeam)

	// kick out other players working at this slot
	gm.detachOtherPlayers(p, fmt.Sprintf("%s removed the module you were working on", p.name))
}

func (gm *GameModel) tickScheduler() {
	for gm.running == true {
		<-time.After(1 * time.Second)
		gm.tick()
	}
}

// this is a naive approach, would be more performant to deal in deltas and only use tick to increment total, not recalculate rate
// TODO approach should be that on any module gain or loss that teams procPow is recalculated
// this entails making a pool of all nodes connected to POE and running the below logic
func (gm *GameModel) tick() {
	// reset each teams ProcPow
	// for _, team := range gm.Teams {
	// 	team.ProcPow = 0
	// }

	// // go through each node
	// for _, node := range gm.Map.Nodes {
	// 	// store the powerpermod of that node
	// 	modVal := node.getPowerPerMod()

	// 	// look at each slot
	// 	for _, slot := range node.Slots {
	// 		// for each module give the owner team appropriate power boost
	// 		if slot.Module != nil {
	// 			gm.Teams[slot.Module.TeamName].ProcPow += modVal
	// 		}
	// 	}
	// }

	// advance each teams VicPoints
	winners := make([]string, 0)

	for _, team := range gm.Teams {
		gm.updateCoinPerTick(team)
		team.VicPoints += team.coinPerTick
		if team.VicPoints >= gm.PointGoal {
			winners = append(winners, team.Name)
		}
		// gm.psBroadcast(nwmessage.PsAlert(fmt.Sprintf("Team %s has completed %d calculations", team.Name, team.VicPoints)))
	}

	if len(winners) > 0 {
		for _, name := range winners {
			gm.psBroadcast(nwmessage.PsAlert(fmt.Sprintf("Team %s wins!", name)))
		}
		gm.stopGame()
	}

	gm.broadcastScore()
}

func (gm *GameModel) startGame() {
	if !gm.running {
		gm.running = true
		gm.psBroadcast(nwmessage.PsAlert("All teams have POEs. Game is starting!"))
		// start go routine to handle ticking
		go gm.tickScheduler()
	}
}

func (gm *GameModel) stopGame() {
	gm.running = false

	// ticking goroutine should auto collapse when running is false
}

func (gm *GameModel) resetMap(m *nodeMap) {
	gm.Map = m

	// Tell everyone to clear their maps
	for _, p := range gm.Players {
		p.Outgoing <- nwmessage.GraphReset()
	}

	// Clear map specific data:
	// gm.POEs = make(map[playerID]*node)
	for _, t := range gm.Teams {
		t.poe = nil
	}
	// log.Println(gm.POEs)

	// send our new state
	gm.broadcastState()
}

func (gm *GameModel) broadcastScore() {
	for _, p := range gm.Players {
		p.Outgoing <- nwmessage.ScoreState(gm.packScores())
	}
}

func (gm *GameModel) makeRouteMap() *trafficMap {
	// collect routes, TODO redundat loop
	traffic := newTrafficMap()

	for _, p := range gm.Players {
		if p.Route != nil {
			traffic.addRoute(p.Route, p.TeamName)
		}
	}

	return traffic
}

func (gm *GameModel) broadcastState() {
	routeMap := gm.makeRouteMap()

	for _, p := range gm.Players {
		// TODO feels super hackey to have to pass in routemap but was quickest solution for now.
		state := gm.calcState(p, routeMap)
		p.Outgoing <- nwmessage.GraphState(state)
	}
}

func (gm *GameModel) broadcastGraphReset() {
	for _, p := range gm.Players {
		p.Outgoing <- nwmessage.GraphReset()
	}
}

func (gm *GameModel) packScores() string {
	scoreMsg, err := json.Marshal(gm.Teams)
	if err != nil {
		log.Println(err)
	}
	return string(scoreMsg)
}

// calcState takes a player argument on the assumption that at some point we'll want to show different states to different players
func (gm *GameModel) calcState(p *Player, tMap *trafficMap) string {

	// calculate player location
	var playerLoc nodeID
	if p.Route == nil {
		playerLoc = -1
	} else {
		playerLoc = p.Route.Endpoint.ID
	}

	// compose state message
	state := stateMessage{
		nodeMap:    gm.Map,
		Alerts:     gm.pendingAlerts[p.ID],
		PlayerLoc:  playerLoc,
		trafficMap: tMap,
	}

	// fmt.Printf("state: %v", state)

	stateMsg, err := json.Marshal(state)
	// fmt.Printf("statemessage: %s\n", stateMsg)
	// clear dumped alerts
	gm.pendingAlerts[p.ID] = nil

	// fmt.Printf("\nState Message %v\n", stateMsg)
	if err != nil {
		log.Println(err)
	}

	return string(stateMsg)
}

func (gm *GameModel) pushActionAlert(color string, location nodeID) {
	for k := range gm.pendingAlerts {
		gm.pendingAlerts[k] = append(gm.pendingAlerts[k], alert{color, location})
	}
}

// send a pseudoServer message to all players
func (gm *GameModel) psBroadcast(msg nwmessage.Message) {
	msg.Sender = "pseudoServer"

	for _, player := range gm.Players {
		player.Outgoing <- msg
	}
}

// broadcast to all but one player
func (gm *GameModel) psBroadcastExcept(p *Player, msg nwmessage.Message) {
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
		return fmt.Errorf("Node '%d' does not exist", ni)
	}

	node := gm.Map.Nodes[ni]
	if node.Feature.Type != feature.POE {
		return fmt.Errorf("No Point of Entry feature at Node, '%d'", ni)
	}

	if node.Feature.TeamName != "" {
		return errors.New("That point of entry is already taken")
	}

	// set the teams poe
	t.poe = node

	node.Feature.dummyClaim(t.Name, "MIN")
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
// 	p.GetName() = n
// 	return nil
// }

// func (gm *GameModel) setPlayerPOE(p *Player, n nodeID) bool {
// 	// TODO move this node validity check to a nodeMap method
// 	// if nodeID is valid

// 	if gm.Map.nodeExists(n) {

// 		gm.POEs[p.ID] = gm.Map.Nodes[n]

// 		return true
// 	}

// 	return false
// }

// AddPlayer ...
func (gm *GameModel) AddPlayer(p *Player) error {
	if _, ok := gm.Players[p.ID]; ok {
		return errors.New("player '" + p.GetName() + "' is already in this game")
	}

	p.inGame = true
	gm.Players[p.ID] = p
	gm.pendingAlerts[p.ID] = make([]alert, 0) // make alerts slot for new player

	gm.setLanguage(p, "python")

	// send initiall map state

	p.Outgoing <- nwmessage.GraphReset()
	routeMap := gm.makeRouteMap()
	p.Outgoing <- nwmessage.GraphState(gm.calcState(p, routeMap))

	// send initial prompt state
	p.SendPrompt()
	return nil
}

// RemovePlayer ...
func (gm *GameModel) RemovePlayer(p *Player) error {
	fmt.Printf("<gm.RemovePlayer> Removing player, %s\n", p.name)
	if _, ok := gm.Players[p.ID]; !ok {
		return errors.New("player '" + p.GetName() + "' is not registered")
	}

	// remove player infor from gamemodel
	if p.TeamName != "" {
		gm.Teams[p.TeamName].removePlayer(p)
	}

	// delete(gm.POEs, p.ID)

	delete(gm.Players, p.ID)

	delete(gm.pendingAlerts, p.ID)

	// remove game infor from player object

	p.TeamName = ""
	p.Route = nil
	p.slotNum = -1
	p.inGame = false

	p.Outgoing <- nwmessage.GraphReset()

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
	// if t.poe != nil {
	// 	gm.setPlayerPOE(p, t.poe.ID)
	// }

	return nil
}

func (gm *GameModel) tryConnectPlayerToNode(p *Player, n nodeID) (*route, error) {

	// TODO report errors here
	// source, poeOK := gm.POEs[p.ID]

	// log.Printf("source: %v, poeOK: %v, gm.POEs: %v", source, poeOK, gm.POEs)
	source := gm.Teams[p.TeamName].poe
	if source == nil {
		return nil, errors.New("No point of entry")
	}

	if !gm.Map.nodeExists(n) {
		return nil, fmt.Errorf("%v is not a valid node", n)
	}

	// log.Printf("player %v attempting to connect to node %v from POE %v", p.GetName(), n, gm.POEs[p.ID].ID)

	target := gm.Map.Nodes[n]

	routeNodes := gm.Map.routeToNode(gm.Teams[p.TeamName], source, target)
	if routeNodes != nil {
		// log.Println("Successful Connect")
		// log.Printf("Route to target: %v", routeNodes)
		p.breakConnection(false)
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

// func (gm *GameModel) breakConnection(p *Player, alert bool) {
// 	// if _, exists := gm.Routes[p.ID]; exists {
// 	if p.Route == nil {
// 		// log.Panic("No route for player")
// 		return
// 	}

// 	p.Route.Endpoint.removePlayer(p)
// 	p.slotNum = -1
// 	p.Route = nil

// 	if alert {
// 		p.Outgoing <- nwmessage.PsError(errors.New("Connection interrupted!"))
// 	}
// }

func (gm *GameModel) evalTrafficForTeam(n *node, t *team) {
	// if the module no longer supports routing for this modules team
	if !n.hasMachineFor(t) {
		for _, player := range gm.Players {
			// check each player who is on team's route
			if player.TeamName == t.Name {
				// and if it contained that node, break the players connection
				if player.Route != nil {
					if _, ok := player.Route.containsNode(n); ok {
						player.breakConnection(true)
					}
				}
			}
		}
		// if this is a POE, announce that teams elimination
		if t.poe == n {
			gm.psBroadcast(nwmessage.PsAlert(fmt.Sprintf("(%s) has been ELIMINATED!", t.Name)))

		}

	}
}

func (gm *GameModel) setLanguage(p *Player, l string) error {
	_, ok := gm.languages[l]

	if !ok {
		return fmt.Errorf("'%v' is not a supported in this match. Use 'langs' to list available languages")
	}

	p.language = l

	p.Outgoing <- nwmessage.Message{
		Type:   "languageState",
		Sender: "server",
		Data:   p.language,
	}
	return nil
}
