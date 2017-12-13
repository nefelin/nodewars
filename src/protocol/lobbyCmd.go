package protocol

import (
	"fmt"
	"log"
	"nwmessage"
	"nwmodel"
	"strconv"
	"strings"
)

// type playerCommand func(p *Player, gm *GameModel, args []string, code string) Message
type lobbyCmd func(*nwmodel.Player, *Dispatcher, []string) nwmessage.Message

var msgMap = map[string]lobbyCmd{
	"new":     cmdNewGame,
	"newgame": cmdNewGame,
	"join":    cmdJoinGame,
}

func actionConsumer(d *Dispatcher) {
	for {
		m := <-d.Lobby.aChan
		pID, err := strconv.Atoi(m.Sender)

		if err != nil {
			log.Println(err)
		}

		p := d.Lobby.players[pID]

		msg := strings.Split(m.Data, " ")

		// log.Println("Finding handlerFunc...")
		if handlerFunc, ok := msgMap[msg[0]]; ok {
			// log.Println("Calling handlerFunc")
			res := handlerFunc(p, d, msg[1:])
			if res.Data != "" {
				p.Outgoing <- res
			}
		} else if gameID, ok := d.locations[pID]; ok {
			d.games[gameID].Recv(m)
		} else {
			p.Outgoing <- nwmessage.PsUnknown("lobby: " + msg[0])
		}
	}
}

func cmdTest(p *nwmodel.Player, d *Dispatcher, args []string) nwmessage.Message {
	return nwmessage.Message{}
}

func cmdNewGame(p *nwmodel.Player, d *Dispatcher, args []string) nwmessage.Message {

	if _, ok := d.games[args[0]]; ok {
		return nwmessage.PsError(fmt.Errorf("A game named '%s' already exists", args[0]))
	}

	if gameID, ok := d.locations[p.ID]; ok {
		return nwmessage.PsError(fmt.Errorf("You can't create a game. You're already playing in '%s'", gameID))
	}

	// create the game
	newGame := nwmodel.NewDefaultModel()

	// add the creator to the game
	newGame.AddPlayer(p)

	// register new game with dispatch
	d.games[args[0]] = newGame

	// have the dispatcher assign the player to newGame
	d.locations[p.ID] = args[0]

	return nwmessage.PsSuccess(fmt.Sprintf("New game, '%s', created", args[0]))
}

func cmdJoinGame(p *nwmodel.Player, d *Dispatcher, args []string) nwmessage.Message {
	if gameID, ok := d.locations[p.ID]; ok {
		return nwmessage.PsError(fmt.Errorf("You're already playing in '%s'", gameID))
	}

	if _, ok := d.games[args[0]]; !ok {
		return nwmessage.PsError(fmt.Errorf("No game named '%s' exists", args[0]))
	}

	// add the creator to the game
	d.games[args[0]].AddPlayer(p)

	// have the dispatcher assign the player to newGame
	d.locations[p.ID] = args[0]

	return nwmessage.PsSuccess(fmt.Sprintf("Joined '%s'", args[0]))
}
