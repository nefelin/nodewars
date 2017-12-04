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
	//"tell": cmdTell,
	//"t": cmdTell,
	"tc": cmdTc,

	// // player settings
	"team": cmdTeam,
	"name": cmdName,

	// // world interaction
	"con":     cmdConnect,
	"connect": cmdConnect,

	"mk":      cmdMake,
	"make":    cmdMake,
	"makemod": cmdMake,

	"um":     cmdUnmake,
	"un":     cmdUnmake,
	"unmake": cmdUnmake,

	// debug only
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

func cmdTc(p *Player, args []string, c string) Message {
	if p.Team == nil {
		return msgNoTeam
	}

	chatMsg := p.name() + "> " + strings.Join(args, " ")
	go p.Team.broadcast(Message{
		Type: "(team)",
		Data: chatMsg,
	})
	return Message{}
}

func cmdTeam(p *Player, args []string, c string) Message {
	log.Println("cmdTeam called")
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

func cmdSetPOE(p *Player, args []string, c string) Message {
	if p.Team == nil {
		return msgNoTeam
	}
	newPOE, _ := strconv.Atoi(args[0])
	_ = gm.setPlayerPOE(p, newPOE)

	// debug only :
	gm.POEs[p.ID].addModule(newModuleBy(p.Team))

	gm.broadcastState()
	return Message{}
}

func cmdMake(p *Player, args []string, c string) Message {
	// TODO demand team for make and unmake actions

	if len(args) < 1 {
		return psError(errors.New("Expected 1 argument, received 0"))
	}

	target, err := strconv.Atoi(args[0])
	if err != nil {
		return psError(fmt.Errorf("%v invalid integer", args[0]))
	}

	switch {
	case p.Route == nil:
		return msgNoConnection

	// TODO CHECK ACTUAL SLOT AVAILABILITY

	// TODO maybe switch this to actual slot names with hex ids or something
	case target > p.Route.Endpoint.capacity()-1:
		return psError(errors.New("Slot '" + args[0] + "' does not exist"))

	case p.Route.Endpoint.isFull():
		return psError(errors.New("No slots available"))
	}

	// Success
	p.Route.Endpoint.addModule(newModuleBy(p.Team))
	gm.broadcastState()
	return Message{}
}

func cmdRefac(p *Player, args []string, c string) Message {

	return Message{}
}

func cmdUnmake(p *Player, args []string, c string) Message {

	return Message{}
}
