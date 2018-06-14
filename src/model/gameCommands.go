package model

import (
	"argument"
	"challenges"
	"command"
	"errors"
	"feature"
	"fmt"
	"model/node"
	"model/player"
	"nwmessage"
	"room"
	"sort"
	"strings"
)

type gameCommand struct {
	command.Info
	handler func(p *player.Player, gm *GameModel, args []interface{}) error
}

func (c gameCommand) Exec(cli nwmessage.Client, context room.Room, args []interface{}) error {
	p, ok := cli.(*player.Player)
	if !ok {
		panic("Error asserting Player in exec")
	}

	gm, ok := context.(*GameModel)
	if !ok {
		panic("Error asserting GameModel in exec")
	}

	return c.handler(p, gm, args)
}

func RegisterCommands(r *command.Registry) {
	gameCommands := []gameCommand{
		{
			command.Info{
				CmdName:   "begin",
				ShortDesc: "Begins the game. Immediately or in n seconds",
				ArgsReq:   argument.ArgList{},
				ArgsOpt: argument.ArgList{
					{Name: "n_seconds", Type: argument.Int},
				},
				CmdContexts: []room.Type{room.Game},
			},
			cmdStartGame,
		},

		{
			command.Info{
				CmdName:   "tc",
				ShortDesc: "Sends a message to all teammates",
				ArgsReq: argument.ArgList{
					{Name: "msg", Type: argument.GreedyString},
				},
				ArgsOpt:     argument.ArgList{},
				CmdContexts: []room.Type{room.Game},
			},
			cmdTeamChat,
		},

		// "say": {
		// 	CmdName:   "say",
		// 	ShortDesc: "Sends all players at the same node",
		// 	ArgsReq: argument.ArgList{
		// 		{Name: "msg", Type: argument.GreedyString},
		// 	},
		// 	ArgsOpt: argument.ArgList{},
		// 	Handler: cmdSay,
		// },

		{
			command.Info{
				CmdName:   "join",
				ShortDesc: "Joins a team",
				LongDesc:  "Joins either a specified team is one is provided or, if no argument is given, a team is selected automatically",
				ArgsReq:   argument.ArgList{},
				ArgsOpt: argument.ArgList{
					{Name: "team_name", Type: argument.String},
				},
				CmdContexts: []room.Type{room.Game},
			},
			cmdJoinTeam,
		},

		{
			command.Info{
				CmdName:   "con",
				ShortDesc: "Connect to the specified node",
				ArgsReq: argument.ArgList{
					{Name: "node_id", Type: argument.Int},
				},
				ArgsOpt:     argument.ArgList{},
				CmdContexts: []room.Type{room.Game},
			},
			cmdConnect,
		},

		{
			command.Info{
				CmdName:   "foc",
				ShortDesc: "Controls map focus",
				LongDesc:  "Focuses on the specified node or resets focus to include all nodes",
				ArgsReq:   argument.ArgList{},
				ArgsOpt: argument.ArgList{
					{Name: "node_id", Type: argument.Int},
				},
				CmdContexts: []room.Type{room.Game},
			},
			cmdGraphFocus,
		},

		{
			command.Info{
				CmdName:   "lang",
				ShortDesc: "Select a programming language",
				ArgsReq: argument.ArgList{
					{Name: "lang_name", Type: argument.String},
				},
				ArgsOpt:     argument.ArgList{},
				CmdContexts: []room.Type{room.Game},
			},
			cmdLang,
		},

		{
			command.Info{
				CmdName:     "langs",
				ShortDesc:   "List languages allowed in this game",
				ArgsReq:     argument.ArgList{},
				ArgsOpt:     argument.ArgList{},
				CmdContexts: []room.Type{room.Game},
			},
			cmdListLanguages,
		},

		{
			command.Info{
				CmdName:   "make",
				ShortDesc: "Submits code to claim or steal machine",
				LongDesc:  "An argument must be provided if the current machine is a Feature. The argument is the type of Feature player would like to install",
				ArgsReq:   argument.ArgList{},
				ArgsOpt: argument.ArgList{
					{Name: "feature", Type: argument.String},
				},
				CmdContexts: []room.Type{room.Game},
			},
			cmdMake,
		},

		{
			command.Info{
				CmdName:     "test",
				ShortDesc:   "Runs code using custom stdin instead of challenge",
				ArgsReq:     argument.ArgList{},
				ArgsOpt:     argument.ArgList{},
				CmdContexts: []room.Type{room.Game},
			},
			cmdTestCode,
		},

		{
			command.Info{
				CmdName:     "reset",
				ShortDesc:   "Submist code to reset machine",
				ArgsReq:     argument.ArgList{},
				ArgsOpt:     argument.ArgList{},
				CmdContexts: []room.Type{room.Game},
			},
			cmdResetMachine,
		},

		{
			command.Info{
				CmdName:   "at",
				ShortDesc: "Attach to a machine in the current node",
				LongDesc:  "An optional second argument can be provided, if that argument is 'n' or 'no', then no boilerplate will be loaded, thus preserving the current editor state.",
				ArgsReq: argument.ArgList{
					{Name: "addr", Type: argument.String},
				},
				ArgsOpt: argument.ArgList{
					{Name: "bp_flag", Type: argument.String},
				},
				CmdContexts: []room.Type{room.Game},
			},
			cmdAttach,
		},

		{
			command.Info{
				CmdName:     "who",
				ShortDesc:   "Prints who's in the current game",
				ArgsReq:     argument.ArgList{},
				ArgsOpt:     argument.ArgList{},
				CmdContexts: []room.Type{room.Game},
			},
			cmdWho,
		},

		{
			command.Info{
				CmdName:     "ls",
				ShortDesc:   "Prints details about current node",
				ArgsReq:     argument.ArgList{},
				ArgsOpt:     argument.ArgList{},
				CmdContexts: []room.Type{room.Game},
			},
			cmdLs,
		},
	}

	for _, c := range gameCommands {
		err := r.AddEntry(c)
		if err != nil {
			panic(err)
		}
	}
}

func cmdStartGame(p *player.Player, gm *GameModel, args []interface{}) error {

	var err error

	if len(args) > 0 {
		count := args[0].(int)
		err = gm.startGame(count)
		gm.psBroadcast(nwmessage.PsAlert(fmt.Sprintf("Game will start in %d seconds!\n", count)))
	} else {
		gm.psBroadcast(nwmessage.PsAlert(fmt.Sprintf("Game has started!")))
		err = gm.startGame(0)
	}

	if err != nil {
		return err
	}

	return nil
}

func cmdTell(p *player.Player, gm *GameModel, args []interface{}) error {

	recipName := args[0].(string)
	msgText := args[1].(string)

	var recip *player.Player
	for _, player := range gm.Players {
		if player.Name() == recipName {
			recip = player
		}
	}

	if recip == nil {
		return fmt.Errorf("No such player, '%s'", recipName)
	}

	chatMsg := fmt.Sprintf("%s > %s", p.Name(), msgText)

	recip.Outgoing(nwmessage.PsChat(p.Name(), chatMsg, "(private)"))
	return nil
}

func cmdTeamChat(p *player.Player, gm *GameModel, args []interface{}) error {
	chatMsg := args[0].(string)

	if p.TeamName == "" {
		return nwmessage.ErrorNoTeam()
	}

	gm.Teams[p.TeamName].broadcast(nwmessage.PsChat(p.Name(), "team", chatMsg))

	return nil

}

func cmdSay(p *player.Player, gm *GameModel, args []interface{}) error {
	chatMsg := args[0].(string)

	node := gm.PlayerLocation(p)
	if node == nil {
		return errors.New("Can only 'say' while connected to a node")
	}

	msg := nwmessage.PsChat(p.Name(), "node", chatMsg)

	for _, p := range gm.playersAt(node) {
		p.Outgoing(msg)

	}

	return nil
}

func cmdJoinTeam(p *player.Player, gm *GameModel, args []interface{}) error {
	var t teamName

	if len(args) == 0 || args[0] == "" {
		tt := gm.trailingTeam()
		p.Outgoing(nwmessage.PsNeutral(fmt.Sprintf("Auto-assigning to team '%s'", tt)))
		t = teamName(tt)
	} else {
		t = args[0].(teamName)
	}

	err := gm.assignPlayerToTeam(p, t)
	if err != nil {
		return err
	}

	p.Outgoing(nwmessage.PsSuccess(fmt.Sprintf("You've joined team '%s'", t)))

	if len(gm.Teams[p.TeamName].poes) < 1 { // if team's got no peos
		p.Outgoing(nwmessage.PsNeutral("Your team doesn't have a point of entry yet.\nUse 'sp node_id' to set one and begin playing"))
	} else {
		// if we do have poes, connect player to a randome one
		var tp *node.Node
		for n := range gm.Teams[p.TeamName].poes {
			tp = n
			break
		}

		if tp != nil {
			p.Outgoing(nwmessage.PsNeutral(fmt.Sprintf("Team's point of entry is node %d.\nConnecting you there now...\n", tp.ID)))
			// log.Printf("player joined team, trying to log into %v", tp.ID)
			var i []interface{} = make([]interface{}, 1)
			i[0] = tp.ID
			err = cmdConnect(p, gm, i)
			if err != nil {
				panic(err)
			}
			gm.broadcastState()
		}

	}

	p.Outgoing(nwmessage.TeamState(p.TeamName))
	return nil
}

func cmdLang(p *player.Player, gm *GameModel, args []interface{}) error {
	lang := args[0].(string)

	err := gm.setLanguage(p, lang)
	if err != nil {
		return err
	}

	// if the player's attached somewhere, update the buffer
	mac := gm.CurrentMachine(p)
	if mac != nil {
		// TODO syntax
		langDetails := gm.options.languages[p.Language()]
		boilerplate := langDetails.Boilerplate
		comment := langDetails.CommentPrefix
		sampleIO := mac.Challenge.SampleIO
		description := mac.Challenge.ShortDesc

		editText := fmt.Sprintf("%s Challenge:\n%s %s\n%s Sample IO: %s\n\n%s", comment, comment, description, comment, sampleIO, boilerplate)
		p.Outgoing(nwmessage.EditState(editText))

	}

	p.Outgoing(nwmessage.PsSuccess(fmt.Sprintf("Language set to %s", lang)))
	return nil
}

func cmdListLanguages(p *player.Player, gm *GameModel, args []interface{}) error {
	var langs sort.StringSlice

	for l := range gm.options.languages {
		langs = append(langs, l)
	}

	langs.Sort()

	p.Outgoing(nwmessage.PsNeutral("This game supports:\n" + strings.Join(langs, "\n")))
	return nil
}

func cmdConnect(p *player.Player, gm *GameModel, args []interface{}) error {
	target := args[0].(int)
	if p.TeamName == "" {
		return errors.New("Join a team first")
	}

	// break any pre-existing connection before connecting elsewhere
	_, err := gm.tryConnectPlayerToNode(p, target)
	if err != nil {
		return err
	}

	gm.broadcastState()
	p.Outgoing(nwmessage.PsSuccess(fmt.Sprintf("Connection established:")))

	return cmdLs(p, gm, args)
}

func cmdWho(p *player.Player, gm *GameModel, args []interface{}) error {

	// Sort team names
	whoStr := ""
	teamNames := make([]string, 0)
	for tname := range gm.Teams {
		teamNames = append(teamNames, tname)
	}
	sort.Strings(teamNames)

	// build rosters and display string
	for _, n := range teamNames {
		t := gm.Teams[n]
		whoStr += n + ":\n"
		for mem := range t.players {
			whoStr += "\t" + mem.Name() + "\n"
		}
	}

	//TODO list unassigned players? leave as invisible for observation? list just a number maybe

	p.Outgoing(nwmessage.PsNeutral(whoStr))
	return nil
}

func cmdLs(p *player.Player, gm *GameModel, args []interface{}) error {

	if gm.routes[p] == nil {
		return nwmessage.ErrorNoConnection()
	}

	node := gm.routes[p].Endpoint()
	retMsg := node2Str(node, p)
	pHere := gm.playersAt(node)

	if len(pHere) > 1 {
		//make slice of names (excluding this player)
		names := make([]string, 0, len(pHere)-1)
		for _, player := range pHere {
			if player.Name() != p.Name() {
				names = append(names, fmt.Sprintf("%s (%s)", player.Name(), player.TeamName))
			}
		}

		//join the slice to string
		addMsg := "\nAlso here: " + strings.Join(names, ", ")

		//add to message
		retMsg += addMsg
	}

	p.Outgoing(nwmessage.PsNeutral(retMsg))
	return nil
}

// func cmdSetPOE(p *player.Player, gm *GameModel, args []interface{}) error {
// p := cl.(*player.Player)
// gm := context.(*GameModel)
// 	if p.TeamName == "" {
// 		p.Outgoing(nwmessage.PsNoTeam())
// return nil
// 	}

// 	nodeId, err := strconv.Atoi(args[0])
// 	if err != nil {
// 		return fmt.Errorf("expected integer, got '%v'", args[0])
// 	}

// 	newPoe, err = gm.Map.getNode(nodeId)
// 	if err != nil {
// 		return err
// 	}

// 	err = gm.Teams[p.TeamName].addPoe(newPOE)
// 		if err != nil {
// 		return err
// 	}

// 	for player := range gm.Teams[p.TeamName].players {
// 		_, _ = gm.tryConnectPlayerToNode(player, newPOE)
// 	}

// 	// if all teams have their poe set
// 	var ready int
// 	for _, team := range gm.Teams {
// 		if len(team.poes) > 0 {
// 			ready++
// 		}
// 	}

// 	// start the game
// 	if len(gm.Teams) == ready {
// 		gm.startGame()
// 	}

// 	gm.broadcastState()
// 	p.Outgoing(nwmessage.PsSuccess(fmt.Sprintf("%s team's point of entry set to node %d\nConnecting you there now...", p.TeamName, newPOE)))
// return nil
// }

func cmdTestCode(p *player.Player, gm *GameModel, args []interface{}) error {

	if p.Editor() == "" {
		return errors.New("No code submitted")
	}

	go func() {
		defer p.SendPrompt()

		response := challenges.GetOutput(p.Language(), p.Editor(), p.Stdin())
		p.Outgoing(nwmessage.PsSuccess("Finished running (check output box)"))

		p.Outgoing(nwmessage.ResultState(response))

	}()

	p.Outgoing(nwmessage.PsBegin(fmt.Sprintf("Testing code...")))
	return nil
}

func cmdAttach(p *player.Player, gm *GameModel, args []interface{}) error {
	macAddress := args[0].(string)

	if gm.routes[p] == nil {
		return nwmessage.ErrorNoConnection()
	}

	err := gm.routes[p].Endpoint().CanAttach(p.TeamName, macAddress)

	if err != nil {
		return err
	}

	// passed checks, set player mac to target

	// remove old attachments
	gm.detachPlayer(p)

	// add this attachment
	p.SetMacAddress(macAddress)
	mac := gm.CurrentMachine(p)
	gm.attachPlayer(p, mac)

	// if the mac has an enemy module, player's language is set to that module's
	langLock := false
	if !mac.AcceptsLanguageFrom(p, p.Language()) {
		langLock = true
		gm.setLanguage(p, mac.Solution.Language)
		p.Outgoing(nwmessage.LangSupportState([]string{mac.Solution.Language}))

	} else {
		supportedLangs := make([]string, len(gm.options.languages))
		var i int
		for lang := range gm.options.languages {
			supportedLangs[i] = lang
			i++
		}

		p.Outgoing(nwmessage.LangSupportState(supportedLangs))

	}

	// get language details
	langDetails := gm.options.languages[p.Language()]
	boilerplate := langDetails.Boilerplate
	comment := langDetails.CommentPrefix

	p.SetChallenge(mac.Challenge)
	p.SetStdin(mac.Challenge.SampleIO[0].Input, true)

	var lockStr string
	if langLock {
		lockStr = fmt.Sprintf("\n\n%sHOSTILE MACHINE, SOLUTION MUST BE IN [%s]", comment, strings.ToUpper(gm.CurrentMachine(p).Solution.Language))
	}

	var flag string
	if len(args) > 1 {
		flag = args[1].(string)
	}

	if flag != "n" && flag != "no" {
		editText := boilerplate + lockStr
		p.SetEditor(editText, true)
	}

	retText := fmt.Sprintf("Attached to machine at %s: \ncontents:%v", macAddress, mac2Str(mac, p))
	retText += lockStr

	p.Outgoing(nwmessage.PsSuccess(retText))
	return nil
}

func cmdMake(p *player.Player, gm *GameModel, args []interface{}) error {

	err := gm.CanSubmit(p)
	if err != nil {
		return err
	}

	mac := gm.CurrentMachine(p)
	node := gm.PlayerLocation(p)

	// Abstract this TODO
	var feaType feature.Type
	if mac.IsGateway() {
		if mac.Type == feature.None {
			if len(args) < 1 {
				return errors.New("Must provide feature argument when attached to gateway, see 'help features'")
			}

			var err error
			feaType, err = feature.FromString(args[0].(string))

			if err != nil {
				return fmt.Errorf("Invalid feature type, '%s'", args[0].(string))
			}
		} else {
			if len(args) > 0 {
				p.Outgoing(nwmessage.PsError(errors.New("Ignoring argument! Cannot change the type on an installed feature")))

			}
		}
	}

	// passed error checks
	go func() {
		response, err := p.SubmitCode(mac.Challenge.ID)

		if err != nil {
			p.Outgoing(nwmessage.PsError(err))

			return
		}

		mac.Lock()
		gm.tryClaimMachine(p, node, mac, response, feaType)
		mac.Unlock()

		p.SendPrompt()
	}()

	p.Outgoing(nwmessage.PsBegin("Compiling..."))
	return nil
}

func cmdResetMachine(p *player.Player, gm *GameModel, args []interface{}) error {
	node := gm.PlayerLocation(p)
	mac := gm.CurrentMachine(p)

	err := gm.CanSubmit(p)
	if err != nil {
		return err
	}

	go func() {

		solution, err := p.SubmitCode(mac.Challenge.ID)

		if err != nil {
			p.Outgoing(nwmessage.PsError(err))
			return
		}

		mac.Lock()
		gm.tryResetMachine(p, node, mac, solution)
		mac.Unlock()

		p.SendPrompt()
	}()

	p.Outgoing(nwmessage.PsBegin("Resetting machine..."))
	return nil
}

func cmdGraphFocus(p *player.Player, gm *GameModel, args []interface{}) error {

	if len(args) < 1 {
		// send resetfocus message
		p.Outgoing(nwmessage.PsNeutral("Resetting map focus..."))
		p.Outgoing(nwmessage.GraphFocusReset())
		return nil
	}

	id := args[0].(int)

	if gm.Map.GetNode(id) == nil {
		return fmt.Errorf("Invalid node id, '%d'", id)
	}

	// send focus message
	p.Outgoing(nwmessage.PsNeutral(fmt.Sprintf("Focusing on node, '%d'", id)))
	p.Outgoing(nwmessage.GraphFocus(id))
	return nil
}
