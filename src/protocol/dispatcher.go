package protocol

import (
	"command"
	"docs"
	"errors"
	"fmt"
	"help"
	"model"
	"model/player"
	"nwmessage"
	"regrequest"
	"room"

	"github.com/gorilla/websocket"
)

type playerID = int
type roomID = string

// Dispatcher ...
type Dispatcher struct {
	players           map[*websocket.Conn]*player.Player
	locations         map[*player.Player]room.Room
	games             map[roomID]room.Room
	registrationQueue chan regrequest.Request
	clientMessages    chan nwmessage.ClientMessage
	helpRegistry      *help.Registry
	cmdRegistry       *command.Registry
}

func NewDispatcher() *Dispatcher {
	d := &Dispatcher{
		players:           make(map[*websocket.Conn]*player.Player),
		locations:         make(map[*player.Player]room.Room),
		games:             make(map[roomID]room.Room),
		registrationQueue: make(chan regrequest.Request),
		clientMessages:    make(chan nwmessage.ClientMessage),
	}

	d.helpRegistry = help.NewRegistry()                 // make new help
	d.cmdRegistry = command.NewRegistry(d.helpRegistry) // make new command collection (all commands added will have their help info added to d.helpRegistry)
	docs.RegisterTopics(d.helpRegistry)                 // register general topics with this helpRegistry

	RegisterCommands(d.cmdRegistry, d)
	model.RegisterCommands(d.cmdRegistry)

	go dispatchConsumer(d)
	return d
}

func (d *Dispatcher) Name() string {
	return "Main Lobby"
}

func (d *Dispatcher) Type() room.Type {
	return room.Lobby
}

func (d *Dispatcher) GetPlayers() []*player.Player {
	list := make([]*player.Player, len(d.players))
	var i int
	for _, p := range d.players {
		list[i] = p
		i++
	}
	return list
}

func (d *Dispatcher) handleRegRequest(r regrequest.Request) {
	act := "Register"
	if r.Action == regrequest.Deregister {
		act = "Deregister"
	}
	fmt.Printf("Handling RegRequest: %s\n", act)
	switch r.Action {
	case regrequest.Register:
		d.AddPlayer(r.Player)
		close(r.ResChan)
	case regrequest.Deregister:
		d.RemovePlayer(r.Player)
	}
}

func (d *Dispatcher) AddPlayer(p *player.Player) error {
	d.players[p.Socket()] = p
	status := fmt.Sprintf("%d players online", len(d.players))
	p.Outgoing(nwmessage.PsNeutral(fmt.Sprintf("%s%s\n%s", logo, welcomeMsg, status)))
	p.Outgoing(nwmessage.PsPrompt(p.Prompt()))
	return nil
}

func (d *Dispatcher) RemovePlayer(p *player.Player) error {

	if game, ok := d.locations[p]; ok {
		game.RemovePlayer(p)
		delete(d.locations, p)
	}

	delete(d.players, p.Socket())

	p.Cleanup()

	return nil
}

func (d *Dispatcher) createGame(name string, r func() (*model.GameModel, error)) error {
	if _, ok := d.games[name]; ok {
		return fmt.Errorf("A game named '%s' already exists", name)
	}

	game, err := r()
	if err != nil {
		return err
	}

	d.games[name] = game
	return nil
}

func (d *Dispatcher) destroyGame() {}

// manipulating/examining player objects

func (d *Dispatcher) joinRoom(p *player.Player, r roomID) error {
	if d.locations[p] != nil {
		return errors.New("Can't join game, already in a game")
	}

	d.locations[p] = d.games[r]
	d.games[r].AddPlayer(p)
	return nil
}

func (d *Dispatcher) leaveRoom(p *player.Player) error {
	game, ok := d.locations[p]
	if !ok {
		return errors.New("You're not in a game")
	}

	game.RemovePlayer(p)
	delete(d.locations, p)
	return nil
}
