package protocol

import (
	"errors"
	"fmt"
	"log"
	"nwmessage"
	"nwmodel"
	"sort"
	"strconv"
	"strings"
)

// type playerCommand func(p *Player, gm *GameModel, args []string, code string) Message
type playerCmd func(*nwmodel.Player, *Dispatcher, []string) nwmessage.Message

var lobbyCmdList = map[string]playerCmd{
	// TODO leaveGame should demand confirmation
	"leave": cmdLeaveGame,

	"t":    cmdTell,
	"tell": cmdTell,

	"ls":   cmdListGames,
	"list": cmdListGames,

	"join": cmdJoinGame,

	"name": cmdSetName,

	"new": cmdNewGame,

	"rm": cmdKillGame,

	"who": cmdWho,
}

var globalCmdList = map[string]bool{
	// TODO leaveGame should demand confirmation
	"leave": true,
	"t":     true,
	"tell":  true,
}

// func actionConsumer(d *Dispatcher) {
// 	for {

// 		select {
// 		// if we get a new player, register and pass back
// 		// to the connection handler
// 		case regReq := <-d.registrationQueue:
// 			regReq.retChan <- d.registerPlayer(regReq.ws)

// 		// if we get a player command, handle that
// 		case m := <-d.Lobby.aChan:
// 			pID, err := strconv.Atoi(m.Sender)

// 			if err != nil {
// 				log.Println(err)
// 			}

// 			p := d.Lobby.players[pID]

// 			log.Printf("This is odd, what is p here: %v", p)

// 			msg := strings.Split(m.Data, " ")

// 			// log.Println("recvd messg")
// 			if handlerFunc, ok := globalCmdList[msg[0]]; ok {
// 				// if the player's in a
// 				gameName, ok := d.locations[pID]
// 				if ok {
// 					// if the players in a game we should grab the player object from the game...
// 					p = d.games[gameName].GetPlayers()[pID]
// 				}

// 				res := handlerFunc(p, d, msg[1:])
// 				if res.Data != "" {
// 					p.Outgoing <- res
// 				}
// 				p.Outgoing <- nwmessage.PromptState(p.GetName() + "@(lobby)>")

// 			} else if gameName, ok := d.locations[pID]; ok {
// 				// p = d.games[gameName].GetPlayers()[p.ID]
// 				d.games[gameName].Recv(m)
// 			} else if handlerFunc, ok := lobbyCmdList[msg[0]]; ok {
// 				res := handlerFunc(p, d, msg[1:])
// 				if res.Data != "" {
// 					p.Outgoing <- res
// 				}
// 				p.Outgoing <- nwmessage.PromptState(p.GetName() + "@(lobby)>")
// 			} else {
// 				// if it's not a known lobby command and the player
// 				// isn't in a game, treat it as a chat.
// 				chatMsg := fmt.Sprintf("%s: %s", p.GetName(), strings.Join(msg, " "))

// 				// log.Printf("lobbCmd d.Lobby.GetPlayers: %v", d.Lobby.GetPlayers())
// 				for _, player := range d.Lobby.GetPlayers() {
// 					player.Outgoing <- nwmessage.PsChat(chatMsg, "(lobby)")
// 				}
// 				p.Outgoing <- nwmessage.PromptState(p.GetName() + "@(lobby)>")
// 			}
// 		}
// 	}
// }

func actionConsumer(d *Dispatcher) {
	for {
		select {
		// if we get a new player, register and pass back
		// to the connection handler
		case regReq := <-d.registrationQueue:
			regReq.retChan <- d.registerPlayer(regReq.ws)

		// if we get a player command, handle that
		case m := <-d.Lobby.aChan:

			pID, err := strconv.Atoi(m.Sender)

			if err != nil {
				log.Println(err)
			}

			msg := strings.Split(m.Data, " ")

			gameName, inGame := d.locations[pID]

			// if players not in a game (i.e in lobby)
			if !inGame {
				// PLAYER IN LOBBY

				// get player object
				p := d.Lobby.players[pID]

				// if its a valid lobby command execute
				if handlerFunc, ok := lobbyCmdList[msg[0]]; ok {
					res := handlerFunc(p, d, msg[1:])
					if res.Data != "" {
						p.Outgoing <- res
					}

				} else {
					// if it's not a known lobby command and the player
					// isn't in a game, treat it as a chat.
					chatMsg := strings.Join(msg, " ")

					for _, player := range d.Lobby.GetPlayers() {
						player.Outgoing <- nwmessage.PsChat(p.GetName(), "global", chatMsg)
					}
				}
				p.Outgoing <- nwmessage.PromptState(p.GetName() + "@(lobby)>")
			} else {
				// PLAYER IN GAME

				// get a valid player object
				p := d.games[gameName].GetPlayers()[pID]

				// is the command entered a global command?
				if _, ok := globalCmdList[msg[0]]; ok {
					// get the handler function
					handlerFunc := lobbyCmdList[msg[0]]
					res := handlerFunc(p, d, msg[1:])
					if res.Data != "" {
						p.Outgoing <- res
					}
				} else {
					// it's not a global command, let the game handle it
					d.games[gameName].Recv(m)
				}
			}
		}
	}
}

func cmdSetName(p *nwmodel.Player, d *Dispatcher, args []string) nwmessage.Message {
	if len(args) < 1 {
		return nwmessage.PsError(errors.New("Expected 1 argument, received none"))
	}

	// check lobby for name collision
	for _, player := range d.Lobby.players {
		if args[0] == player.GetName() {
			return nwmessage.PsError(fmt.Errorf("The name '%s' is already taken", args[0]))
		}
	}

	// check ongoing games for name collision
	for _, gm := range d.games {
		for _, player := range gm.GetPlayers() {
			if args[0] == player.GetName() {
				return nwmessage.PsError(fmt.Errorf("The name '%s' is already taken", args[0]))
			}
		}
	}

	p.SetName(args[0])
	return nwmessage.PsSuccess("Name set to '" + p.GetName() + "'")
}

func cmdNewGame(p *nwmodel.Player, d *Dispatcher, args []string) nwmessage.Message {

	if len(args) == 0 || args[0] == "" {
		return nwmessage.PsError(errors.New("Need a name for new game"))
	}

	if _, ok := d.games[args[0]]; ok {
		return nwmessage.PsError(fmt.Errorf("A game named '%s' already exists", args[0]))
	}

	if gameName, ok := d.locations[p.ID]; ok {
		return nwmessage.PsError(fmt.Errorf("You can't create a game. You're already playing in '%s'", gameName))
	}

	// create the game
	newGame := nwmodel.NewDefaultModel()

	// register new game with dispatch
	d.games[args[0]] = newGame

	// tell dispatcher about change of locations
	d.locations[p.ID] = args[0]
	// take player out of lobby
	d.Lobby.RemovePlayer(p)
	// put player in the game
	newGame.AddPlayer(p)

	return nwmessage.PsSuccess(fmt.Sprintf("New game, '%s', created and joined", args[0]))
}

func cmdKillGame(p *nwmodel.Player, d *Dispatcher, args []string) nwmessage.Message {

	if len(args) == 0 || args[0] == "" {
		return nwmessage.PsError(errors.New("Need a name for game to remove"))
	}

	game, ok := d.games[args[0]]

	if !ok {
		return nwmessage.PsError(fmt.Errorf("No game, '%s', exists", args[0]))
	}

	// check to make sure game is empty
	if len(game.GetPlayers()) != 0 {
		return nwmessage.PsError(fmt.Errorf("The game, '%s', is not empty", args[0]))
	}

	// clean up game
	// TODO is this sufficient?
	delete(d.games, args[0])

	return nwmessage.PsSuccess(fmt.Sprintf("The game, '%s', has been removed", args[0]))
}

func cmdWho(p *nwmodel.Player, d *Dispatcher, args []string) nwmessage.Message {
	// var location Room
	// location, ok := d.games[d.locations[p.ID]]

	// if !ok {
	// 	location = d.Lobby
	// }
	location := d.Lobby

	if len(location.GetPlayers()) == 0 {
		return nwmessage.PsNeutral("There are no players here")
	}

	var playerNames sort.StringSlice

	for _, p := range location.GetPlayers() {
		playerNames = append(playerNames, p.GetName())
	}

	playerNames.Sort()

	retMsg := "Players here:\n" + strings.Join(playerNames, ", ")
	return nwmessage.PsNeutral(retMsg)
}

func cmdJoinGame(p *nwmodel.Player, d *Dispatcher, args []string) nwmessage.Message {
	if len(args) < 1 {
		return nwmessage.PsError(errors.New("Expected 1 argument, received none"))
	}

	gameName, ok := d.locations[p.ID]
	if ok {
		return nwmessage.PsError(fmt.Errorf("You're already playing in '%s'", gameName))
	}

	game, ok := d.games[args[0]]
	if !ok {
		return nwmessage.PsError(fmt.Errorf("No game named '%s' exists", args[0]))
	}

	// tell dispatcher about change of locations
	d.locations[p.ID] = args[0]
	// take player out of lobby
	d.Lobby.RemovePlayer(p)
	// put player in the game
	game.AddPlayer(p)

	return nwmessage.PsSuccess(fmt.Sprintf("Joined game, '%s'", args[0]))
}

func cmdLeaveGame(p *nwmodel.Player, d *Dispatcher, args []string) nwmessage.Message {

	gameName, ok := d.locations[p.ID]
	if !ok {
		return nwmessage.PsError(errors.New("Your not in a game"))
	}

	// add the creator to the game
	d.games[gameName].RemovePlayer(p)
	// have the dispatcher assign the player to newGame
	delete(d.locations, p.ID)
	// put the player back in the lobby
	d.Lobby.AddPlayer(p)

	return nwmessage.PsSuccess(fmt.Sprintf("Left game, '%s'", gameName))
}

func cmdTell(p *nwmodel.Player, d *Dispatcher, args []string) nwmessage.Message {

	if len(args) < 2 {
		return nwmessage.PsError(errors.New("Need a recipient and a message"))
	}

	var recip *nwmodel.Player

	// check lobby for recipient:
	for _, player := range d.Lobby.players {
		if player.GetName() == args[0] {
			recip = player
		}
	}

	// if not in lobby check all games
	if recip == nil {
		for _, game := range d.games {
			for _, player := range game.GetPlayers() {
				if player.GetName() == args[0] {
					recip = player
				}
			}
		}
	}

	if recip == nil {
		return nwmessage.PsError(fmt.Errorf("No such player, '%s'", args[0]))
	}

	chatMsg := strings.Join(args[1:], " ")

	recip.Outgoing <- nwmessage.PsChat(p.GetName(), "private", chatMsg)
	return nwmessage.PsNeutral(fmt.Sprintf("(you to %s): %s", recip.GetName(), chatMsg))
}

func cmdListGames(p *nwmodel.Player, d *Dispatcher, args []string) nwmessage.Message {
	gameList := ""

	if len(d.games) == 0 {
		return nwmessage.PsNeutral("No games running. Type, 'new game_name', to start one")
	}

	for gameName, game := range d.games {
		gameList += fmt.Sprintf("'%s' - Players: %d\n", gameName, len(game.GetPlayers()))
	}

	return nwmessage.PsNeutral(strings.TrimSpace("Available games:\n" + gameList))
}
