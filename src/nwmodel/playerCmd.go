package nwmodel

import (
	"errors"
	"fmt"
	"log"
	"nwmessage"
	"sort"
	"strconv"
	"strings"
)

// type playerCommand func(p *Player, gm *GameModel, args []string, code string) nwmessage.Message
type playerCommand func(*Player, *GameModel, []string, string) nwmessage.Message

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
	// disconnect
	// "disconnect": cmdDisconnect,
	// "dis":        cmdDisconnect,

	"lang": cmdLang,

	"langs": cmdListLanguages,

	"make": cmdMake,

	// TODO hacky AF
	"setStdinFromTheFrontEndButNotARealCommand": cmdStdin,

	"pow":   cmdScore,
	"power": cmdScore,

	"test": cmdTestCode,

	"rm":     cmdRemoveModule,
	"remove": cmdRemoveModule,

	// "rf":    cmdRefac,
	// "ref":   cmdRefac,
	// "refac": cmdRefac,

	"at":     cmdAttach,
	"att":    cmdAttach,
	"attach": cmdAttach,

	// "nm": cmdNewMap,
	// what am I attached to? what's my task? what's my language set to?
	// "st":   cmdStatus,
	// "stat": cmdStatus,

	"who": cmdWho,

	"ls": cmdLs, // list modules/slot. out of spec but for expediency

	"sp": cmdSetPOE,
	// "boilerplate": cmdLoadBoilerplate,
	// "bp":          cmdLoadBoilerplate,
}

func actionConsumer(gm *GameModel) {
	for {
		m := <-gm.aChan
		senderID, err := strconv.Atoi(m.Sender)

		if err != nil {
			log.Println(err)
		}

		p := gm.Players[senderID]

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
				continue
			}
			// if the games not locked, allow map to me modified.
			res := handlerFunc(p, gm, msg[1:], m.Code)
			if res.Data != "" {
				p.Outgoing <- res
			}
		} else if handlerFunc, ok := gameCmdList[msg[0]]; ok {
			res := handlerFunc(p, gm, msg[1:], m.Code)
			if res.Data != "" {
				p.Outgoing <- res
			}
		} else {
			p.Outgoing <- nwmessage.PsUnknown(msg[0])
		}
		p.Outgoing <- nwmessage.PromptState(p.Prompt())
	}
}

func cmdYell(p *Player, gm *GameModel, args []string, c string) nwmessage.Message {

	if len(args) == 0 {
		return nwmessage.PsError(errors.New("Need a message to yell"))
	}

	chatMsg := strings.Join(args, " ")

	gm.psBroadcast(nwmessage.PsChat(p.GetName(), "global", chatMsg))
	return nwmessage.Message{}
}

// func cmdTell(p *Player, gm *GameModel, args []string, c string) nwmessage.Message {

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

func cmdTeamChat(p *Player, gm *GameModel, args []string, c string) nwmessage.Message {
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

func cmdSay(p *Player, gm *GameModel, args []string, c string) nwmessage.Message {

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

// func cmdSetName(p *Player, gm *GameModel, args []string, c string) nwmessage.Message {
// 	return nwmessage.PsError(errors.New("Can't change name in a game"))
// }

func cmdJoinTeam(p *Player, gm *GameModel, args []string, c string) nwmessage.Message {
	// log.Println("cmdJoinTeam called")
	// TODO if args[0] == "auto", join smallest team, also use for team
	if len(args) == 0 {
		// if p.TeamName == "" {
		// 	return nwmessage.PsNoTeam()
		// }
		tt := gm.trailingTeam()
		p.Outgoing <- nwmessage.PsNeutral(fmt.Sprintf("Auto-assigning to team '%s'", tt))
		return cmdJoinTeam(p, gm, []string{tt}, c)
	}

	err := gm.assignPlayerToTeam(p, teamName(args[0]))
	if err != nil {
		return nwmessage.PsError(err)
	}

	retStr := fmt.Sprintf("You're on the " + args[0] + " team")

	tp := gm.Teams[p.TeamName].poe

	if tp == nil {
		retStr += "\nYour team doesn't have a point of entry yet.\nUse 'sp node_id' to set one and begin playing"
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

func cmdLang(p *Player, gm *GameModel, args []string, c string) nwmessage.Message {
	if len(args) == 0 {
		return nwmessage.PsError(errors.New("Expected one argument, received zero"))
	}

	err := gm.setLanguage(p, args[0])
	if err != nil {
		return nwmessage.PsError(err)
	}

	// if the player's attached somewhere, update the buffer
	if p.slotNum != -1 {
		if p.slot().Module != nil && p.slot().Module.TeamName != p.TeamName {
			return nwmessage.PsError(errors.New("Can't change language on enemy module"))
		}
		pSlot := p.slot()
		langDetails := gm.languages[p.language]
		boilerplate := langDetails.Boilerplate
		comment := langDetails.CommentPrefix
		sampleIO := pSlot.challenge.SampleIO
		description := pSlot.challenge.Description

		editText := fmt.Sprintf("%s Challenge:\n%s %s\n%s Sample IO: %s\n\n%s", comment, comment, description, comment, sampleIO, boilerplate)
		p.Outgoing <- nwmessage.EditState(editText)
	}

	return nwmessage.PsSuccess(fmt.Sprintf("Language set to %s", args[0]))
}

func cmdListLanguages(p *Player, gm *GameModel, args []string, c string) nwmessage.Message {
	msgContent := ""

	for k := range gm.languages {
		msgContent += k + "\n"
	}

	return nwmessage.PsNeutral(msgContent)
}

func cmdConnect(p *Player, gm *GameModel, args []string, c string) nwmessage.Message {
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

func cmdWho(p *Player, gm *GameModel, args []string, c string) nwmessage.Message {

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

func cmdLs(p *Player, gm *GameModel, args []string, c string) nwmessage.Message {

	if p.Route == nil {
		return nwmessage.PsNoConnection()
	}

	retMsg := p.Route.Endpoint.forMsg()
	pHere := p.Route.Endpoint.playersHere

	if len(pHere) > 1 {
		//make slice of names (excluding this player)
		names := make([]string, 0, len(pHere)-1)
		for _, pID := range pHere {
			player := gm.Players[pID]
			if player.name != p.GetName() {
				names = append(names, fmt.Sprintf("%s (%s)", player.name, player.TeamName))
			}
		}

		//join the slice to string
		addMsg := "\nAlso here: " + strings.Join(names, ", ")

		//add to message
		retMsg += addMsg
	}

	return nwmessage.PsNeutral(retMsg)
}

func cmdSetPOE(p *Player, gm *GameModel, args []string, c string) nwmessage.Message {
	if p.TeamName == "" {
		return nwmessage.PsNoTeam()
	}

	newPOE, err := strconv.Atoi(args[0])
	if err != nil {
		return nwmessage.PsError(fmt.Errorf("expected integer, got '%v'", args[0]))
	}

	err = gm.setTeamPoe(gm.Teams[p.TeamName], newPOE)
	if err != nil {
		return nwmessage.PsError(err)
	}

	// TODO handle initial poe module more elegently
	gm.POEs[p.ID].buildDummyModule(p)
	gm.calcPoweredNodes(gm.Teams[p.TeamName])

	for player := range gm.Teams[p.TeamName].players {
		_, _ = gm.tryConnectPlayerToNode(player, newPOE)
	}

	// if all teams have their poe set
	var ready int
	for _, team := range gm.Teams {
		if team.poe != nil {
			ready++
		}
	}

	// start the game
	if len(gm.Teams) == ready {
		gm.startGame()
	}

	gm.broadcastState()
	return nwmessage.PsSuccess(fmt.Sprintf("%s team's point of entry set to node %d\nConnecting you there now...", p.TeamName, newPOE))
}

// deprecate this in favor of stdin box TODO
func cmdStdin(p *Player, gm *GameModel, args []string, c string) nwmessage.Message {
	// disallow blank stdin
	if len(args) == 0 {
		p.stdin = ""
	}

	p.stdin = strings.Join(args, " ")

	return nwmessage.Message{}
}

func cmdTestCode(p *Player, gm *GameModel, args []string, c string) nwmessage.Message {

	if len(args) > 0 {
		p.stdin = strings.Join(args, " ")
	}

	// TODO handle compiler error
	if c == "" {
		return nwmessage.PsError(errors.New("No code submitted"))
	}

	// passed error checks on args

	go func() {
		response := getOutput(p.language, c, p.stdin)
		p.Outgoing <- nwmessage.PsSuccess("Finished running (check output box)")

		if response.Message.Type == "error" {
			p.Outgoing <- nwmessage.ResultState(fmt.Sprintf("Error: %s", response.Message.Data))
			return
		}

		p.Outgoing <- nwmessage.ResultState(fmt.Sprintf("%s", response.Raw))
	}()

	return nwmessage.PsBegin(fmt.Sprintf("Testing code..."))
}

func cmdScore(p *Player, gm *GameModel, args []string, c string) nwmessage.Message {
	var scoreStrs sort.StringSlice

	for teamName, team := range gm.Teams {
		scoreStrs = append(scoreStrs, fmt.Sprintf("%s:\nProcessing Power: %.2f\nCalculations Completed: %.2f/%.0f", teamName, team.ProcPow, team.VicPoints, gm.PointGoal))
	}

	scoreStrs.Sort()

	return nwmessage.PsNeutral(strings.Join(scoreStrs, "\n"))
}

// TODO refactor cmdAttach for clarity and redundancy
func cmdAttach(p *Player, gm *GameModel, args []string, c string) nwmessage.Message {
	slotNum, err := validateOneIntArg(args)
	if err != nil {
		return nwmessage.PsError(err)
	}

	// if err = validateSlotIs("either", p, slotNum); err != nil {
	// 	return nwmessage.PsError(err)
	// }

	switch {
	case p.Route == nil:
		return nwmessage.PsNoConnection()

	case slotNum > len(p.Route.Endpoint.Slots)-1 || slotNum < 0:
		return nwmessage.PsError(fmt.Errorf("Slot '%v' does not exist", slotNum))
	}

	// passed checks, set player slot to target
	p.slotNum = slotNum
	pSlot := p.slot()

	// if the slot has an enemy module, player's language is set to that module's
	langLock := false
	if pSlot.Module != nil && pSlot.Module.TeamName != p.TeamName {
		langLock = true
		gm.setLanguage(p, pSlot.Module.language)
	}

	// Send slot info to edit buffer
	msgPostfix := "\nchallenge details loaded to codebox"

	// get language details
	langDetails := gm.languages[p.language]
	boilerplate := langDetails.Boilerplate
	comment := langDetails.CommentPrefix
	sampleIO := pSlot.challenge.SampleIO
	description := pSlot.challenge.Description

	// resp := nwmessage.PsPrompt(p.Outgoing, p.Socket, "Overwriting edit buffer with challenge details,\nhit any key to continue, (n) to leave buffer in place: ")
	resp := ""
	if resp != "n" && resp != "no" {
		editText := fmt.Sprintf("%s Challenge:\n%s %s\n%s Sample IO: %s\n\n%s", comment, comment, description, comment, sampleIO, boilerplate)

		if langLock {
			editText += fmt.Sprintf("\n\n%sENEMY MODULE, SOLUTION MUST BE IN [%s]", comment, strings.ToUpper(p.slot().Module.language))
		}

		p.Outgoing <- nwmessage.EditState(editText)
	} else {
		msgPostfix = "\n" + fmt.Sprintf("%s %s\n%s Sample IO: %s", comment, description, comment, sampleIO)
	}

	retText := fmt.Sprintf("Attached to slot %d: \ncontents:%v", slotNum, pSlot.forMsg())
	retText += msgPostfix
	if langLock {
		// TODO add this message to codebox
		retText += fmt.Sprintf("\nalert: SOLUTION MUST BE IN %v", pSlot.Module.language)
	}
	return nwmessage.PsSuccess(retText)
}

func cmdMake(p *Player, gm *GameModel, args []string, c string) nwmessage.Message {

	// TODO handle compiler error
	if c == "" {
		return nwmessage.PsError(errors.New("No code submitted"))
	}

	slot := p.slot()

	if slot == nil {
		return nwmessage.PsError(errors.New("Not attached to slot"))
	}

	// enforce module language
	if slot.Module != nil && slot.Module.TeamName != p.TeamName && slot.Module.language != p.language {
		nwmessage.PsError(fmt.Errorf("This module is written in %s, your code must be written in %s", slot.Module.language, slot.Module.language))
	}

	// passed error checks on args

	go func(p *Player, gm *GameModel, c string) {
		slot := p.slot()
		log.Printf("Make go routine, slot.challenge.ID: %s", slot.challenge.ID)
		response := submitTest(slot.challenge.ID, p.language, c)
		p.compiling = false
		p.Outgoing <- nwmessage.TerminalUnpause()

		if response.Message.Type == "error" {
			p.Outgoing <- nwmessage.PsCompileFail()
			p.Outgoing <- nwmessage.ResultState(fmt.Sprintf("Results: %s\nErrors: %s", response.Graded, response.Message.Data))
			return
		}

		newModHealth := response.passed()

		if newModHealth == 0 {
			p.Outgoing <- nwmessage.PsError(fmt.Errorf("Module failed all tests"))
			p.Outgoing <- nwmessage.ResultState(fmt.Sprintf("%s", response.Graded))
			return
		}
		// if there's no error show graded results, regardless of what happens with the module:
		p.Outgoing <- nwmessage.ResultState(fmt.Sprintf("%s", response.Graded))

		// LOCK SLOT
		slot.Lock()
		defer slot.Unlock()

		if slot.Module != nil {
			// in case we're refactoring a friendly module
			if slot.Module.TeamName == p.TeamName {

				slot.Module.Health = newModHealth
				slot.Module.language = p.language

				gm.broadcastState()
				gm.psBroadcastExcept(p, nwmessage.PsAlert(fmt.Sprintf("%s of (%s) refactored a friendly module in node %d", p.GetName(), p.TeamName, p.Route.Endpoint.ID)))
				p.Outgoing <- nwmessage.PsSuccess(fmt.Sprintf("Refactored friendly module to %d/%d [%s]", slot.Module.Health, slot.Module.MaxHealth, slot.Module.language))
				return
			}

			// hostile module
			switch {
			case newModHealth < slot.Module.Health:
				p.Outgoing <- nwmessage.PsError(fmt.Errorf("Module too weak to install: %d/%d, need at least %d/%d", response.passed(), len(response.Graded), slot.Module.Health, slot.Module.MaxHealth))
				return

			case newModHealth == slot.Module.Health:
				p.Outgoing <- nwmessage.PsAlert(fmt.Sprintf("You need to pass one more test to steal,\nbut your %d/%d is enough to remove.\nKeep trying if you think you can do\nbetter or type 'remove' to proceed", newModHealth, slot.Module.MaxHealth))
				return

			case newModHealth > slot.Module.Health:
				// // track old owner to evaluate traffic after module loss
				oldTeamName := slot.Module.TeamName

				// // refactor module to new owner and health
				// slot.Module.TeamName = p.TeamName
				// slot.Module.Health = newModHealth

				// // evaluate routing of player trffic through node
				// gm.evalTrafficForTeam(p.Route.Endpoint, oldTeam)
				// // ensure new owner is powering node
				// gm.Teams[p.TeamName].powerOn(p.Route.Endpoint)

				gm.refactorModule(slot.Module, p, newModHealth)

				// update map state
				gm.broadcastState()
				gm.broadcastAlertFlash(p.TeamName)

				// broadcast terminal messages
				gm.psBroadcastExcept(p, nwmessage.PsAlert(fmt.Sprintf("%s of (%s) stole a (%s) module in node %d", p.GetName(), p.TeamName, oldTeamName, p.Route.Endpoint.ID)))
				p.Outgoing <- nwmessage.PsSuccess(fmt.Sprintf("You stole (%v)'s module, new module health: %d/%d", oldTeamName, slot.Module.Health, slot.Module.MaxHealth))
				return
			}

		}

		// slot is empty, simply install...
		newMod := newModule(p, response, p.language)

		gm.buildModule(p, newMod)

		gm.broadcastState()
		gm.broadcastAlertFlash(p.TeamName)
		gm.psBroadcastExcept(p, nwmessage.PsAlert(fmt.Sprintf("%s of (%s) constructed a module in node %d", p.GetName(), p.TeamName, p.Route.Endpoint.ID)))
		p.Outgoing <- nwmessage.PsSuccess(fmt.Sprintf("Module constructed in [%s], Health: %d/%d", slot.Module.language, slot.Module.Health, slot.Module.MaxHealth))
		return
	}(p, gm, c)

	p.compiling = true
	p.Outgoing <- nwmessage.TerminalPause()
	return nwmessage.PsBegin("Compiling...")
}

func cmdRemoveModule(p *Player, gm *GameModel, args []string, c string) nwmessage.Message {

	slot := p.slot()

	if slot == nil {
		return nwmessage.PsError(errors.New("Not attached to slot"))
	}

	if slot.Module == nil {
		return nwmessage.PsError(errors.New("Slot is empty"))
	}

	// if we're removing a friendly module, just do it:
	if p.TeamName == slot.Module.TeamName {
		// TODO hacky, refactor
		if len(args) == 0 {
			args = append(args, "")
		}

		flag := args[0]

		log.Printf("cmdRemoveModule flag: %v", flag)
		switch {
		case flag == "":
			beginRemoveModuleConf(p, gm)
			return nwmessage.Message{}

		case flag == "-y" || flag == "-ye" || flag == "-yes":
			// err := p.Route.Endpoint.removeModule(p.slotNum)
			// if err != nil {
			// 	return nwmessage.PsError(err)
			// }

			// gm.evalTrafficForTeam(p.Route.Endpoint, gm.Teams[p.TeamName])
			gm.removeModule(p)
			gm.broadcastState()
			return nwmessage.PsSuccess("Module removed")

		case flag == "-no":
			return nwmessage.PsError(errors.New("Removal aborted"))

		default:
			return nwmessage.PsError(errors.New("Unknown flag, use -y for automatic confirmation"))
		}
	}

	if c == "" {
		return nwmessage.PsError(errors.New("No code submitted"))
	}

	// All checks passed:
	// passed error checks on args

	go func(p *Player, gm *GameModel, c string) {
		slot := p.slot()

		response := submitTest(slot.challenge.ID, p.language, c)

		p.compiling = false
		p.Outgoing <- nwmessage.TerminalUnpause()

		if response.Message.Type == "error" {
			p.Outgoing <- nwmessage.PsCompileFail()
			p.Outgoing <- nwmessage.ResultState(fmt.Sprintf("Results: %s\nErrors: %s", response.Graded, response.Message.Data))
			return
		}

		newModHealth := response.passed()

		if newModHealth == 0 {
			p.Outgoing <- nwmessage.PsError(fmt.Errorf("Module failed all tests"))
			p.Outgoing <- nwmessage.ResultState(fmt.Sprintf("%s", response.Graded))
			return
		}

		// if there's no error, show graded results, regardless of what happens with the module:
		p.Outgoing <- nwmessage.ResultState(fmt.Sprintf("%s", response.Graded))

		// LOCK SLOT
		slot.Lock()
		defer slot.Unlock()

		if newModHealth >= slot.Module.Health {
			oldTeamName := slot.Module.TeamName

			// err := p.Route.Endpoint.removeModule(p.slotNum)
			// if err != nil {
			// 	p.Outgoing <- nwmessage.PsError(err)
			// 	return
			// }
			// gm.evalTrafficForTeam(p.Route.Endpoint, oldTeam)
			gm.removeModule(p)

			gm.broadcastState()
			gm.broadcastAlertFlash(p.TeamName)
			gm.psBroadcastExcept(p, nwmessage.PsAlert(fmt.Sprintf("%s of (%s) removed a (%s) module in node %d", p.GetName(), p.TeamName, oldTeamName, p.Route.Endpoint.ID)))
			p.Outgoing <- nwmessage.PsSuccess("Module removed")
			return

		}

		p.Outgoing <- nwmessage.PsError(fmt.Errorf(
			"Solution too weak: %d/%d, need %d/%d to remove",
			response.passed(), len(response.Graded), slot.Module.Health, slot.Module.MaxHealth,
		))
		return
	}(p, gm, c)

	p.compiling = true
	p.Outgoing <- nwmessage.TerminalPause()
	return nwmessage.PsBegin("Removing module...")

}

func cmdLoadMod(p *Player, gm *GameModel, args []string, c string) nwmessage.Message {
	return nwmessage.Message{}
}

func validateOneIntArg(args []string) (int, error) {
	if len(args) < 1 {
		return 0, fmt.Errorf("Expected 1 argument, received %v", len(args))
	}

	target, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, fmt.Errorf("expected integer, got '%v'", args[0])
	}

	return target, nil
}

// Async Confirmation Dialogues
func beginRemoveModuleConf(p *Player, gm *GameModel) {
	p.Outgoing <- nwmessage.PsDialogue("Removing friendly module, (y)es to confirm\nany other key to abort: ")
	p.dialogue = nwmessage.NewDialogue([]nwmessage.Fn{
		func(d *nwmessage.Dialogue, s string) nwmessage.Message {
			if s == "y" || s == "ye" || s == "yes" {
				d.SetProp("flag", "-yes")
			} else {
				d.SetProp("flag", "-no")
			}

			p.dialogue = nil

			return cmdRemoveModule(p, gm, []string{d.GetProp("flag")}, "")
		},
	})
}

// func validateSlotIs(wants string, p *Player, slotNum int) error {
// 	// check validity of player.Route and slot number
// 	switch {
// 	case p.Route == nil:
// 		return errors.New(noConnectStr)

// 	case slotNum > len(p.Route.Endpoint.slots)-1 || slotNum < 0:
// 		return fmt.Errorf("Slot '%v' does not exist", slotNum)
// 	}

// 	switch wants {
// 	case "full":
// 		if p.Route.Endpoint.slots[slotNum].Module == nil {
// 			return fmt.Errorf("slot '%v' is empty", slotNum)
// 		}
// 		return nil
// 	case "empty":
// 		if p.Route.Endpoint.slots[slotNum].Module != nil {
// 			return fmt.Errorf("slot '%v' is full", slotNum)
// 		}
// 		return nil
// 	}
// 	return nil
// }

// func slotValidateNotEmpty(p *Player, slotNum int) error {
// 	switch {
// 	case p.Route == nil:
// 		return errors.New(noConnectStr)

// 	case slotNum > len(p.Route.Endpoint.slots)-1 || slotNum < 0:
// 		return fmt.Errorf("Slot '%v' does not exist", slotNum)

// 	case p.Route.Endpoint.slots[slotNum].Module == nil:
// 		// log.Printf("slots target: %v", p.Route.Endpoint.slots[target])
// 		return fmt.Errorf("slot '%v' is empty", slotNum)
// 	}
// 	return nil
// }
