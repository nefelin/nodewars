package protocol

import (
	"errors"
	"fmt"
	"nwmodel"
	"regrequest"

	"github.com/gorilla/websocket"
)

type playerID = int
type roomID = string

// Dispatcher ...
type Dispatcher struct {
	players           map[*websocket.Conn]*nwmodel.Player
	locations         map[*nwmodel.Player]Room
	games             map[roomID]Room
	registrationQueue chan regrequest.Request
	clientMessages    chan nwmodel.ClientMessage
}

func NewDispatcher() *Dispatcher {
	d := &Dispatcher{
		players:           make(map[*websocket.Conn]*nwmodel.Player),
		locations:         make(map[*nwmodel.Player]Room),
		games:             make(map[roomID]Room),
		registrationQueue: make(chan regrequest.Request),
		clientMessages:    make(chan nwmodel.ClientMessage),
	}

	go dispatchConsumer(d)
	return d
}

type Room interface {
	Name() string
	Type() string // TODO switch this to a roomtype definition
	Recv(msg nwmodel.ClientMessage) error
	AddPlayer(p *nwmodel.Player) error
	RemovePlayer(p *nwmodel.Player) error
	GetPlayers() []*nwmodel.Player
}

func (d *Dispatcher) Name() string {
	return "Main Lobby"
}

func (d *Dispatcher) Type() string {
	return "Lobby"
}

func (d *Dispatcher) GetPlayers() []*nwmodel.Player {
	list := make([]*nwmodel.Player, len(d.players))
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

func (d *Dispatcher) AddPlayer(p *nwmodel.Player) error {
	d.players[p.Socket] = p
	return nil
}

func (d *Dispatcher) RemovePlayer(p *nwmodel.Player) error {

	if game, ok := d.locations[p]; ok {
		game.RemovePlayer(p)
		delete(d.locations, p)
	}

	delete(d.players, p.Socket)

	p.Socket.Close()
	close(p.Outgoing)

	return nil
}

func (d *Dispatcher) createGame(r Room) error {
	if _, ok := d.games[r.Name()]; ok {
		return fmt.Errorf("A game named '%s' already exists", r.Name())
	}

	d.games[r.Name()] = r
	return nil
}

func (d *Dispatcher) destroyGame() {}

// manipulating/examining player objects

func (d *Dispatcher) joinRoom(p *nwmodel.Player, r roomID) error {
	if d.locations[p] != nil {
		return errors.New("Can't join game, already in a game")
	}

	d.locations[p] = d.games[r]
	d.games[r].AddPlayer(p)
	return nil
}

func (d *Dispatcher) leaveRoom(p *nwmodel.Player) error {
	game, ok := d.locations[p]
	if !ok {
		return errors.New("Your not in a game")
	}

	game.RemovePlayer(p)
	delete(d.locations, p)
	return nil
}
