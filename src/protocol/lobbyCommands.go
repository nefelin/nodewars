package protocol

// import (
// 	"commands"
// 	"fmt"
// 	"nwmessage"
// 	"nwmodel"
// 	"sort"
// 	"strings"
// )

// // Recv ...
// func (d *Dispatcher) Recv(m nwmodel.ClientMessage) error {
// 	if m.Data == "" {
// 		return nil
// 	}

// 	return commandList.Exec(d, m)
// }

import (
	"argument"
	"commands"
	"errors"
	"fmt"
	"math/rand"
	"nwmessage"
	"nwmodel"
	"sort"
	"strconv"
	"strings"
)

var commandList = commands.CommandGroup{
	"chat": {
		Name:      "chat",
		ShortDesc: "Toggles chat mode (all text entered is broadcast)",
		ArgsReq:   argument.ArgList{},
		ArgsOpt:   argument.ArgList{},
		Handler:   cmdToggleChat,
	},

	"join": {

		Name:      "join",
		ShortDesc: "Joins the specified game",
		ArgsReq: argument.ArgList{
			{Name: "game_name", Type: argument.String},
		},
		ArgsOpt: argument.ArgList{},
		Handler: cmdJoinGame,
	},

	"leave": {

		Name:      "leave",
		ShortDesc: "Leaves the current game",
		ArgsReq:   argument.ArgList{},
		ArgsOpt:   argument.ArgList{},
		Handler:   cmdLeaveGame,
	},

	"ls": {

		Name:      "ls",
		ShortDesc: "List the games that are currently running",
		ArgsReq:   argument.ArgList{},
		ArgsOpt:   argument.ArgList{},
		Handler:   cmdListGames,
	},

	"name": {

		Name:      "name",
		ShortDesc: "Sets the player's name",
		ArgsReq: argument.ArgList{
			{Name: "new_name", Type: argument.String},
		},
		ArgsOpt: argument.ArgList{},
		Handler: cmdSetName,
	},

	"ng": {

		Name:      "ng",
		ShortDesc: "Creates a new game",
		ArgsReq:   argument.ArgList{},
		ArgsOpt: argument.ArgList{
			{Name: "game_name", Type: argument.String},
		},
		Handler: cmdNewGame,
	},

	"kill": {

		Name:      "kill",
		ShortDesc: "Removes a game (must be empty)",
		ArgsReq: argument.ArgList{
			{Name: "game_name", Type: argument.String},
		},
		ArgsOpt: argument.ArgList{},
		Handler: cmdKillGame,
	},

	"tell": {

		Name:      "tell",
		ShortDesc: "Sends a private message to another player",
		ArgsReq: argument.ArgList{
			{Name: "recip", Type: argument.String},
			{Name: "msg", Type: argument.GreedyString},
		},
		ArgsOpt: argument.ArgList{},
		Handler: cmdTell,
	},

	"who": {

		Name:      "who",
		ShortDesc: "Shows who's in the lobby",
		ArgsReq:   argument.ArgList{},
		ArgsOpt:   argument.ArgList{},
		Handler:   cmdWho,
	},

	"yell": {

		Name:      "yell",
		ShortDesc: "Sends a message to all player (in the same game/lobby)",
		ArgsReq: argument.ArgList{
			{Name: "msg", Type: argument.GreedyString},
		},
		ArgsOpt: argument.ArgList{},
		Handler: cmdYell,
	},
}

func cmdToggleChat(p *nwmodel.Player, context interface{}, args []interface{}) error {
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

func cmdYell(p *nwmodel.Player, context interface{}, args []interface{}) error {
	d := context.(*Dispatcher)
	msg := args[0].(string)

	for _, player := range d.GetPlayers() {
		player.Outgoing <- nwmessage.PsChat(p.GetName(), "global", msg)
	}
	return nil
}

func cmdSetName(p *nwmodel.Player, context interface{}, args []interface{}) error {
	d := context.(*Dispatcher)
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

func cmdNewGame(p *nwmodel.Player, context interface{}, args []interface{}) error {
	d := context.(*Dispatcher)

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

func cmdKillGame(p *nwmodel.Player, context interface{}, args []interface{}) error {
	d := context.(*Dispatcher)

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

func cmdWho(p *nwmodel.Player, context interface{}, args []interface{}) error {
	d := context.(*Dispatcher)
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

func cmdJoinGame(p *nwmodel.Player, context interface{}, args []interface{}) error {
	d := context.(*Dispatcher)
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

func cmdLeaveGame(p *nwmodel.Player, context interface{}, args []interface{}) error {
	d := context.(*Dispatcher)

	err := d.leaveRoom(p)

	if err != nil {
		return err
	}

	p.Outgoing <- nwmessage.PsSuccess(fmt.Sprintf("You have left the game"))
	return nil
}

func cmdTell(p *nwmodel.Player, context interface{}, args []interface{}) error {
	d := context.(*Dispatcher)

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

func cmdListGames(p *nwmodel.Player, context interface{}, args []interface{}) error {
	d := context.(*Dispatcher)
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

// func (cg *CommandGroup) Exec(d *Dispatcher, m nwmodel.ClientMessage) error {
// 	fullCmd := strings.Split(m.Data, " ")

// 	// handle help
// 	if fullCmd[0] == "help" {
// 		if len(fullCmd) == 1 {

// 			m.Sender.Outgoing <- nwmessage.PsNeutral(cg.AllHelp())

// 		} else {
// 			help, err := cg.Help(fullCmd[1:])

// 			if err != nil {
// 				m.Sender.Outgoing <- nwmessage.PsError(err)
// 			}
// 			m.Sender.Outgoing <- nwmessage.PsNeutral(help)

// 		}
// 		return nil

// 	}

// 	// if we find the command, try to execute
// 	if cmd, ok := commandList[fullCmd[0]]; ok {

// 		args, err := cmd.ValidateArgs(fullCmd[1:])
// 		if err != nil {
// 			// if we have trouble validating args
// 			m.Sender.Outgoing <- nwmessage.PsError(fmt.Errorf("%s\nusage: %s", err.Error(), cmd.Usage()))
// 		} else {
// 			// otherwise actually execute the command
// 			err = cmd.handler(m.Sender, d, args)
// 			if err != nil {
// 				m.Sender.Outgoing <- nwmessage.PsError(err)
// 			}
// 		}

// 		return nil
// 	}

// 	// if we don't find the command, pass an error back to caller in case caller wants to do something else
// 	return unknownCommand(fullCmd[0])
// }

// // Help composes a help string for the given command
// func (cg CommandGroup) Help(args []string) (string, error) {
// 	if cmd, ok := cg[args[0]]; ok {
// 		return cmd.LongHelp(), nil
// 	}
// 	return "", unknownCommand(args[0])
// }

// // AllHelp composes help for all commands in the group
// func (cg CommandGroup) AllHelp() string {
// 	cmds := make([]string, len(cg))
// 	var i int
// 	for key := range cg {
// 		cmds[i] = key
// 		i++
// 	}

// 	sort.Strings(cmds)
// 	// offset := cg.longestKey()
// 	helpStr := make([]string, len(cmds)+1)
// 	helpStr[0] = "Available commands:"

// 	for i, cmd := range cmds {
// 		helpStr[i+1] = cg[cmd].ShortHelp()
// 	}

// 	return strings.Join(helpStr, "\n")
// }

// func unknownCommand(cmd string) error {
// 	return fmt.Errorf("Unknown command, '%s'", cmd)
// }
