package nwmodel

import (
	"argument"
	"commands"
	"errors"
	"feature"
	"fmt"
	"nwmessage"
	"receiver"
	"sort"
	"strings"
)

var gameCommands = commands.CommandGroup{
	// "chat": {
	// 	Name:      "chat",
	// 	ShortDesc: "Toggles chat mode (all text entered is broadcast)",
	// 	ArgsReq:   argument.ArgList{},
	// 	ArgsOpt:   argument.ArgList{},
	// 	Handler:   cmdToggleChat,
	// },

	"yell": {
		Name:      "yell",
		ShortDesc: "Sends a message to all player (in the same game/lobby)",
		ArgsReq: argument.ArgList{
			{Name: "msg", Type: argument.GreedyString},
		},
		ArgsOpt: argument.ArgList{},
		Handler: cmdYell,
	},

	"tc": {
		Name:      "tc",
		ShortDesc: "Sends a message to all teammates",
		ArgsReq: argument.ArgList{
			{Name: "msg", Type: argument.GreedyString},
		},
		ArgsOpt: argument.ArgList{},
		Handler: cmdYell,
	},

	// "say": {
	// 	Name:      "say",
	// 	ShortDesc: "Sends all players at the same node",
	// 	ArgsReq: argument.ArgList{
	// 		{Name: "msg", Type: argument.GreedyString},
	// 	},
	// 	ArgsOpt: argument.ArgList{},
	// 	Handler: cmdSay,
	// },

	"join": {
		Name:      "join",
		ShortDesc: "Joins a team",
		LongDesc:  "Joins either a specified team is one is provided or, if no argument is given, a team is selected automatically",
		ArgsReq:   argument.ArgList{},
		ArgsOpt: argument.ArgList{
			{Name: "team_name", Type: argument.String},
		},
		Handler: cmdJoinTeam,
	},

	"con": {
		Name:      "con",
		ShortDesc: "Connect to the specified node",
		ArgsReq: argument.ArgList{
			{Name: "node_id", Type: argument.Int},
		},
		ArgsOpt: argument.ArgList{},
		Handler: cmdConnect,
	},

	"lang": {
		Name:      "lang",
		ShortDesc: "Select a programming language",
		ArgsReq: argument.ArgList{
			{Name: "lang_name", Type: argument.String},
		},
		ArgsOpt: argument.ArgList{},
		Handler: cmdLang,
	},

	"langs": {
		Name:      "langs",
		ShortDesc: "List languages allowed in this game",
		ArgsReq:   argument.ArgList{},
		ArgsOpt:   argument.ArgList{},
		Handler:   cmdListLanguages,
	},

	"make": {
		Name:      "make",
		ShortDesc: "Submits code to claim or steal machine",
		LongDesc:  "An argument must be provided if the current machine is a Feature. The argument is the type of Feature player would like to install",
		ArgsReq:   argument.ArgList{},
		ArgsOpt: argument.ArgList{
			{Name: "feature", Type: argument.String},
		},
		Handler: cmdMake,
	},

	"test": {
		Name:      "test",
		ShortDesc: "Runs code using custom stdin instead of challenge",
		ArgsReq:   argument.ArgList{},
		ArgsOpt:   argument.ArgList{},
		Handler:   cmdTestCode,
	},

	"res": {
		Name:      "res",
		ShortDesc: "Submist code to reset machine",
		ArgsReq:   argument.ArgList{},
		ArgsOpt:   argument.ArgList{},
		Handler:   cmdResetMachine,
	},

	"at": {
		Name:      "at",
		ShortDesc: "Attach to a machine in the current node",
		LongDesc:  "An optional second argument can be provided, if that argument is 'n' or 'no', then no boilerplate will be loaded, thus preserving the current editor state.",
		ArgsReq: argument.ArgList{
			{Name: "addr", Type: argument.String},
		},
		ArgsOpt: argument.ArgList{
			{Name: "bp_flag", Type: argument.String},
		},
		Handler: cmdAttach,
	},

	"who": {
		Name:      "who",
		ShortDesc: "Prints who's in the current game",
		ArgsReq:   argument.ArgList{},
		ArgsOpt:   argument.ArgList{},
		Handler:   cmdWho,
	},

	"ls": {
		Name:      "ls",
		ShortDesc: "Prints details about current node",
		ArgsReq:   argument.ArgList{},
		ArgsOpt:   argument.ArgList{},
		Handler:   cmdLs,
	},
}

func cmdYell(cl nwmessage.Client, context receiver.Receiver, args []interface{}) error {
	p := cl.(*Player)
	gm := context.(*GameModel)

	chatMsg := args[0].(string)

	gm.psBroadcast(nwmessage.PsChat(p.GetName(), "global", chatMsg))
	return nil
}

func cmdTell(cl nwmessage.Client, context receiver.Receiver, args []interface{}) error {
	p := cl.(*Player)
	gm := context.(*GameModel)

	recipName := args[0].(string)
	msgText := args[1].(string)

	var recip *Player
	for _, player := range gm.Players {
		if player.GetName() == recipName {
			recip = player
		}
	}

	if recip == nil {
		return fmt.Errorf("No such player, '%s'", recipName)
	}

	chatMsg := fmt.Sprintf("%s > %s", p.GetName(), msgText)

	recip.Outgoing(nwmessage.PsChat(p.name, chatMsg, "(private)"))
	return nil
}

func cmdTeamChat(cl nwmessage.Client, context receiver.Receiver, args []interface{}) error {
	p := cl.(*Player)
	gm := context.(*GameModel)
	chatMsg := args[0].(string)

	if p.TeamName == "" {
		return nwmessage.ErrorNoTeam()
	}

	gm.Teams[p.TeamName].broadcast(nwmessage.PsChat(p.GetName(), "team", chatMsg))

	return nil

}

func cmdSay(cl nwmessage.Client, context receiver.Receiver, args []interface{}) error {
	p := cl.(*Player)
	gm := context.(*GameModel)
	chatMsg := args[0].(string)

	node := p.location()
	if node == nil {
		return errors.New("Can only 'say' while connected to a node")
	}

	msg := nwmessage.PsChat(p.GetName(), "node", chatMsg)

	for _, p := range gm.playersAt(node) {
		p.Outgoing(msg)

	}

	return nil
}

func cmdJoinTeam(cl nwmessage.Client, context receiver.Receiver, args []interface{}) error {
	p := cl.(*Player)
	gm := context.(*GameModel)
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
		var tp *node
		for n := range gm.Teams[p.TeamName].poes {
			tp = n
			break
		}

		if tp != nil {
			p.Outgoing(nwmessage.PsNeutral(fmt.Sprintf("Team's point of entry is node %d.\nConnecting you there now...\n", tp.ID)))
			// log.Printf("player joined team, trying to log into %v", tp.ID)
			var i []interface{} = make([]interface{}, 1)
			i[0] = tp.ID
			err = cmdConnect(cl, context, i)
			if err != nil {
				panic(err)
			}
			gm.broadcastState()
		}

	}

	p.Outgoing(nwmessage.TeamState(p.TeamName))
	return nil
}

func cmdLang(cl nwmessage.Client, context receiver.Receiver, args []interface{}) error {
	p := cl.(*Player)
	gm := context.(*GameModel)
	lang := args[0].(string)

	err := gm.setLanguage(p, lang)
	if err != nil {
		return err
	}

	// if the player's attached somewhere, update the buffer
	mac := p.currentMachine()
	if mac != nil {
		// TODO syntax
		langDetails := gm.languages[p.language]
		boilerplate := langDetails.Boilerplate
		comment := langDetails.CommentPrefix
		sampleIO := mac.challenge.SampleIO
		description := mac.challenge.ShortDesc

		editText := fmt.Sprintf("%s Challenge:\n%s %s\n%s Sample IO: %s\n\n%s", comment, comment, description, comment, sampleIO, boilerplate)
		p.Outgoing(nwmessage.EditState(editText))

	}

	p.Outgoing(nwmessage.PsSuccess(fmt.Sprintf("Language set to %s", lang)))
	return nil
}

func cmdListLanguages(cl nwmessage.Client, context receiver.Receiver, args []interface{}) error {
	p := cl.(*Player)
	gm := context.(*GameModel)
	var langs sort.StringSlice

	for l := range gm.languages {
		langs = append(langs, l)
	}

	langs.Sort()

	p.Outgoing(nwmessage.PsNeutral("This game supports:\n" + strings.Join(langs, "\n")))
	return nil
}

func cmdConnect(cl nwmessage.Client, context receiver.Receiver, args []interface{}) error {
	p := cl.(*Player)
	gm := context.(*GameModel)
	targetNode := args[0].(int)
	if p.TeamName == "" {
		return errors.New("Join a team first")
	}

	// break any pre-existing connection before connecting elsewhere
	_, err := gm.tryConnectPlayerToNode(p, targetNode)
	if err != nil {
		return err
	}

	gm.broadcastState()
	p.Outgoing(nwmessage.PsSuccess(fmt.Sprintf("Connection established:")))

	return cmdLs(p, gm, args)
}

func cmdWho(cl nwmessage.Client, context receiver.Receiver, args []interface{}) error {
	p := cl.(*Player)
	gm := context.(*GameModel)

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
			whoStr += "\t" + mem.GetName() + "\n"
		}
	}

	//TODO list unassigned players? leave as invisible for observation? list just a number maybe

	p.Outgoing(nwmessage.PsNeutral(whoStr))
	return nil
}

func cmdLs(cl nwmessage.Client, context receiver.Receiver, args []interface{}) error {
	p := cl.(*Player)
	gm := context.(*GameModel)

	if p.Route == nil {
		return nwmessage.ErrorNoConnection()
	}

	node := p.Route.Endpoint()
	retMsg := node.StringFor(p)
	pHere := gm.playersAt(node)

	if len(pHere) > 1 {
		//make slice of names (excluding this player)
		names := make([]string, 0, len(pHere)-1)
		for _, player := range pHere {
			if player.GetName() != p.GetName() {
				names = append(names, fmt.Sprintf("%s (%s)", player.GetName(), player.TeamName))
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

// func cmdSetPOE(cl nwmessage.Client, context receiver.Receiver, args []interface{}) error {
// p := cl.(*Player)
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

func cmdTestCode(cl nwmessage.Client, context receiver.Receiver, args []interface{}) error {
	p := cl.(*Player)
	// gm := context.(*GameModel)

	if p.EditorState == "" {
		return errors.New("No code submitted")
	}

	go func() {
		defer p.SendPrompt()

		response := getOutput(p.language, p.EditorState, p.StdinState)
		p.Outgoing(nwmessage.PsSuccess("Finished running (check output box)"))

		p.Outgoing(nwmessage.ResultState(response))

	}()

	p.Outgoing(nwmessage.PsBegin(fmt.Sprintf("Testing code...")))
	return nil
}

func cmdScore(cl nwmessage.Client, context receiver.Receiver, args []interface{}) error {
	p := cl.(*Player)
	gm := context.(*GameModel)
	var scoreStrs sort.StringSlice

	for teamName, team := range gm.Teams {
		scoreStrs = append(scoreStrs, fmt.Sprintf("%s:\nCoinCoin Production: %.2f\nCoinCoin Stockpiled: %.2f/%.0f", teamName, team.coinPerTick, team.VicPoints, gm.PointGoal))
	}

	scoreStrs.Sort()

	p.Outgoing(nwmessage.PsNeutral(strings.Join(scoreStrs, "\n")))
	return nil
}

func cmdAttach(cl nwmessage.Client, context receiver.Receiver, args []interface{}) error {
	p := cl.(*Player)
	gm := context.(*GameModel)
	macAddress := args[0].(string)

	if p.Route == nil {
		return nwmessage.ErrorNoConnection()
	}

	_, addOk := p.Route.Endpoint().addressMap[macAddress]

	if !addOk {
		return fmt.Errorf("Invalid address, '%s'", macAddress)
	}

	// passed checks, set player mac to target

	// remove old attachments
	p.macDetach()

	// add this attachment
	p.macAddress = macAddress
	mac := p.currentMachine()
	mac.addPlayer(p)

	// if the mac has an enemy module, player's language is set to that module's
	langLock := false
	if !mac.isNeutral() && !mac.belongsTo(p.TeamName) {
		langLock = true
		gm.setLanguage(p, mac.language)
		p.Outgoing(nwmessage.LangSupportState([]string{mac.language}))

	} else {
		supportedLangs := make([]string, len(gm.languages))
		var i int
		for lang := range gm.languages {
			supportedLangs[i] = lang
			i++
		}

		p.Outgoing(nwmessage.LangSupportState(supportedLangs))

	}

	// get language details
	langDetails := gm.languages[p.language]
	boilerplate := langDetails.Boilerplate
	comment := langDetails.CommentPrefix

	p.challengeState(mac.challenge)
	p.stdinState(mac.challenge.SampleIO[0].Input)

	var lockStr string
	if langLock {
		lockStr = fmt.Sprintf("\n\n%sHOSTILE MACHINE, SOLUTION MUST BE IN [%s]", comment, strings.ToUpper(p.currentMachine().language))
	}

	var flag string
	if len(args) > 1 {
		flag = args[1].(string)
	}

	if flag != "n" && flag != "no" {
		editText := boilerplate + lockStr
		p.editState(editText)
	}

	retText := fmt.Sprintf("Attached to machine at %s: \ncontents:%v", macAddress, mac.StringFor(p))
	retText += lockStr

	p.Outgoing(nwmessage.PsSuccess(retText))
	return nil
}

func cmdMake(cl nwmessage.Client, context receiver.Receiver, args []interface{}) error {
	p := cl.(*Player)
	gm := context.(*GameModel)

	err := p.canSubmit()
	if err != nil {
		return err
	}

	mac := p.currentMachine()
	node := p.location()

	// Abstract this TODO
	var feaType feature.Type
	if mac.isFeature() {
		if mac.Type == feature.None {
			if len(args) < 1 {
				return errors.New("Make requires one argument when attached to an untyped feature")
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
		response, err := p.submitCode()

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

func cmdResetMachine(cl nwmessage.Client, context receiver.Receiver, args []interface{}) error {
	p := cl.(*Player)
	gm := context.(*GameModel)
	node := p.location()

	err := p.canSubmit()
	if err != nil {
		return err
	}

	go func() {
		mac := p.currentMachine()
		response, err := p.submitCode()

		if err != nil {
			p.Outgoing(nwmessage.PsError(err))
			return
		}

		mac.Lock()
		gm.tryResetMachine(p, node, mac, response)
		mac.Unlock()

		p.SendPrompt()
	}()

	p.Outgoing(nwmessage.PsBegin("Resetting machine..."))
	return nil
}

// Async Confirmation Dialogues
// func beginRemoveModuleConf(p *Player, gm *GameModel) {
// 	p.Outgoing(nwmessage.PsDialogue("Resetting friendly machine, (y)es to confirm\nany other key to abort: "))

// 	p.dialogue = nwmessage.NewDialogue([]nwmessage.Fn{
// 		func(d *nwmessage.Dialogue, s string) nwmessage.Message {
// 			if s == "y" || s == "ye" || s == "yes" {
// 				d.SetProp("flag", "-yes")
// 			} else {
// 				d.SetProp("flag", "-no")
// 			}

// 			p.dialogue = nil

// 			return cmdResetMachine(p, gm, []string{d.GetProp("flag")})
// 		},
// 	})
// }
