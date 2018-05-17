package nwmodel

import (
	"errors"
	"fmt"
	"log"
	"nwmessage"
	"sort"
	"strconv"
	"strings"

	"feature"
)

// type playerCommand func(p *Player, gm *GameModel, args []string) nwmessage.Message
type playerCommand func(*Player, *GameModel, []string) nwmessage.Message

var gameCmdList = map[string]playerCommand{
	// chat functions
	"y":    cmdYell,
	"yell": cmdYell,
	// "t":        cmdTell,
	// "tell":     cmdTell,
	"tc": cmdTeamChat,
	// "teamchat": cmdTeamChat,

	// "s":   cmdSay,
	"say": cmdSay,

	// // player settings
	"join": cmdJoinTeam,
	"team": cmdJoinTeam,

	// Just an error, setname only works in lobby. should be handled differently TODO
	// "name": cmdSetName,

	// // world interaction
	"con": cmdConnect,

	"lang": cmdLang,

	"langs": cmdListLanguages,

	"make": cmdMake,

	"pow":   cmdScore,
	"power": cmdScore,

	"test": cmdTestCode,

	"res":   cmdResetMachine,
	"reset": cmdResetMachine,

	"at":     cmdAttach,
	"att":    cmdAttach,
	"attach": cmdAttach,

	"who": cmdWho,

	"ls": cmdLs,
	// "ls": cmdDetail{
	// 		usage: "an n target",
	// 		description: "Adds n new nodes linked to target node",
	// 		argsReq: []argType{Int, Int}
	// 	},

	// "sp": cmdSetPOE,
}

func (gm *GameModel) parseCommand(m ClientMessage) {
	p := m.Sender

	msg := strings.Split(m.Data, " ")

	// TODO clean nightmare below
	if p.compiling != false {
		// Would be more elegant to freeze prompt while this happens....
		p.Outgoing <- nwmessage.PsError(errors.New("Code compiling. Wait for completion..."))
	} else if p.dialogue != nil {
		p.Outgoing <- p.dialogue.Run(msg[0])
	} else if handlerFunc, ok := mapCmdList[msg[0]]; ok {
		if gm.running {
			// make this message more situation agnostic TODO
			p.Outgoing <- nwmessage.PsError(errors.New("Cannot alter map once game has started"))
			return
		}
		// if the games not locked, allow map to me modified.
		res := handlerFunc(p, gm, msg[1:])
		if res.Data != "" {
			p.Outgoing <- res
		}
	} else if handlerFunc, ok := gameCmdList[msg[0]]; ok {
		res := handlerFunc(p, gm, msg[1:])
		if res.Data != "" {
			p.Outgoing <- res
		}
	} else {
		p.Outgoing <- nwmessage.PsUnknown(msg[0])
	}

	if p.dialogue == nil {
		p.SendPrompt()
	}
}

// func actionConsumer(gm *GameModel) {
// 	for {
// 		m := <-gm.aChan
// 		senderID, err := strconv.Atoi(m.Sender)

// 		if err != nil {
// 			log.Println(err)
// 		}

// 		p := gm.Players[senderID]

// 		msg := strings.Split(m.Data, " ")

// 		// TODO clean nightmare below
// 		if p.compiling != false {
// 			// Would be more elegant to freeze prompt while this happens....
// 			p.Outgoing <- nwmessage.PsError(errors.New("Code compiling. Wait for completion..."))
// 		} else if p.dialogue != nil {
// 			p.Outgoing <- p.dialogue.Run(msg[0])
// 		} else if handlerFunc, ok := mapCmdList[msg[0]]; ok {
// 			if gm.running {
// 				// make this message more situation agnostic TODO
// 				p.Outgoing <- nwmessage.PsError(errors.New("Cannot alter map once game has started"))
// 				continue
// 			}
// 			// if the games not locked, allow map to me modified.
// 			res := handlerFunc(p, gm, msg[1:])
// 			if res.Data != "" {
// 				p.Outgoing <- res
// 			}
// 		} else if handlerFunc, ok := gameCmdList[msg[0]]; ok {
// 			res := handlerFunc(p, gm, msg[1:])
// 			if res.Data != "" {
// 				p.Outgoing <- res
// 			}
// 		} else {
// 			p.Outgoing <- nwmessage.PsUnknown(msg[0])
// 		}

// 		if p.dialogue == nil {
// 			p.SendPrompt()
// 		}
// 	}
// }

func cmdYell(p *Player, gm *GameModel, args []string) nwmessage.Message {

	if len(args) == 0 {
		return nwmessage.PsError(errors.New("Need a message to yell"))
	}

	chatMsg := strings.Join(args, " ")

	gm.psBroadcast(nwmessage.PsChat(p.GetName(), "global", chatMsg))
	return nwmessage.Message{}
}

// func cmdTell(p *Player, gm *GameModel, args []string) nwmessage.Message {

// 	if len(args) < 2 {
// 		return nwmessage.PsError(errors.New("Need a recipient and a message"))
// 	}
// 	var recip *Player
// 	for _, player := range gm.Players {
// 		if player.GetName() == args[0] {
// 			recip = player
// 		}
// 	}

// 	if recip == nil {
// 		return nwmessage.PsError(fmt.Errorf("No such player, '%s'", args[0]))
// 	}

// 	chatMsg := p.GetName() + " > " + strings.Join(args[1:], " ")

// 	recip.Outgoing <- nwmessage.PsChat(chatMsg, "(private)")
// 	return nwmessage.Message{}
// }

func cmdTeamChat(p *Player, gm *GameModel, args []string) nwmessage.Message {
	if p.TeamName == "" {
		return nwmessage.PsNoTeam()
	}

	if len(args) == 0 {
		return nwmessage.PsError(errors.New("Need a message for team chat"))
	}

	chatMsg := strings.Join(args, " ")

	gm.Teams[p.TeamName].broadcast(nwmessage.PsChat(p.GetName(), "team", chatMsg))

	return nwmessage.Message{}

}

func cmdSay(p *Player, gm *GameModel, args []string) nwmessage.Message {

	if p.Route == nil {
		return nwmessage.PsError(errors.New("Can only 'say' while connected to a node"))
	}

	if len(args) == 0 {
		return nwmessage.PsError(errors.New("Need a message to say"))
	}

	chatMsg := strings.Join(args, " ")

	msg := nwmessage.PsChat(p.GetName(), "node", chatMsg)

	for _, pID := range p.Route.Endpoint.playersHere {
		gm.Players[pID].Outgoing <- msg
	}

	return nwmessage.Message{}
}

// func cmdSetName(p *Player, gm *GameModel, args []string) nwmessage.Message {
// 	return nwmessage.PsError(errors.New("Can't change name in a game"))
// }

func cmdJoinTeam(p *Player, gm *GameModel, args []string) nwmessage.Message {
	// log.Println("cmdJoinTeam called")
	// TODO if args[0] == "auto", join smallest team, also use for team
	if len(args) == 0 || args[0] == "" {
		tt := gm.trailingTeam()
		p.Outgoing <- nwmessage.PsNeutral(fmt.Sprintf("Auto-assigning to team '%s'", tt))
		return cmdJoinTeam(p, gm, []string{tt})
	}

	err := gm.assignPlayerToTeam(p, teamName(args[0]))
	if err != nil {
		return nwmessage.PsError(err)
	}

	retStr := fmt.Sprintf("You're on the " + args[0] + " team")

	if len(gm.Teams[p.TeamName].poes) < 1 {
		retStr += "\nYour team doesn't have a point of entry yet.\nUse 'sp node_id' to set one and begin playing"
	}

	var tp *node
	for n := range gm.Teams[p.TeamName].poes {
		tp = n
		break
	}

	if tp != nil {
		retStr += fmt.Sprintf("\nTeam's point of entry is node %d.\nConnecting you there now...", tp.ID)
		// log.Printf("player joined team, trying to log into %v", tp.ID)
		_, err = gm.tryConnectPlayerToNode(p, tp.ID)
		if err != nil {
			log.Println(err)
		}
		gm.broadcastState()
	}

	p.Outgoing <- nwmessage.TeamState(p.TeamName)
	return nwmessage.PsSuccess(retStr)
}

func cmdLang(p *Player, gm *GameModel, args []string) nwmessage.Message {
	if len(args) == 0 {
		return nwmessage.PsError(errors.New("Expected one argument, received zero"))
	}

	err := gm.setLanguage(p, args[0])
	if err != nil {
		return nwmessage.PsError(err)
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
		p.Outgoing <- nwmessage.EditState(editText)
	}

	return nwmessage.PsSuccess(fmt.Sprintf("Language set to %s", args[0]))
}

func cmdListLanguages(p *Player, gm *GameModel, args []string) nwmessage.Message {
	var langs sort.StringSlice

	for l := range gm.languages {
		langs = append(langs, l)
	}

	langs.Sort()

	return nwmessage.PsNeutral("This game supports:\n" + strings.Join(langs, "\n"))
}

func cmdConnect(p *Player, gm *GameModel, args []string) nwmessage.Message {
	if p.TeamName == "" {
		return nwmessage.PsError(errors.New("Join a team first"))
	}

	if len(args) == 0 {
		return nwmessage.PsError(errors.New("Need a node ID to connect to"))
	}

	targetNode, err := strconv.Atoi(args[0])
	if err != nil {
		return nwmessage.PsError(err)
	}

	// break any pre-existing connection before connecting elsewhere
	route, err := gm.tryConnectPlayerToNode(p, targetNode)
	if err != nil {
		return nwmessage.PsError(err)
	}

	gm.broadcastState()
	return nwmessage.PsSuccess(fmt.Sprintf("Connected to established : %s", route.forMsg()))
}

func cmdWho(p *Player, gm *GameModel, args []string) nwmessage.Message {

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

	return nwmessage.PsNeutral(whoStr)
}

func cmdLs(p *Player, gm *GameModel, args []string) nwmessage.Message {

	if p.Route == nil {
		return nwmessage.PsNoConnection()
	}

	retMsg := p.Route.Endpoint.StringFor(p)
	pHere := p.Route.Endpoint.playersHere

	if len(pHere) > 1 {
		//make slice of names (excluding this player)
		names := make([]string, 0, len(pHere)-1)
		for _, pID := range pHere {
			player := gm.Players[pID]
			if player.GetName() != p.GetName() {
				names = append(names, fmt.Sprintf("%s (%s)", player.GetName(), player.TeamName))
			}
		}

		//join the slice to string
		addMsg := "\nAlso here: " + strings.Join(names, ", ")

		//add to message
		retMsg += addMsg
	}

	return nwmessage.PsNeutral(retMsg)
}

// func cmdSetPOE(p *Player, gm *GameModel, args []string) nwmessage.Message {
// 	if p.TeamName == "" {
// 		return nwmessage.PsNoTeam()
// 	}

// 	nodeId, err := strconv.Atoi(args[0])
// 	if err != nil {
// 		return nwmessage.PsError(fmt.Errorf("expected integer, got '%v'", args[0]))
// 	}

// 	newPoe, err = gm.Map.getNode(nodeId)
// 	if err != nil {
// 		return nwmessage.PsError(err)
// 	}

// 	err = gm.Teams[p.TeamName].addPoe(newPOE)
// 		if err != nil {
// 		return nwmessage.PsError(err)
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
// 	return nwmessage.PsSuccess(fmt.Sprintf("%s team's point of entry set to node %d\nConnecting you there now...", p.TeamName, newPOE))
// }

func cmdTestCode(p *Player, gm *GameModel, args []string) nwmessage.Message {

	if p.EditorState == "" {
		return nwmessage.PsError(errors.New("No code submitted"))
	}

	go func() {
		defer p.SendPrompt()

		response := getOutput(p.language, p.EditorState, p.StdinState)
		p.Outgoing <- nwmessage.PsSuccess("Finished running (check output box)")
		p.Outgoing <- nwmessage.ResultState(response)

	}()

	return nwmessage.PsBegin(fmt.Sprintf("Testing code..."))
}

func cmdScore(p *Player, gm *GameModel, args []string) nwmessage.Message {
	var scoreStrs sort.StringSlice

	for teamName, team := range gm.Teams {
		scoreStrs = append(scoreStrs, fmt.Sprintf("%s:\nCoin Production: %.2f\nCoin Stockpiled: %.2f/%.0f", teamName, team.coinPerTick, team.VicPoints, gm.PointGoal))
	}

	scoreStrs.Sort()

	return nwmessage.PsNeutral(strings.Join(scoreStrs, "\n"))
}

func cmdAttach(p *Player, gm *GameModel, args []string) nwmessage.Message {

	if len(args) < 1 {
		return nwmessage.PsError(errors.New("Attach requires one argument"))
	}

	if p.Route == nil {
		return nwmessage.PsNoConnection()
	}

	macAddress := args[0]
	_, addOk := p.Route.Endpoint.addressMap[macAddress]

	if !addOk {
		return nwmessage.PsError(fmt.Errorf("Invalid address, '%s'", macAddress))
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
		p.Outgoing <- nwmessage.LangSupportState([]string{mac.language})
	} else {
		supportedLangs := make([]string, len(gm.languages))
		var i int
		for lang := range gm.languages {
			supportedLangs[i] = lang
			i++
		}

		p.Outgoing <- nwmessage.LangSupportState(supportedLangs)
	}

	// get language details
	langDetails := gm.languages[p.language]
	boilerplate := langDetails.Boilerplate
	comment := langDetails.CommentPrefix

	p.challengeState(mac.challenge)
	p.stdinState(mac.challenge.SampleIO[0].Input)

	// resp := nwmessage.PsPrompt(p.Outgoing, p.Socket, "Overwriting edit buffer with challenge details,\nhit any key to continue, (n) to leave buffer in place: ")
	var resp string
	if len(args) > 1 {
		resp = args[1]
	}

	var lockStr string
	if langLock {
		lockStr = fmt.Sprintf("\n\n%sHOSTILE MACHINE, SOLUTION MUST BE IN [%s]", comment, strings.ToUpper(p.currentMachine().language))
	}

	// machineStr := "machine"
	// if mac.isFeature() {
	// 	machineStr = "feature"
	// }

	if resp != "-n" && resp != "-no" {
		editText := boilerplate + lockStr
		p.editState(editText)
	}

	retText := fmt.Sprintf("Attached to machine at %s: \ncontents:%v", macAddress, mac.StringFor(p))
	retText += lockStr

	return nwmessage.PsSuccess(retText)
}

func cmdMake(p *Player, gm *GameModel, args []string) nwmessage.Message {

	err := p.canSubmit()
	if err != nil {
		return nwmessage.PsError(err)
	}

	mac := p.currentMachine()

	// Abstract this TODO
	var feaType feature.Type
	if mac.isFeature() {
		if mac.Type == feature.None {
			if len(args) < 1 {
				return nwmessage.PsError(errors.New("Make requires one argument when attached to an untyped feature"))
			}

			var err error
			feaType, err = feature.FromString(args[0])

			if err != nil {
				return nwmessage.PsError(fmt.Errorf("Invalid feature type, '%s'", args[0]))
			}
		} else {
			if len(args) > 0 {
				p.Outgoing <- nwmessage.PsError(errors.New("Ignoring argument! Cannot change the type on an installed feature"))
			}
		}
	}

	// passed error checks
	go func() {
		response, err := p.submitCode()

		if err != nil {
			p.Outgoing <- nwmessage.PsError(err)
			return
		}

		mac.Lock()
		gm.tryClaimMachine(p, response, feaType)
		mac.Unlock()

		p.SendPrompt()
	}()

	return nwmessage.PsBegin("Compiling...")
}

func cmdResetMachine(p *Player, gm *GameModel, args []string) nwmessage.Message {
	err := p.canSubmit()
	if err != nil {
		return nwmessage.PsError(err)
	}

	go func() {
		mac := p.currentMachine()
		response, err := p.submitCode()

		if err != nil {
			p.Outgoing <- nwmessage.PsError(err)
			return
		}

		mac.Lock()
		gm.tryResetMachine(p, response)
		mac.Unlock()

		p.SendPrompt()
	}()

	return nwmessage.PsBegin("Resetting machine...")
}

// Async Confirmation Dialogues
func beginRemoveModuleConf(p *Player, gm *GameModel) {
	p.Outgoing <- nwmessage.PsDialogue("Resetting friendly machine, (y)es to confirm\nany other key to abort: ")
	p.dialogue = nwmessage.NewDialogue([]nwmessage.Fn{
		func(d *nwmessage.Dialogue, s string) nwmessage.Message {
			if s == "y" || s == "ye" || s == "yes" {
				d.SetProp("flag", "-yes")
			} else {
				d.SetProp("flag", "-no")
			}

			p.dialogue = nil

			return cmdResetMachine(p, gm, []string{d.GetProp("flag")})
		},
	})
}
