package nwmodel

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type playerCommand func(p *Player, args []string, code string) Message

var msgMap = map[string]playerCommand{
	// chat functions
	"y":    cmdYell,
	"yell": cmdYell,
	"t":    cmdTell,
	"tell": cmdTell,
	"tc":   cmdTc,

	// // player settings
	"team": cmdTeam,
	"name": cmdName,

	// // world interaction
	"con":     cmdConnect,
	"connect": cmdConnect,

	"lang": cmdLanguage,

	"langs": cmdListLanguages,

	"mk":      cmdMake,
	"mak":     cmdMake,
	"make":    cmdMake,
	"makemod": cmdMake,

	"std":   cmdStdin,
	"stdin": cmdStdin,

	"tst":  cmdTestCode,
	"test": cmdTestCode,

	"rm": cmdRemoveModule,

	"rf":    cmdRefac,
	"ref":   cmdRefac,
	"refac": cmdRefac,

	"at":     cmdAttach,
	"attach": cmdAttach,

	// what am I attached to? what's my task? what's my language set to?
	// "st":   cmdStatus,
	// "stat": cmdStatus,

	"wh":  cmdWho,
	"who": cmdWho,
	// non finalized
	// "pr":    cmdProbe,
	// "probe": cmdProbe,

	"ls": cmdLs, // list modules/slot. out of spec but for expediency
	"sp": cmdSetPOE,
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

	return psSuccess("You're on the " + args[0] + " team")
}

func cmdLanguage(p *Player, args []string, c string) Message {
	if len(args) == 0 {
		return psSuccess("Your name is " + p.language)
	}

	p.setLanguage(args[0])

	return Message{
		Type:   "languageState",
		Sender: "server",
		Data:   p.language,
	}
}

func cmdListLanguages(p *Player, args []string, c string) Message {
	msgContent := ""

	for k := range gm.languages {
		msgContent += k + "\n"
	}

	return psSuccess(msgContent)
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

	if len(args) < 1 {
		return psError(errors.New("Expected 1 argument, received 0"))
	}

	targetNode, err := strconv.Atoi(args[0])
	if err != nil {
		return psError(err)
	}

	route, err := gm.tryConnectPlayerToNode(p, targetNode)
	if err != nil {
		return psError(err)
	}

	gm.broadcastState()
	return psSuccess(fmt.Sprintf("Connected to established : %s", route.forMsg()))
}

func cmdWho(p *Player, args []string, c string) Message {
	// lists all players in the current node
	if p.Route.Endpoint == nil {
		return psError(errors.New(noConnectStr))
	}

	// TODO maintain a list of connected players at either node or slot
	pHere := ""
	for _, otherPlayer := range gm.Players {
		if otherPlayer.Route != nil {
			if otherPlayer.Route.Endpoint == p.Route.Endpoint {
				playerDesc := otherPlayer.name()
				if otherPlayer.slotNum > -1 {
					playerDesc += " at slot: " + strconv.Itoa(otherPlayer.slotNum)
				}
				pHere += playerDesc + "\n"
			}
		}
	}
	return psSuccess(pHere)
}

func cmdLs(p *Player, args []string, c string) Message {
	if p.Route == nil {
		return msgNoConnection
	}
	return psSuccess(p.Route.Endpoint.forMsg())
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
	// if the slot has a module, player's language is set to that module's
	if pSlot.module != nil {
		p.setLanguage(pSlot.module.language)
	}

	//TODO: build buffer from language boilerplate and test sampleIO

	return psSuccess(fmt.Sprintf("Attached to slot %d: %v, Working in: %s", slotNum, pSlot.forProbe(), p.language))
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

	if slot.isFull() {
		return psError(errors.New("Slot is not empty"))
	}

	// passed error checks on args
	p.outgoing <- psBegin("Making module...")

	response := submitTest(slot.challenge.ID, p.language, playerCode)

	if response.Message.Type == "error" {
		return psError(errors.New(response.Message.Data))
	}

	newModHealth := response.passed()

	log.Printf("response to submitted test: %v", response)
	var newMod *module
	if newModHealth > 0 {
		newMod = newModule(p, response, p.language)

		err := p.Route.Endpoint.addModule(newMod, p.slotNum)
		if err != nil {
			return psError(err)
		}

		gm.broadcastState()
		return psSuccess(fmt.Sprintf("Module constructed\nHealth: %d/%d", newMod.Health, newMod.MaxHealth))

	}

	return psError(fmt.Errorf("Module too weak to install: %d/%d", response.passed(), len(response.PassFail)))

}

func cmdRemoveModule(p *Player, args []string, playerCode string) Message {

	if playerCode == "" {
		return psError(errors.New("No code submitted"))
	}

	slot := p.slot()

	if slot == nil {
		return psError(errors.New("Not attached to slot"))
	}

	if !slot.isFull() {
		return psError(errors.New("Slot is empty"))
	}

	// All checks passed:
	// passed error checks on args

	p.outgoing <- psBegin("Removing module...")

	// if module doesn't belong to your team, attack
	if p.Team != slot.module.Team {

		log.Printf("remove cID: %v", slot.challenge.ID)
		response := submitTest(slot.challenge.ID, p.language, playerCode)

		if response.Message.Type == "error" {
			return psError(errors.New(response.Message.Data))
		}

		newModHealth := response.passed()

		log.Printf("response to submitted test: %v", response)

		if newModHealth >= slot.module.Health {

			err := p.Route.Endpoint.removeModule(p.slotNum)
			if err != nil {
				return psError(err)
			}

			gm.broadcastState()
			return psSuccess("Module removed")

		}

		return psError(fmt.Errorf(
			"Solution too weak: %d/%d, need %d/%d to remove",
			response.passed(), len(response.PassFail), slot.module.Health, slot.module.MaxHealth,
		))
	}

	err := p.Route.Endpoint.removeModule(p.slotNum)
	if err != nil {
		return psError(err)
	}

	gm.broadcastState()
	return psSuccess("Module removed")

}

func cmdRefac(p *Player, args []string, playerCode string) Message {

	if playerCode == "" {
		return psError(errors.New("No code submitted"))
	}

	slot := p.slot()

	if slot == nil {
		return psError(errors.New("Not attached to slot"))
	}

	if !slot.isFull() {
		return psError(errors.New("Slot is empty"))
	}

	if slot.module.Health == slot.module.MaxHealth {
		return psError(errors.New("Cannot refactor completed module, you can try to remove (rm)741"))
	}

	// All checks passed:
	// passed error checks on args
	p.outgoing <- psBegin("Refactoring module...")

	response := submitTest(slot.challenge.ID, p.language, playerCode)

	if response.Message.Type == "error" {
		return psError(errors.New(response.Message.Data))
	}

	newModHealth := response.passed()

	log.Printf("response to submitted refactor: %v", response)

	if newModHealth > slot.module.Health {

		// who owns module before refactor:
		log.Printf("refac slot: %v", slot)
		oldTeam := slot.module.Team
		var retMsg Message

		slot.module.Health = newModHealth
		slot.module.Team = p.Team

		// if the module changed hands...
		if oldTeam != p.Team {
			p.Route.Endpoint.evalTrafficForTeam(oldTeam)
			retMsg = psSuccess(fmt.Sprintf("Module refactored to (%v) with health %d/%d", slot.module.Team.Name, slot.module.Health, slot.module.MaxHealth))
		} else {
			retMsg = psSuccess(fmt.Sprintf("Module refactored with health %d/%d", slot.module.Health, slot.module.MaxHealth))
		}

		gm.broadcastState()
		return retMsg
	}

	return psError(fmt.Errorf(
		"Solution too weak: %d/%d, need %d/%d to refactor",
		response.passed(), len(response.PassFail), slot.module.Health+1, slot.module.MaxHealth,
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
