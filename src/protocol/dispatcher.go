package protocol

import (
	"errors"
	"fmt"
	"nwmessage"
	"nwmodel/player"
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
}

func NewDispatcher() *Dispatcher {
	d := &Dispatcher{
		players:           make(map[*websocket.Conn]*player.Player),
		locations:         make(map[*player.Player]room.Room),
		games:             make(map[roomID]room.Room),
		registrationQueue: make(chan regrequest.Request),
		clientMessages:    make(chan nwmessage.ClientMessage),
	}

	go dispatchConsumer(d)
	return d
}

func (d *Dispatcher) Name() string {
	return "Main Lobby"
}

func (d *Dispatcher) Type() string {
	return "Lobby"
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

func (d *Dispatcher) createGame(r room.Room) error {
	if _, ok := d.games[r.Name()]; ok {
		return fmt.Errorf("A game named '%s' already exists", r.Name())
	}

	d.games[r.Name()] = r
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
		return errors.New("Your not in a game")
	}

	game.RemovePlayer(p)
	delete(d.locations, p)
	return nil
}
