package nwmodel

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
)

type playerCommand func(p *Player, args []string, code string) Message

var msgMap = map[string]playerCommand{
	// chat functions
	"y":        cmdYell,
	"yell":     cmdYell,
	"t":        cmdTell,
	"tell":     cmdTell,
	"tc":       cmdTc,
	"teamchat": cmdTc,

	// // player settings
	"team": cmdTeam,
	"name": cmdName,

	// // world interaction
	"con":     cmdConnect,
	"connect": cmdConnect,
	// disconnect
	// "disconnect": cmdDisconnect,
	// "dis":        cmdDisconnect,

	"lang": cmdLanguage,

	"langs": cmdListLanguages,

	"mk":      cmdMake,
	"mak":     cmdMake,
	"make":    cmdMake,
	"makemod": cmdMake,

	"stdin": cmdStdin,

	"test": cmdTestCode,

	"rm":     cmdRemoveModule,
	"remove": cmdRemoveModule,

	// "rf":    cmdRefac,
	// "ref":   cmdRefac,
	// "refac": cmdRefac,

	"at":     cmdAttach,
	"att":    cmdAttach,
	"attach": cmdAttach,

	"nm": cmdNewMap,
	// what am I attached to? what's my task? what's my language set to?
	// "st":   cmdStatus,
	// "stat": cmdStatus,

	"who": cmdWho,

	"ls":          cmdLs, // list modules/slot. out of spec but for expediency
	"listmod":     cmdLs, // list modules/slot. out of spec but for expediency
	"sp":          cmdSetPOE,
	"teampoe":     cmdSetPOE,
	"boilerplate": cmdLoadBoilerplate,
	"bp":          cmdLoadBoilerplate,
}

func cmdHandler(m *Message, p *Player) Message {
	// when we receive a message
	// log.Println("Splitting...")
	msg := strings.Split(m.Data, " ")

	// log.Println("Finding handlerFunc...")
	if handlerFunc, ok := msgMap[msg[0]]; ok {

		// log.Println("Calling handlerFunc")
		return handlerFunc(p, msg[1:], m.Code)
	}

	return psUnknown(msg[0])

}

func cmdYell(p *Player, args []string, c string) Message {

	chatMsg := p.name() + " > " + strings.Join(args, " ")

	gm.psBroadcast(Message{
		Type: "(global)",
		Data: chatMsg,
	})
	return Message{}
}

func cmdTell(p *Player, args []string, c string) Message {

	var recip *Player
	for _, player := range gm.Players {
		if player.Name == args[0] {
			recip = player
		}
	}

	if recip == nil {
		return psError(fmt.Errorf("No such player, '%s'", args[0]))
	}

	chatMsg := p.name() + " > " + strings.Join(args[1:], " ")

	recip.outgoing <- Message{
		Type:   "(private)",
		Data:   chatMsg,
		Sender: pseudoStr,
	}
	return Message{}
}

func cmdTc(p *Player, args []string, c string) Message {
	if p.Team == nil {
		return msgNoTeam
	}

	chatMsg := p.name() + "> " + strings.Join(args, " ")

	p.Team.broadcast(Message{
		Type: "(team)",
		Data: chatMsg,
	})
	return Message{}

}

func cmdTeam(p *Player, args []string, c string) Message {
	// log.Println("cmdTeam called")
	// TODO if args[0] == "auto", join smallest team, also use for team
	if len(args) == 0 {
		if p.Team == nil {
			return msgNoTeam
		}
		return psSuccess(("You're on the " + p.Team.Name + " team"))
	}

	err := gm.assignPlayerToTeam(p, teamName(args[0]))
	if err != nil {
		return psError(err)
	}

	p.outgoing <- Message{
		Type:   "teamState",
		Sender: "server",
		Data:   args[0],
	}

	if p.Team.poe != nil {
		log.Printf("player joined team, tryin to log into %v", p.Team.poe.ID)
		_, err = gm.tryConnectPlayerToNode(p, p.Team.poe.ID)
		if err != nil {
			log.Println(err)
		}
		gm.broadcastState()
	}

	return psSuccess("You're on the " + args[0] + " team")
}

func cmdLanguage(p *Player, args []string, c string) Message {
	if len(args) == 0 {
		return psSuccess("Your name is " + p.language)
	}

	err := p.setLanguage(args[0])
	if err != nil {
		return psError(err)
	}

	// if the player's attached somewhere, update the buffer
	if p.slotNum != -1 {
		p.outgoing <- editStateMsg(boilerPlateFor(p) + challengeBufferFor(p))
	}

	return psSuccess(fmt.Sprintf("Language set to %s", args[0]))
}

func cmdListLanguages(p *Player, args []string, c string) Message {
	msgContent := ""

	for k := range gm.languages {
		msgContent += k + "\n"
	}

	return psMessage(msgContent)
}

func cmdLoadBoilerplate(p *Player, args []string, c string) Message {
	p.outgoing <- editStateMsg(boilerPlateFor(p))
	return psSuccess(fmt.Sprintf("%s boilerplate loaded", p.language))
}

func cmdName(p *Player, args []string, c string) Message {
	if len(args) == 0 {
		return psSuccess("Your name is " + p.name())
	}

	err := gm.setPlayerName(p, args[0])
	if err != nil {
		return psError(err)
	}

	return psSuccess("Name set to '" + p.name() + "'")
}

func cmdConnect(p *Player, args []string, c string) Message {
	if p.Team == nil {
		return psError(errors.New("Join a team first"))
	}

	if len(args) == 0 {
		if p.Route != nil {
			return psSuccess(fmt.Sprintf("You're connected to node %d which connects to %v", p.Route.Endpoint.ID, p.Route.Endpoint.Connections))
		}
		return psError(errors.New("No connection, provide a nodeID argument to connect"))
	}

	targetNode, err := strconv.Atoi(args[0])
	if err != nil {
		return psError(err)
	}

	// break any pre-existing connection before connecting elsewhere
	route, err := gm.tryConnectPlayerToNode(p, targetNode)
	if err != nil {
		return psError(err)
	}

	gm.broadcastState()
	return psSuccess(fmt.Sprintf("Connected to established : %s", route.forMsg()))
}

func cmdWho(p *Player, args []string, c string) Message {
	// lists all players in the current node
	// if p.Route.Endpoint == nil {
	// 	return psError(errors.New(noConnectStr))
	// }

	// // TODO maintain a list of connected players at either node or slot
	// pHere := ""
	// for _, otherPlayer := range gm.Players {
	// 	if otherPlayer.Route != nil {
	// 		if otherPlayer.Route.Endpoint == p.Route.Endpoint {
	// 			playerDesc := otherPlayer.name()
	// 			if otherPlayer.slotNum > -1 {
	// 				playerDesc += " at slot: " + strconv.Itoa(otherPlayer.slotNum)
	// 			}
	// 			pHere += playerDesc + "\n"
	// 		}
	// 	}
	// }

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
			whoStr += "\t" + mem.Name + "\n"
		}
	}

	return psMessage(whoStr)
}

func cmdLs(p *Player, args []string, c string) Message {
	if p.Route == nil {
		return msgNoConnection
	}

	retMsg := p.Route.Endpoint.forMsg()
	pHere := p.Route.Endpoint.playersHere
	if len(pHere) > 1 {
		//make slice of names (excluding this player)
		names := make([]string, 0, len(pHere)-1)
		for _, player := range pHere {
			if player != p {
				names = append(names, player.Name)
			}
		}

		//join the slice to string
		addMsg := "\nAlso here: " + strings.Join(names, ", ")

		//add to message
		retMsg += addMsg
	}
	return psMessage(retMsg)
}

func cmdSetPOE(p *Player, args []string, c string) Message {
	if p.Team == nil {
		return msgNoTeam
	}

	newPOE, err := strconv.Atoi(args[0])
	if err != nil {
		return psError(fmt.Errorf("expected integer, got '%v'", args[0]))
	}

	err = gm.setTeamPoe(p.Team, newPOE)
	if err != nil {
		return psError(err)
	}

	// debug only :
	// fix this TODO
	_ = gm.POEs[p.ID].addModule(newModule(p, ChallengeResponse{}, p.language), 0)

	for player := range p.Team.players {
		_, _ = gm.tryConnectPlayerToNode(player, newPOE)
	}

	// TODO connect players on POE set. that means we have to create the module before each player's POE gets set...
	// which is happening right now as a result of setTeamPoe

	gm.broadcastState()
	return psSuccess(fmt.Sprintf("Team %s's point of entry set to node %d", p.Team.Name, newPOE))
}

func cmdStdin(p *Player, args []string, playerCode string) Message {
	// disallow blank stdin
	if p.stdin == "" {
		p.stdin = "default stdin"
	}

	if len(args) == 0 {
		return psMessage("stdin is: " + p.stdin)
	}

	p.stdin = strings.Join(args, " ")

	return psMessage("stdin set to: " + p.stdin)
}

func cmdTestCode(p *Player, args []string, playerCode string) Message {

	if len(args) > 0 {
		p.stdin = strings.Join(args, " ")
	}

	// TODO handle compiler error
	if playerCode == "" {
		return psError(errors.New("No code submitted"))
	}

	if p.stdin == "" {
		p.stdin = "default stdin"
	}

	// passed error checks on args
	p.outgoing <- psBegin(fmt.Sprintf("Running test with stdin: %v", p.stdin))

	// disallow blank stdin
	if p.stdin == "" {
		p.stdin = "default stdin"
	}

	response := getOutput(p.language, playerCode, p.stdin)

	if response.Message.Type == "error" {
		return psError(errors.New(response.Message.Data))
	}

	return psSuccess(fmt.Sprintf("Output: %v", response))
}

// TODO refactor cmdAttach for clarity and redundancy
func cmdAttach(p *Player, args []string, playerCode string) Message {
	slotNum, err := validateOneIntArg(args)
	if err != nil {
		return psError(err)
	}

	if err = validateSlotIs("either", p, slotNum); err != nil {
		return psError(err)
	}

	// passed checks, set player slot to target
	p.slotNum = slotNum
	pSlot := p.slot()

	// if the slot has an enemy module, player's language is set to that module's
	langLock := false
	if pSlot.module != nil && pSlot.module.Team != p.Team {
		langLock = true
		p.setLanguage(pSlot.module.language)
	}

	// Send slot info to edit buffer
	msgPostfix := "\nchallenge details loaded to codebox"

	// get language details
	langDetails := gm.languages[p.language]
	boilerplate := langDetails.Boilerplate
	comment := langDetails.CommentPrefix
	sampleIO := pSlot.challenge.SampleIO
	description := pSlot.challenge.Description

	resp := psPrompt(p, "Overwriting edit buffer with challenge details,\nhit any key to continue, (n) to leave buffer in place: ")
	if resp != "n" && resp != "no" {
		editMessage := fmt.Sprintf("%s\n%s %s\n%s Sample IO: %s", boilerplate, comment, description, comment, sampleIO)
		if langLock {
			editMessage += fmt.Sprintf("\n\n%sENEMY MODULE, SOLUTION MUST BE IN [%s]", comment, strings.ToUpper(p.slot().module.language))
		}

		p.outgoing <- editStateMsg(editMessage)
	} else {
		msgPostfix = "\n" + challengeBufferFor(p)
	}

	retText := fmt.Sprintf("Attached to slot %d: \ncontents:%v", slotNum, pSlot.forMsg())
	retText += msgPostfix
	if langLock {
		// TODO add this message to codebox
		retText += fmt.Sprintf("\nalert: SOLUTION MUST BE IN %v", pSlot.module.language)
	}
	return psSuccess(retText)
}

func boilerPlateFor(p *Player) string {
	langDetails := gm.languages[p.language]
	return langDetails.Boilerplate
}

func challengeBufferFor(p *Player) string {
	langDetails := gm.languages[p.language]
	pSlot := p.slot()
	if pSlot == nil {
		return ""
	}

	comment := langDetails.CommentPrefix
	sampleIO := pSlot.challenge.SampleIO
	description := pSlot.challenge.Description
	return fmt.Sprintf("%s %s\n%s Sample IO: %s", comment, description, comment, sampleIO)
}

func cmdMake(p *Player, args []string, playerCode string) Message {

	// TODO handle compiler error
	if playerCode == "" {
		return psError(errors.New("No code submitted"))
	}

	slot := p.slot()

	if slot == nil {
		return psError(errors.New("Not attached to slot"))
	}

	// enforce module language
	if slot.module != nil && slot.module.Team != p.Team && slot.module.language != p.language {
		psError(fmt.Errorf("This module is written in %s, your code must be written in %s", slot.module.language, slot.module.language))
	}

	// passed error checks on args
	p.outgoing <- psBegin("Compiling...")

	response := submitTest(slot.challenge.ID, p.language, playerCode)

	if response.Message.Type == "error" {
		return psError(errors.New(response.Message.Data))
	}

	newModHealth := response.passed()

	if newModHealth == 0 {
		return psError(fmt.Errorf("Failed to make module, test results: %d/%d", response.passed(), len(response.PassFail)))
	}

	if slot.module != nil {
		// in case we're refactoring a friendly module
		if slot.module.Team == p.Team {

			slot.module.Health = newModHealth
			slot.module.language = p.language

			gm.broadcastState()
			gm.psBroadcastExcept(p, psAlert(fmt.Sprintf("%s of (%s) refactored a friendly module in node %d", p.Name, p.Team.Name, p.Route.Endpoint.ID)))
			return psSuccess(fmt.Sprintf("Refactored friendly module to %d/%d [%s]", slot.module.Health, slot.module.MaxHealth, slot.module.language))
		}

		// hostile module
		switch {
		case newModHealth < slot.module.Health:
			return psError(fmt.Errorf("Module too weak to install: %d/%d, need at least %d/%d", response.passed(), len(response.PassFail), slot.module.Health, slot.module.MaxHealth))

		case newModHealth == slot.module.Health:
			return psAlert("You need to pass one more test to steal,\nbut your %d/%d is enough to remove.\nKeep trying if you think you can do\nbetter or type 'remove' to proceed")

		case newModHealth > slot.module.Health:
			oldTeam := slot.module.Team
			slot.module.Team = p.Team
			slot.module.Health = newModHealth
			gm.broadcastState()
			gm.broadcastAlertFlash(p.Team.Name)
			gm.psBroadcastExcept(p, psAlert(fmt.Sprintf("%s of (%s) stole a (%s) module in node %d", p.Name, p.Team.Name, oldTeam.Name, p.Route.Endpoint.ID)))
			return psSuccess(fmt.Sprintf("You stole (%v)'s module, new module health: %d/%d", oldTeam.Name, slot.module.Health, slot.module.MaxHealth))
		}

	}
	// slot is empty, simply install...
	newMod := newModule(p, response, p.language)
	err := p.Route.Endpoint.addModule(newMod, p.slotNum)
	if err != nil {
		return psError(err)
	}
	gm.broadcastState()
	gm.broadcastAlertFlash(p.Team.Name)
	gm.psBroadcastExcept(p, psAlert(fmt.Sprintf("%s of (%s) constructed a module in node %d", p.Name, p.Team.Name, p.Route.Endpoint.ID)))
	return psSuccess(fmt.Sprintf("Module constructed in [%s], Health: %d/%d", slot.module.language, slot.module.Health, slot.module.MaxHealth))
}

func cmdNewMap(p *Player, args []string, playerCode string) Message {
	// TODO fix d3 to update...
	nodeCount, err := validateOneIntArg(args)
	if err != nil {
		return psError(err)
	}

	nodeIDCount = 0
	gm.Map = newRandMap(nodeCount)
	p.outgoing <- Message{
		Type:   "graphReset",
		Sender: serverStr,
		Data:   "",
	}
	gm.broadcastState()
	return psSuccess("Generating new map...")
}

func cmdRemoveModule(p *Player, args []string, playerCode string) Message {

	slot := p.slot()

	if slot == nil {
		return psError(errors.New("Not attached to slot"))
	}

	if !slot.isFull() {
		return psError(errors.New("Slot is empty"))
	}

	// if we're removing a friendly module, just do it:
	if p.Team == slot.module.Team {
		resp := psPrompt(p, "Friendly module, confirm removal? (y/n)")
		if resp == "y" || resp == "ye" || resp == "yes" {
			err := p.Route.Endpoint.removeModule(p.slotNum)
			if err != nil {
				return psError(err)
			}

			gm.broadcastState()
			return psSuccess("Module removed")
		}
		return psError(errors.New("removal aborted"))
	}

	if playerCode == "" {
		return psError(errors.New("No code submitted"))
	}

	// All checks passed:
	// passed error checks on args

	p.outgoing <- psBegin("Removing module...")

	log.Printf("remove cID: %v", slot.challenge.ID)
	response := submitTest(slot.challenge.ID, p.language, playerCode)

	if response.Message.Type == "error" {
		return psError(errors.New(response.Message.Data))
	}

	newModHealth := response.passed()

	log.Printf("response to submitted test: %v", response)

	if newModHealth >= slot.module.Health {
		oldTeamName := slot.module.Team.Name

		err := p.Route.Endpoint.removeModule(p.slotNum)
		if err != nil {
			return psError(err)
		}

		gm.broadcastState()
		gm.broadcastAlertFlash(p.Team.Name)
		gm.psBroadcastExcept(p, psAlert(fmt.Sprintf("%s of (%s) removed a (%s) module in node %d", p.Name, p.Team.Name, oldTeamName, p.Route.Endpoint.ID)))
		return psSuccess("Module removed")

	}

	return psError(fmt.Errorf(
		"Solution too weak: %d/%d, need %d/%d to remove",
		response.passed(), len(response.PassFail), slot.module.Health, slot.module.MaxHealth,
	))

}

func cmdLoadMod(p *Player, args []string, c string) Message {
	return Message{}
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

func validateSlotIs(wants string, p *Player, slotNum int) error {
	// check validity of player.Route and slot number
	switch {
	case p.Route == nil:
		return errors.New(noConnectStr)

	case slotNum > len(p.Route.Endpoint.slots)-1 || slotNum < 0:
		return fmt.Errorf("Slot '%v' does not exist", slotNum)
	}

	switch wants {
	case "full":
		if p.Route.Endpoint.slots[slotNum].module == nil {
			return fmt.Errorf("slot '%v' is empty", slotNum)
		}
		return nil
	case "empty":
		if p.Route.Endpoint.slots[slotNum].module != nil {
			return fmt.Errorf("slot '%v' is full", slotNum)
		}
		return nil
	}
	return nil
}

// func slotValidateNotEmpty(p *Player, slotNum int) error {
// 	switch {
// 	case p.Route == nil:
// 		return errors.New(noConnectStr)

// 	case slotNum > len(p.Route.Endpoint.slots)-1 || slotNum < 0:
// 		return fmt.Errorf("Slot '%v' does not exist", slotNum)

// 	case p.Route.Endpoint.slots[slotNum].module == nil:
// 		// log.Printf("slots target: %v", p.Route.Endpoint.slots[target])
// 		return fmt.Errorf("slot '%v' is empty", slotNum)
// 	}
// 	return nil
// }
