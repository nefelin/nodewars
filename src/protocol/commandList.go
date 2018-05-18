package protocol

import (
	"argument"
	"commandinfo"
	"errors"
	"fmt"
	"math/rand"
	"nwmessage"
	"nwmodel"
	"sort"
	"strconv"
	"strings"
)

// var commandList = map[string]lobbyCommand{
var commandList = LobbyCommandGroup{
	"chat": {
		Info: commandinfo.Info{
			Name:      "chat",
			ShortDesc: "Toggles chat mode (all text entered is broadcast)",
			ArgsReq:   argument.ArgList{},
			ArgsOpt:   argument.ArgList{},
		},
		handler: cmdToggleChat,
	},

	"join": {
		Info: commandinfo.Info{
			Name:      "join",
			ShortDesc: "Joins the specified game",
			ArgsReq: argument.ArgList{
				{Name: "game_name", Type: argument.String},
			},
			ArgsOpt: argument.ArgList{},
		},
		handler: cmdJoinGame,
	},

	"leave": {
		Info: commandinfo.Info{
			Name:      "leave",
			ShortDesc: "Leaves the current game",
			ArgsReq:   argument.ArgList{},
			ArgsOpt:   argument.ArgList{},
		},
		handler: cmdLeaveGame,
	},

	"ls": {
		Info: commandinfo.Info{
			Name:      "ls",
			ShortDesc: "List the games that are currently running",
			ArgsReq:   argument.ArgList{},
			ArgsOpt:   argument.ArgList{},
		},
		handler: cmdListGames,
	},

	"name": {
		Info: commandinfo.Info{
			Name:      "name",
			ShortDesc: "Sets the player's name",
			ArgsReq: argument.ArgList{
				{Name: "new_name", Type: argument.String},
			},
			ArgsOpt: argument.ArgList{},
		},
		handler: cmdSetName,
	},

	"ng": {
		Info: commandinfo.Info{
			Name:      "ng",
			ShortDesc: "Creates a new game",
			ArgsReq:   argument.ArgList{},
			ArgsOpt: argument.ArgList{
				{Name: "game_name", Type: argument.String},
			},
		},
		handler: cmdNewGame,
	},

	"kill": {
		Info: commandinfo.Info{
			Name:      "kill",
			ShortDesc: "Removes a game (must be empty)",
			ArgsReq: argument.ArgList{
				{Name: "game_name", Type: argument.String},
			},
			ArgsOpt: argument.ArgList{},
		},
		handler: cmdKillGame,
	},

	"tell": {
		Info: commandinfo.Info{
			Name:      "tell",
			ShortDesc: "Sends a private message to another player",
			ArgsReq: argument.ArgList{
				{Name: "recip", Type: argument.String},
				{Name: "msg", Type: argument.GreedyString},
			},
			ArgsOpt: argument.ArgList{},
		},
		handler: cmdTell,
	},

	"who": {
		Info: commandinfo.Info{
			Name:      "who",
			ShortDesc: "Shows who's in the lobby",
			ArgsReq:   argument.ArgList{},
			ArgsOpt:   argument.ArgList{},
		},
		handler: cmdWho,
	},

	"yell": {
		Info: commandinfo.Info{
			Name:      "yell",
			ShortDesc: "Sends a message to all player (in the same game/lobby)",
			ArgsReq: argument.ArgList{
				{Name: "msg", Type: argument.GreedyString},
			},
			ArgsOpt: argument.ArgList{},
		},
		handler: cmdYell,
	},
}

// func assertStringSlice(f []interface{}) []string {
// 	ret := make([]string, len(f))

// 	for i, v := range f {
// 		ret[i] = v.(string)
// 	}

// 	return ret
// }

func cmdToggleChat(p *nwmodel.Player, d *Dispatcher, args []interface{}) error {
	p.ChatMode = !p.ChatMode

	var flag string
	if p.ChatMode {
		flag = "ON"
	} else {
		flag = "OFF"
	}

	p.Outgoing <- nwmessage.PsNeutral(fmt.Sprintf("ChatMode set to %s", flag))
	return nil
}

func cmdYell(p *nwmodel.Player, d *Dispatcher, args []interface{}) error {
	msg := args[0].(string)

	for _, player := range d.GetPlayers() {
		player.Outgoing <- nwmessage.PsChat(p.GetName(), "global", msg)
	}
	return nil
}

func cmdSetName(p *nwmodel.Player, d *Dispatcher, args []interface{}) error {
	if d.locations[p] != nil {
		return errors.New("Can only change name while in Lobby")
	}
	name := args[0].(string)

	if name == p.GetName() {
		return fmt.Errorf("Your name's already set to '%s'", name)
	}

	// check for name collision
	for _, player := range d.players {
		if name == player.GetName() {
			return fmt.Errorf("The name '%s' is already taken", name)
		}
	}

	p.SetName(name)
	return nil
}

func cmdNewGame(p *nwmodel.Player, d *Dispatcher, args []interface{}) error {

	if _, ok := d.locations[p]; ok {
		return fmt.Errorf("You can't create a game. You're already in a game")
	}

	var gameName string
	if len(args) == 0 {
		gameName = "Game_" + strconv.Itoa(rand.Intn(100000))
		_, exists := d.games[gameName]
		for exists {
			gameName = "Game_" + strconv.Itoa(rand.Intn(100000))
		}
	} else {
		gameName = args[0].(string)
	}

	// create the game
	err := d.createGame(nwmodel.NewDefaultModel(gameName))

	if err != nil {
		return err
	}

	err = d.joinRoom(p, gameName)
	if err != nil {
		return err
	}

	// p.SendPrompt()
	p.Outgoing <- nwmessage.PsSuccess(fmt.Sprintf("New game, '%s', created and joined", gameName))
	return nil
}

func cmdKillGame(p *nwmodel.Player, d *Dispatcher, args []interface{}) error {

	gameName := args[0].(string)

	game, ok := d.games[gameName]

	if !ok {
		return fmt.Errorf("No game named, '%s', exists", gameName)
	}

	// check to make sure game is empty
	if len(game.GetPlayers()) != 0 {
		return fmt.Errorf("The game, '%s', is not empty", gameName)
	}

	// clean up game
	// TODO is this sufficient?
	delete(d.games, gameName)

	p.Outgoing <- nwmessage.PsSuccess(fmt.Sprintf("The game, '%s', has been removed", gameName))
	return nil
}

func cmdWho(p *nwmodel.Player, d *Dispatcher, args []interface{}) error {
	// var location Room
	// location, ok := d.games[d.locations[p.ID]]

	// if !ok {
	// 	location = d.Lobby
	// }
	location := d.locations[p]
	if location == nil {
		location = d
	}

	if len(location.GetPlayers()) == 0 {
		p.Outgoing <- nwmessage.PsNeutral("There are no players here")
		return nil
	}

	var playerNames sort.StringSlice

	for _, p := range location.GetPlayers() {
		playerNames = append(playerNames, p.GetName())
	}

	playerNames.Sort()

	retMsg := "Players here:\n" + strings.Join(playerNames, ", ")
	p.Outgoing <- nwmessage.PsNeutral(retMsg)
	return nil
}

func cmdJoinGame(p *nwmodel.Player, d *Dispatcher, args []interface{}) error {
	gameName := args[0].(string)

	_, ok := d.games[gameName]
	if !ok {
		return fmt.Errorf("No game named '%s' exists", gameName)
	}

	err := d.joinRoom(p, gameName)
	if err != nil {
		return err
	}

	p.Outgoing <- nwmessage.PsSuccess(fmt.Sprintf("Joined game, '%s'", gameName))
	return nil
}

func cmdLeaveGame(p *nwmodel.Player, d *Dispatcher, args []interface{}) error {

	err := d.leaveRoom(p)

	if err != nil {
		return err
	}

	p.Outgoing <- nwmessage.PsSuccess(fmt.Sprintf("You have left the game"))
	return nil
}

func cmdTell(p *nwmodel.Player, d *Dispatcher, args []interface{}) error {

	name := args[0].(string)
	msg := args[1].(string)
	var recip *nwmodel.Player

	// check lobby for recipient:
	for _, player := range d.players {
		if player.GetName() == name {
			recip = player
		}
	}

	if recip == nil {
		return fmt.Errorf("No such player, '%s'", name)
	}

	recip.Outgoing <- nwmessage.PsChat(p.GetName(), "private", msg)
	p.Outgoing <- nwmessage.PsNeutral(fmt.Sprintf("(you to %s): %s", recip.GetName(), msg))
	return nil
}

func cmdListGames(p *nwmodel.Player, d *Dispatcher, args []interface{}) error {
	gameList := ""

	if len(d.games) == 0 {
		p.Outgoing <- nwmessage.PsNeutral("No games running. Type, 'new game_name', to start one")
		return nil
	}

	for gameName, game := range d.games {
		gameList += fmt.Sprintf("'%s' - Players: %d\n", gameName, len(game.GetPlayers()))
	}

	p.Outgoing <- nwmessage.PsNeutral(strings.TrimSpace("Available games:\n" + gameList))
	return nil
}
