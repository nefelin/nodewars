package model

import (
	"challenges"
	"encoding/json"
	"errors"
	"feature"
	"fmt"
	"log"
	"math/rand"
	"model/machines"
	"model/modes"
	"model/node"
	"model/player"
	"model/statemessage"
	"nwmessage"
	"room"
	"sort"
	"strings"
	"time"
	"timer"
)

type playerSet = map[*player.Player]bool

// GameModel holds all state information
type GameModel struct {
	name          string
	Map           *node.Map                          `json:"map"`
	Teams         map[teamName]*team                 `json:"teams"`
	Players       map[player.PlayerID]*player.Player `json:"players"`
	pendingAlerts map[player.PlayerID][]statemessage.Alert
	// PointGoal     float32 `json:"pointGoal"`
	// languages     map[string]challenges.Language

	// Interactions
	attachments map[*machines.Machine]playerSet
	routes      map[*player.Player]*node.Route

	// timelimit should be able to set a timelimit and count points at the end
	options gameOptions
	mode    modes.Mode
	aChan   chan nwmessage.Message
	// timer
	jobTimer *timer.Timer
	clock    int
}

// init methods:

func NewModel(options *gameOptions) (*GameModel, error) {
	gm := &GameModel{
		Teams:         make(map[string]*team),
		Players:       make(map[player.PlayerID]*player.Player),
		aChan:         make(chan nwmessage.Message, 100), // why is this buffered? TODO
		pendingAlerts: make(map[player.PlayerID][]statemessage.Alert),
		attachments:   make(map[*machines.Machine]playerSet),
		routes:        make(map[*player.Player]*node.Route),
		jobTimer:      timer.NewTimer().Start(),
	}

	gm.jobTimer.AddScheduledJob("score", gm.scoreTick, 1*time.Second)
	gm.jobTimer.AddScheduledJob("clock", gm.gameClock, 1*time.Second)

	if options == nil {
		gm.options = newDefaultOptions()
	} else {
		gm.options = *options
	}

	// initialize
	err := gm.init()
	if err != nil {
		return nil, err
	}

	// return
	return gm, nil
}

func (gm *GameModel) init() error {
	// generate map
	var err error
	gm.Map, err = gm.options.mapGen(gm.options.mapSize)
	if err != nil {
		return err
	}

	// add teams TODO (should be dynamic based on options)
	err = gm.addTeams(makeDummyTeams())
	if err != nil {
		return err
	}

	gm.setDefaultLanguage()

	return nil
}

func makeDummyTeams() []*team {
	teams := make([]*team, 2)
	teams[0] = NewTeam("red")
	teams[1] = NewTeam("blue")
	return teams
}

// GameModel methods --------------------------------------------------------------------------

// fulfill Room interface
func (gm *GameModel) GetPlayers() []*player.Player {
	list := make([]*player.Player, len(gm.Players))
	var i int
	for _, p := range gm.Players {
		list[i] = p
		i++
	}
	return list
}

// func (gm *GameModel) Recv(msg nwmessage.ClientMessage) error {
// 	// return gameCommands.Exec(gm, msg)
// }

func (gm *GameModel) Name() string {
	// gm.aChan <- msg
	return gm.name
}

func (gm *GameModel) Type() room.Type {
	// gm.aChan <- msg
	return room.Game
}

// Addteams can only be called once. After addteams is called unused poes are removed
func (gm *GameModel) addTeams(teams []*team) error {
	nodes := gm.Map.CollectEmptyPoes()

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

		poeNode.Feature.DummyClaim(t.Name, "MIN")
		err := t.addPoe(poeNode)
		if err != nil {
			panic(err)
		}
		gm.calcPoweredNodes(t)

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

func (gm *GameModel) setDefaultLanguage() {
	// tries to set initial language to one of these defaults before picking first available
	var langDefaults = []string{
		"python",
		"javascript",
		"golang",
		"c++",
	}

	// try to use a common language as default
	for _, langName := range langDefaults {
		if _, ok := gm.options.languages[langName]; ok {
			gm.options.defaultLang = langName
			return
		}
	}

	// assign semi randomly
	if gm.options.defaultLang == "" {
		for langName := range gm.options.languages {
			gm.options.defaultLang = langName
			break
		}
	}
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
	fmt.Println("<updateCoinPerTick>")

	// for each node in t.powered add power for each module
	// if a slot is producing, set slot.Processing = true so we can animate this

	// reset
	t.coinPerTick = 0
	// go through each node
	for _, node := range t.powered {
		fmt.Printf("checking node: %d\n", node.ID)
		// store the powerpermod of that node
		t.coinPerTick += node.CoinProduction(t.Name)
	}
}

func (gm *GameModel) calcPoweredNodes(t *team) {
	fmt.Println("<calcPoweredNodes>")
	t.powered = nil // clear previous list of powered nodes
	for _, n := range gm.Map.Nodes {

		if n.HasMachineFor(t.Name) {
			var foundPower bool
			for poe := range t.poes {
				if gm.Map.RouteToNode(t.Name, n, poe) != nil {
					foundPower = true
					break
				}
			}
			// fmt.Printf("foundPower: %t\n", foundPower)

			if foundPower {
				n.PowerMachines(t.Name, true)
				t.powered = append(t.powered, n)
			} else {
				n.PowerMachines(t.Name, false)
			}
		}
	}
}

func (gm *GameModel) playersAt(n *node.Node) []*player.Player {
	players := make([]*player.Player, 0)
	for _, p := range gm.Players {
		if gm.PlayerLocation(p) == n {
			players = append(players, p)
		}
	}
	return players
}

// func (gm *GameModel) detachOtherPlayers(p *player.Player, msg string) {
// 	if gm.CurrentMachine(p) == nil {
// 		log.Panic("player.Player is not attached to a machine")
// 	}

// 	for _, player := range gm.playersAt(gm.routes[p].Endpoint()) {
// 		if player != p {
// 			comment := gm.options.languages[player.Language()].CommentPrefix

// 			editMsg := fmt.Sprintf("%s %s", comment, msg)

// 			player.Outgoing(nwmessage.PsAlert(fmt.Sprintf("You have been detached from the machine at %s", gm.CurrentMachine(p).address)))

// 			player.Outgoing(nwmessage.EditState(editMsg))

// 			gm.BreakConnection(player, true)

// 		}
// 	}
// }

func (gm *GameModel) tryClaimMachine(p *player.Player, node *node.Node, mac *machines.Machine, solution challenges.Solution, fType feature.Type) {

	var hostile bool
	var friendly bool

	if !mac.IsNeutral() {
		if !mac.BelongsTo(p.TeamName) {
			hostile = true
		} else {
			friendly = true
		}
	}

	if hostile {
		switch {
		case mac.Health() == mac.MaxHealth:
			p.Outgoing(nwmessage.PsError(fmt.Errorf("Current solution of %d/%d is the best possible so you cannot steal this machine,\nuse 'reset' instead of 'make' to remove opponent's solution.", mac.Health(), mac.MaxHealth)))

			return

		case solution.Strength < mac.Health():
			p.Outgoing(nwmessage.PsError(fmt.Errorf("Solution (%d/%d) too weak to install, need at least %d/%d to steal", solution.Strength, mac.MaxHealth, mac.Health()+1, mac.MaxHealth)))

			return

		case solution.Strength == mac.Health():
			p.Outgoing(nwmessage.PsAlert(fmt.Sprintf("You need to pass one more test to steal,\nbut your %d/%d is enough to reset this machine.\nKeep trying if you think you can do\nbetter or type 'reset' to proceed", solution.Strength, mac.MaxHealth)))

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
	allowed := node.HasMachineFor(newTeam.Name)
	if hostile {
		oldAllowed = node.HasMachineFor(oldTeam.Name)
	}

	// TODO I think we need to do same for old team, but test first

	// refactor module to new owner and health
	mac.Claim(p.TeamName, solution)

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
	if node.HasMachineFor(newTeam.Name) != allowed {
		gm.calcPoweredNodes(newTeam)
	}
	if hostile {
		if node.HasMachineFor(oldTeam.Name) != oldAllowed {
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
		gm.psBroadcastExcept(p, nwmessage.PsAlert(fmt.Sprintf("%s of (%s) stole a (%s) machine in node %d", p.Name(), p.TeamName, oldTeam.Name, node.ID)))
		p.Outgoing(nwmessage.PsSuccess(fmt.Sprintf("You stole (%v)'s machine, new machine health: %d/%d", oldTeam.Name, mac.Health(), mac.MaxHealth)))

	} else if friendly {
		gm.psBroadcastExcept(p, nwmessage.PsAlert(fmt.Sprintf("%s of (%s) refactored a friendly machine in node %d", p.Name(), p.TeamName, node.ID)))
		p.Outgoing(nwmessage.PsSuccess(fmt.Sprintf("Refactored friendly machine to %d/%d [%s]", solution.Strength, mac.MaxHealth, mac.Solution.Language)))

	} else {
		gm.psBroadcastExcept(p, nwmessage.PsAlert(fmt.Sprintf("%s of (%s) constructed a machine in node %d", p.Name(), p.TeamName, node.ID)))
		p.Outgoing(nwmessage.PsSuccess(fmt.Sprintf("Solution installed in [%s], Health: %d/%d", solution.Language, mac.Health(), mac.MaxHealth)))
	}

	if node.DominatedBy(p.TeamName) {
		gm.psBroadcastExcept(p, nwmessage.PsAlert(fmt.Sprintf("%s is now dominating node %d (production bonus)", p.TeamName, node.ID)))
		p.Outgoing(nwmessage.PsSuccess(fmt.Sprintf("Your team now dominates node %d! (production bonus)", node.ID)))
	}
}

func (gm *GameModel) tryResetMachine(p *player.Player, node *node.Node, mac *machines.Machine, solution challenges.Solution) {

	if mac.IsNeutral() {
		p.Outgoing(nwmessage.PsError(errors.New("Machine is already neutral")))

		return
	}

	if solution.Strength < mac.Health() {
		p.Outgoing(nwmessage.PsError(fmt.Errorf("Solution too weak: %d/%d, need %d/%d to remove", solution.Strength, mac.MaxHealth, mac.Health(), mac.MaxHealth)))
		return
	}

	// track old owner to evaluate traffic after module loss
	oldTeam := gm.Teams[mac.TeamName]

	// track whether node allowed routing for active player before refactor
	allowed := node.HasMachineFor(oldTeam.Name)

	// reset the machine
	mac.Reset()

	// evaluate routing of player trffic through node
	gm.evalTrafficForTeam(node, oldTeam)

	// if routing status has changed, recalculate powered nodes
	if node.HasMachineFor(oldTeam.Name) != allowed {
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
	gm.psBroadcastExcept(p, nwmessage.PsAlert(fmt.Sprintf("%s of (%s) reset a (%s) machine in node %d", p.Name(), p.TeamName, oldTeam.Name, node.ID)))
	p.Outgoing(nwmessage.PsSuccess("Machine reset"))

}

func (gm *GameModel) startGame(when int) error {
	if gm.mode == modes.Running {
		return errors.New("Game already running")
	}

	if gm.mode == modes.Over {
		return errors.New("Game already over")
	}

	if when == 0 {
		gm.mode = modes.Running
	} else {
		gm.clock = when * -1
		gm.mode = modes.Countdown
	}

	return nil
}

func (gm *GameModel) stopGame() error {
	if gm.mode != modes.Running {
		return errors.New("Game is not running")
	}

	gm.mode = modes.Over
	return nil
}

func (gm *GameModel) resetMap(m *node.Map) {
	gm.Map = m

	// Tell everyone to clear their maps
	for _, p := range gm.Players {
		p.Outgoing(nwmessage.GraphReset())

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
		p.Outgoing(nwmessage.ScoreState(gm.packScores()))

	}
}

func (gm *GameModel) makeRouteMap() *statemessage.TrafficMap {
	// collect node.Routes, TODO redundat loop
	traffic := statemessage.NewTrafficMap()

	for _, p := range gm.Players {
		if gm.routes[p] != nil {
			traffic.AddRoute(gm.routes[p], p.TeamName)
		}
	}

	return traffic
}

func (gm *GameModel) broadcastState() {
	for _, p := range gm.Players {
		// TODO feels super hackey to have to pass in node.Routemap but was quickest solution for now.
		state := gm.calcState(p)
		p.Outgoing(nwmessage.GraphState(state))

	}
}

func (gm *GameModel) broadcastGraphReset() {
	for _, p := range gm.Players {
		p.Outgoing(nwmessage.GraphReset())

	}
}

func (gm *GameModel) packScores() string {
	teamSlice := make([]team, len(gm.Teams))
	i := 0
	for _, t := range gm.Teams {
		teamSlice[i] = *t
		i++
	}

	scoreMsg, err := json.Marshal(teamSlice)
	if err != nil {
		log.Println(err)
	}
	return string(scoreMsg)
}

// calcState takes a player argument on the assumption that at some point we'll want to show different states to different players
func (gm *GameModel) calcState(p *player.Player) string {
	routeMap := gm.makeRouteMap()

	// calculate player location
	var playerLoc node.NodeID
	if gm.routes[p] == nil {
		playerLoc = -1
	} else {
		playerLoc = gm.routes[p].Endpoint().ID
	}

	// copy map so we can doctor for this player
	// thisMap := new(node.Map)
	// thisMap = *gm.Map

	// compose state message
	state := statemessage.Message{
		Map:        gm.Map,
		Alerts:     gm.pendingAlerts[p.ID],
		PlayerLoc:  playerLoc,
		TrafficMap: routeMap,
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

func (gm *GameModel) pushActionAlert(color string, location node.NodeID) {
	for k := range gm.pendingAlerts {
		gm.pendingAlerts[k] = append(gm.pendingAlerts[k], statemessage.Alert{color, location})
	}
}

// send a pseudoServer message to all players
func (gm *GameModel) psBroadcast(msg nwmessage.Message) {
	msg.Sender = "pseudoServer"

	for _, player := range gm.Players {
		player.Outgoing(msg)

		player.Outgoing(nwmessage.PsPrompt(player.Prompt()))

	}
}

// broadcast to all but one player
func (gm *GameModel) psBroadcastExcept(p *player.Player, msg nwmessage.Message) {
	msg.Sender = "pseudoServer"

	for _, player := range gm.Players {
		//skip if it's our player
		if player == p {
			continue
		}
		player.Outgoing(msg)

		player.Outgoing(nwmessage.PsPrompt(p.Prompt()))

	}
}

func (gm *GameModel) setPlayerName(p *player.Player, n string) error {

	// check to see if name is in use
	for _, player := range gm.Players {
		if player.Name() == n {
			return errors.New("Name '" + n + "' already in use")
		}
	}

	p.SetName(n)
	return nil
}

// AddPlayer ...
func (gm *GameModel) AddPlayer(p *player.Player) error {
	if _, ok := gm.Players[p.ID]; ok {
		return errors.New("player '" + p.Name() + "' is already in this game")
	}

	// p.inGame = true
	gm.Players[p.ID] = p
	gm.pendingAlerts[p.ID] = make([]statemessage.Alert, 0) // make alerts slot for new player

	supportedLangs := make([]string, len(gm.options.languages))
	var i int
	for lang := range gm.options.languages {
		supportedLangs[i] = lang
		i++
	}

	p.Outgoing(nwmessage.LangSupportState(supportedLangs))

	gm.setLanguage(p, gm.options.defaultLang)

	// send initiall map state

	p.Outgoing(nwmessage.GraphReset())

	p.Outgoing(nwmessage.GraphState(gm.calcState(p)))

	// TODO remove for production
	// this is only here to honestly represent gamestate
	// until it can be handled more gracefully
	welcomeStr := "Game is in sandbox mode. Use 'begin' to start keeping score"
	if gm.mode == modes.Running {
		welcomeStr = "Game has started, get in there!"
	}
	p.Outgoing(nwmessage.PsNeutral(welcomeStr))

	// send initial prompt state
	p.Outgoing(nwmessage.PsPrompt(p.Prompt()))
	return nil
}

// RemovePlayer ...
func (gm *GameModel) RemovePlayer(p *player.Player) error {
	fmt.Printf("<gm.RemovePlayer> Removing player, %s\n", p.Name())
	if _, ok := gm.Players[p.ID]; !ok {
		return errors.New("player '" + p.Name() + "' is not registered")
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
	gm.BreakConnection(p, false)
	// p.inGame = false

	p.Outgoing(nwmessage.LangSupportState([]string{}))

	p.Outgoing(nwmessage.GraphReset())

	return nil
}

func (gm *GameModel) assignPlayerToTeam(p *player.Player, tn teamName) error {
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

func (gm *GameModel) tryConnectPlayerToNode(p *player.Player, n node.NodeID) (*node.Route, error) {

	// log.Printf("source: %v, poeOK: %v, gm.POEs: %v", source, poeOK, gm.POEs)
	team := gm.Teams[p.TeamName]

	if len(team.poes) < 1 {
		return nil, errors.New("No point of entry")
	}

	target := gm.Map.GetNode(n)
	if target == nil {
		return nil, fmt.Errorf("Invalid node, '%d'", n)
	}

	for source := range team.poes {
		// log.Printf("player %v attempting to connect to node %v from POE %v", p.Name(), n, source.ID)
		route := gm.Map.RouteToNode(p.TeamName, source, target)
		// log.Printf("route: %+v\n", node.route)
		if route != nil {
			// log.Println("Successful Connect")
			// log.Printf("Route to target: %v", route)
			gm.BreakConnection(p, false)
			err := gm.establishConnection(p, *route)

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
func (gm *GameModel) establishConnection(p *player.Player, r node.Route) error {
	// set's players node.Route to the node.Route generated via node.RouteToNode
	// gm.Routes[p.ID] = &node.Route{Endpoint: n, Nodes: routeNodes}
	// make sure we're not blocked by any firewalls:
	for _, n := range r.Nodes {
		if n.Feature.Type == feature.Firewall {
			if n.MachinesFor(p.TeamName) <= gm.trafficCount(n, p.TeamName) {
				return fmt.Errorf("Connection refused (firewall at node %d)", n.ID)
			}
		}
	}

	gm.routes[p] = &r
	return nil
	// return gm.Routes[p.ID]
}

func (gm *GameModel) trafficCount(n *node.Node, t teamName) int {
	var count int

	for p := range gm.Teams[t].players {
		if gm.routes[p] != nil {
			if gm.routes[p].RunsThrough(n) || gm.routes[p].Endpoint() == n {
				count++
			}
		}
	}
	return count
}

func (gm *GameModel) evalTrafficForTeam(n *node.Node, t *team) {
	// if the module no longer supports routing for this modules team
	if !n.HasMachineFor(t.Name) {
		for _, p := range gm.Players {
			// check each p who is on team's node.Route
			if p.TeamName == t.Name {
				// and if it contained that node, break the players connection
				if gm.routes[p] != nil {
					if gm.routes[p].RunsThrough(n) {
						gm.BreakConnection(p, true)
					}
				}
			}
		}
	}
}

func (gm *GameModel) setLanguage(p *player.Player, l string) error {
	_, ok := gm.options.languages[l]

	if !ok {
		return fmt.Errorf("'%v' is not a supported in this match. Use 'langs' to list available languages")
	}

	mac := gm.CurrentMachine(p)
	if mac != nil && !mac.AcceptsLanguageFrom(p, l) {
		return errors.New("Can't change language while attached to a hostile machine")
	}

	// TODO combine
	p.SetLanguage(l)
	p.Outgoing(nwmessage.EditLangState(l))

	return nil
}

// Machine related methods
func (gm *GameModel) attachPlayer(p *player.Player, m *machines.Machine) {
	if _, ok := gm.attachments[m]; !ok {
		gm.attachments[m] = make(playerSet)
	}
	gm.attachments[m][p] = true
}

func (gm *GameModel) detachPlayer(p *player.Player) {
	m := gm.CurrentMachine(p)
	if m != nil {
		delete(gm.attachments[m], p)
	}
}

func (gm *GameModel) detachAll(m *machines.Machine, msg string) {
	for p := range gm.attachments[m] {
		gm.detachPlayer(p)
	}

}

func (gm *GameModel) resetMachine(m *machines.Machine) {
	msg := fmt.Sprintf("mac:%s is resetting, you have been detached", m.Address)
	gm.detachAll(m, msg)
	m.Reset()

}

func node2Str(n *node.Node, p *player.Player) string {

	// sort keys for consistent presentation
	addList := make([]string, 0)
	for add := range n.Addresses() {
		addList = append(addList, add)
	}
	sort.Strings(addList)

	// compose list of all machines
	macList := ""
	for _, add := range addList {
		atIndicator := ""
		if p.MacAddress() == add {
			atIndicator = "*"
		}
		mac := n.MacAt(add)
		macList += "\n" + add + ":" + mac2Str(mac, p) + atIndicator
	}

	connectList := strings.Trim(strings.Join(strings.Split(fmt.Sprint(n.Connections), " "), ","), "[]")

	return fmt.Sprintf("NodeID: %v\nConnects To: %s\nMachines: %v", n.ID, connectList, macList)
}

func mac2Str(m *machines.Machine, p *player.Player) string {
	// p is the player observing the machine. Unused atm but allows us to show/hide info from different players
	var gateway string

	if m.IsGateway() {
		gateway = " (gateway)"
	}

	details := fmt.Sprintf("[%s] [%s] [%s] [%d/%d]", m.TeamName, m.Solution.Author, m.Solution.Language, m.Health(), m.MaxHealth)

	switch {
	case m.TeamName != "":
		return "(" + details + ")" + gateway
	default:
		return "( -neutral- )" + gateway
	}
}
