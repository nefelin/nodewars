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
	name    string
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

// tries to set initial language to one of these defaults before picking first available
var langDefaults = []string{
	"python",
	"javascript",
	"golang",
	"c++",
}

// init methods:

// NewDefaultModel Generic game model
func NewDefaultModel(name string) *GameModel {
	m := newRandMap(10)
	p := make(map[playerID]*Player)
	// poes := make(map[playerID]*node)

	aChan := make(chan nwmessage.Message, 100)

	gm := &GameModel{
		name:    name,
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

	// fmt.Println("Supported Languages:")
	// for key := range gm.languages {
	// 	fmt.Println(key)
	// }

	err := gm.addTeams(makeDummyTeams())
	if err != nil {
		fmt.Println(err)
	}

	// go actionConsumer(gm)

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
func (gm *GameModel) GetPlayers() []*Player {
	list := make([]*Player, len(gm.Players))
	var i int
	for _, p := range gm.Players {
		list[i] = p
		i++
	}
	return list
}

func (gm *GameModel) Recv(msg ClientMessage) error {
	// gm.aChan <- msg
	gm.parseCommand(msg)
	return nil
}

func (gm *GameModel) Name() string {
	// gm.aChan <- msg
	return gm.name
}

func (gm *GameModel) Type() string {
	// gm.aChan <- msg
	return ""
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

		poeIndex := rand.Intn(len(nodes))
		poeNode := nodes[poeIndex]

		poeNode.Feature.dummyClaim(t.Name, "MIN")
		err := t.addPoe(poeNode)
		if err != nil {
			panic(err)
		}

		nodes = append(nodes[:poeIndex], nodes[poeIndex+1:]...)
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

	// convert map to list
	teamList := make([]*team, 0)
	for _, team := range gm.Teams {
		if len(team.poes) > 0 {
			teamList = append(teamList, team)
		}
	}

	// scramble list
	for i := range teamList {
		j := rand.Intn(i + 1)
		teamList[i], teamList[j] = teamList[j], teamList[i]
	}

	// fmt.Printf("<trailingTeam> gm.Teams %v\n", gm.Teams)
	for _, team := range teamList {
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
		// clear previous list of powered nodes
		t.powered = nil

		if n.hasMachineFor(t) {
			var foundPower bool

			for poe := range t.poes {
				if gm.Map.routeToNode(t, n, poe) != nil {
					foundPower = true
					break
				}
			}

			if foundPower {
				n.powerMachines(t.Name, true)
				t.powered = append(t.powered, n)
			} else {
				n.powerMachines(t.Name, false)
			}
		}
	}
}

func (gm *GameModel) playersAt(n *node) []*Player {
	players := make([]*Player, 0)
	for _, p := range gm.Players {
		if p.location() == n {
			players = append(players, p)
		}
	}
	return players
}

func (gm *GameModel) detachOtherPlayers(p *Player, msg string) {
	if p.currentMachine() == nil {
		log.Panic("Player is not attached to a machine")
	}

	for _, player := range gm.playersAt(p.Route.Endpoint()) {
		if player != p {
			comment := gm.languages[player.language].CommentPrefix

			editMsg := fmt.Sprintf("%s %s", comment, msg)

			player.Outgoing <- nwmessage.PsAlert(fmt.Sprintf("You have been detached from the machine at %s", player.currentMachine().address))
			player.Outgoing <- nwmessage.EditState(editMsg)

			player.breakConnection(true)

		}
	}
}

func (gm *GameModel) tryClaimMachine(p *Player, mac *machine, response GradedResult, fType feature.Type) {
	node := p.Route.Endpoint()
	solutionStrength := response.passed()

	var hostile bool
	var friendly bool

	if !mac.isNeutral() {
		if !mac.belongsTo(p.TeamName) {
			hostile = true
		} else {
			friendly = true
		}
	}

	if hostile {
		switch {
		case mac.Health == mac.MaxHealth:
			p.Outgoing <- nwmessage.PsError(fmt.Errorf("Current solution of %d/%d is the best possible so you cannot steal this machine,\nuse 'reset' instead of 'make' to remove opponent's solution.", mac.Health, mac.MaxHealth))
			return

		case solutionStrength < mac.Health:
			p.Outgoing <- nwmessage.PsError(fmt.Errorf("Solution (%d/%d) too weak to install, need at least %d/%d to steal", response.passed(), len(response.Grades), mac.Health+1, mac.MaxHealth))
			return

		case solutionStrength == mac.Health:
			p.Outgoing <- nwmessage.PsAlert(fmt.Sprintf("You need to pass one more test to steal,\nbut your %d/%d is enough to reset this machine.\nKeep trying if you think you can do\nbetter or type 'reset' to proceed", solutionStrength, mac.MaxHealth))
			return
		}

	}

	var oldTeam *team
	var oldAllowed bool

	// track old owner to evaluate traffic after module loss
	newTeam := gm.Teams[p.TeamName]
	if hostile {
		oldTeam = gm.Teams[mac.TeamName]
	}

	// track whether node allowed routing for active player before refactor
	allowed := node.hasMachineFor(newTeam)
	if hostile {
		oldAllowed = node.hasMachineFor(oldTeam)
	}

	// TODO I think we need to do same for old team, but test first

	// refactor module to new owner and health
	mac.TeamName = p.TeamName
	mac.language = p.language
	mac.Health = solutionStrength

	if mac.Type == feature.None {
		mac.Type = fType
	} else if mac.Type == feature.POE {
		err := newTeam.addPoe(node)
		if err != nil {
			panic(err)
		}

		if hostile {
			oldTeam.remPoe(node)
			gm.psBroadcast(nwmessage.PsAlert(fmt.Sprintf("%s team has lost a Point of Entry to %s!", oldTeam.Name, newTeam.Name)))
			if len(oldTeam.poes) < 1 {
				gm.psBroadcast(nwmessage.PsAlert(fmt.Sprintf("%s team has been eliminated!", oldTeam.Name)))
			}
		}
	}

	// evaluate routing of player trffic through node
	if hostile {
		gm.evalTrafficForTeam(node, oldTeam)
	}

	// if routing status has changed, recalculate powered nodes
	if node.hasMachineFor(newTeam) != allowed {
		gm.calcPoweredNodes(newTeam)
	}
	if hostile {
		if node.hasMachineFor(oldTeam) != oldAllowed {
			gm.calcPoweredNodes(oldTeam)
		}
	}

	// recalculate coin production
	gm.updateCoinPerTick(newTeam)
	if hostile {
		gm.updateCoinPerTick(oldTeam)
	}

	// map alert
	gm.pushActionAlert(p.TeamName, node.ID)

	// update map
	gm.broadcastState()

	// do terminal messaging
	if hostile {
		gm.psBroadcastExcept(p, nwmessage.PsAlert(fmt.Sprintf("%s of (%s) stole a (%s) machine in node %d", p.GetName(), p.TeamName, oldTeam.Name, node.ID)))
		p.Outgoing <- nwmessage.PsSuccess(fmt.Sprintf("You stole (%v)'s machine, new machine health: %d/%d", oldTeam.Name, mac.Health, mac.MaxHealth))
	} else if friendly {
		gm.psBroadcastExcept(p, nwmessage.PsAlert(fmt.Sprintf("%s of (%s) refactored a friendly machine in node %d", p.GetName(), p.TeamName, node.ID)))
		p.Outgoing <- nwmessage.PsSuccess(fmt.Sprintf("Refactored friendly machine to %d/%d [%s]", mac.Health, mac.MaxHealth, mac.language))
	} else {
		gm.psBroadcastExcept(p, nwmessage.PsAlert(fmt.Sprintf("%s of (%s) constructed a machine in node %d", p.GetName(), p.TeamName, node.ID)))
		p.Outgoing <- nwmessage.PsSuccess(fmt.Sprintf("Solution installed in [%s], Health: %d/%d", mac.language, mac.Health, mac.MaxHealth))
	}
}

func (gm *GameModel) tryResetMachine(p *Player, mac *machine, r GradedResult) {
	node := p.Route.Endpoint()
	solutionStrength := r.passed()

	if mac.isNeutral() {
		p.Outgoing <- nwmessage.PsError(errors.New("Machine is already neutral"))
		return
	}

	if solutionStrength < mac.Health {
		p.Outgoing <- nwmessage.PsError(fmt.Errorf(
			"Solution too weak: %d/%d, need %d/%d to remove",
			r.passed(), len(r.Grades), mac.Health, mac.MaxHealth))
		return
	}

	// track old owner to evaluate traffic after module loss
	oldTeam := gm.Teams[mac.TeamName]

	// track whether node allowed routing for active player before refactor
	allowed := node.hasMachineFor(oldTeam)

	// reset the machine
	mac.reset()

	// evaluate routing of player trffic through node
	gm.evalTrafficForTeam(node, oldTeam)

	// if routing status has changed, recalculate powered nodes
	if node.hasMachineFor(oldTeam) != allowed {
		gm.calcPoweredNodes(oldTeam)
	}

	// recalculate teams processsing power
	gm.updateCoinPerTick(oldTeam)

	// if machine was poe, remove team poe pointer
	if mac.Type == feature.POE {
		err := oldTeam.remPoe(node)
		if err != nil {
			panic(err)
		}

		gm.psBroadcast(nwmessage.PsAlert(fmt.Sprintf("%s team has lost a Point of Entry!", oldTeam.Name)))
		if len(oldTeam.poes) < 1 {
			gm.psBroadcast(nwmessage.PsAlert(fmt.Sprintf("%s team has been eliminated!", oldTeam.Name)))
		}
	}

	// kick out other players working at this slot
	// gm.detachOtherPlayers(p, fmt.Sprintf("%s removed the module you were working on", p.name))
	gm.pushActionAlert(p.TeamName, node.ID)
	gm.broadcastState()

	// terminal messaging
	gm.psBroadcastExcept(p, nwmessage.PsAlert(fmt.Sprintf("%s of (%s) reset a (%s) machine in node %d", p.GetName(), p.TeamName, oldTeam.Name, node.ID)))
	p.Outgoing <- nwmessage.PsSuccess("Machine reset")
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
	for _, t := range gm.Teams {
		t.poes = nil
	}

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
		playerLoc = p.Route.Endpoint().ID
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
		player.Outgoing <- nwmessage.PsPrompt(player.Prompt())
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
		player.Outgoing <- nwmessage.PsPrompt(p.Prompt())
	}
}

func (gm *GameModel) setPlayerName(p *Player, n string) error {

	// check to see if name is in use
	for _, player := range gm.Players {
		if player.name == n {
			return errors.New("Name '" + n + "' already in use")
		}
	}

	p.name = n
	return nil
}

// AddPlayer ...
func (gm *GameModel) AddPlayer(p *Player) error {
	if _, ok := gm.Players[p.ID]; ok {
		return errors.New("player '" + p.GetName() + "' is already in this game")
	}

	p.inGame = true
	gm.Players[p.ID] = p
	gm.pendingAlerts[p.ID] = make([]alert, 0) // make alerts slot for new player

	supportedLangs := make([]string, len(gm.languages))
	var i int
	for lang := range gm.languages {
		supportedLangs[i] = lang
		i++
	}

	p.Outgoing <- nwmessage.LangSupportState(supportedLangs)

	var defaultLanguage string
	for _, lang1 := range langDefaults {
		if defaultLanguage != "" {
			break
		}
		for _, lang2 := range supportedLangs {
			// fmt.Printf("comparing %s to %s\n", lang1, lang2)
			if lang1 == lang2 {
				defaultLanguage = lang1
			}
		}
	}
	if defaultLanguage == "" {
		defaultLanguage = supportedLangs[0]
	}

	gm.setLanguage(p, defaultLanguage)

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
	p.breakConnection(false)
	p.inGame = false

	p.Outgoing <- nwmessage.LangSupportState([]string{})
	p.Outgoing <- nwmessage.GraphReset()

	return nil
}

func (gm *GameModel) assignPlayerToTeam(p *Player, tn teamName) error {
	// log.Printf("assignPlayerToTeam, player: %v", p)
	if team, ok := gm.Teams[tn]; !ok {
		return errors.New("'" + tn + "' team does not exist")
	} else if p.TeamName != "" {
		return errors.New("You're already on the " + p.TeamName + " team")
	} else if team.isFull() {
		return errors.New("The " + tn + " team is full")
	} else if len(team.poes) < 1 {
		return errors.New("The " + tn + " team is dead (poe was eliminated)")
	}

	t := gm.Teams[tn]
	t.addPlayer(p)
	// if t.poe != nil {
	// 	gm.setPlayerPOE(p, t.poe.ID)
	// }

	return nil
}

func (gm *GameModel) tryConnectPlayerToNode(p *Player, n nodeID) (*route, error) {

	// log.Printf("source: %v, poeOK: %v, gm.POEs: %v", source, poeOK, gm.POEs)
	team := gm.Teams[p.TeamName]

	if len(team.poes) < 1 {
		return nil, errors.New("No point of entry")
	}

	target := gm.Map.getNode(n)
	if target == nil {
		return nil, fmt.Errorf("%v is not a valid node", n)
	}

	for source := range team.poes {
		// log.Printf("player %v attempting to connect to node %v from POE %v", p.GetName(), n, gm.POEs[p.ID].ID)
		routeNodes := gm.Map.routeToNode(gm.Teams[p.TeamName], source, target)
		if routeNodes != nil {
			// log.Println("Successful Connect")
			// log.Printf("Route to target: %v", routeNodes)
			p.breakConnection(false)
			route, err := gm.establishConnection(p, routeNodes)

			if err != nil {
				return nil, err
			}

			return route, nil
		}

	}
	// log.Println("Cannot Connect")
	return nil, errors.New("No route exists")
}

// TODO should this have gm as receiver? there's no need but makes sense syntactically
func (gm *GameModel) establishConnection(p *Player, routeNodes []*node) (*route, error) {
	// set's players route to the route generated via routeToNode
	// gm.Routes[p.ID] = &route{Endpoint: n, Nodes: routeNodes}
	r := &route{Nodes: routeNodes, player: p}

	// make sure we're not blocked by any firewalls:
	for _, n := range r.Nodes {
		if n.Feature.Type == feature.Firewall {
			if n.machinesFor(p.TeamName) <= gm.trafficCount(n, p.TeamName) {
				return nil, fmt.Errorf("Connection refused (firewall at node %d)", n.ID)
			}
		}
	}

	p.Route = r
	return p.Route, nil
	// return gm.Routes[p.ID]
}

func (gm *GameModel) trafficCount(n *node, t teamName) int {
	var count int

	for p := range gm.Teams[t].players {
		if p.Route != nil {
			if p.Route.runsThrough(n) || p.Route.Endpoint() == n {
				count++
			}
		}
	}
	return count
}

func (gm *GameModel) evalTrafficForTeam(n *node, t *team) {
	// if the module no longer supports routing for this modules team
	if !n.hasMachineFor(t) {
		for _, player := range gm.Players {
			// check each player who is on team's route
			if player.TeamName == t.Name {
				// and if it contained that node, break the players connection
				if player.Route != nil {
					if player.Route.runsThrough(n) {
						player.breakConnection(true)
					}
				}
			}
		}
	}
}

func (gm *GameModel) setLanguage(p *Player, l string) error {
	_, ok := gm.languages[l]

	if !ok {
		return fmt.Errorf("'%v' is not a supported in this match. Use 'langs' to list available languages")
	}

	mac := p.currentMachine()
	if mac != nil && !mac.isNeutral() && !mac.belongsTo(p.TeamName) && mac.language != l {
		return errors.New("Can't change language while attached to a hostile machine")
	}

	p.language = l

	p.Outgoing <- nwmessage.EditLangState(p.language)
	return nil
}
