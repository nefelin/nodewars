package protocol

import (
	"fmt"
	"nwmessage"
	"nwmodel"
	"regrequest"

	"github.com/gorilla/websocket"
)

type playerID = int
type gameID = string

// Dispatcher ...
type Dispatcher struct {
	players           map[*websocket.Conn]*nwmodel.Player
	locations         map[playerID]gameID
	games             map[gameID]Room
	registrationQueue chan regrequest.Request
	Lobby
}

type Lobby struct {
	players map[playerID]*nwmodel.Player
	aChan   chan nwmessage.Message
}

// Room ...
type Room interface {
	Recv(msg nwmessage.Message)
	AddPlayer(p *nwmodel.Player) error
	RemovePlayer(p *nwmodel.Player) error
	GetPlayers() map[int]*nwmodel.Player
}

func NewDispatcher() *Dispatcher {
	d := &Dispatcher{
		players:           make(map[*websocket.Conn]*nwmodel.Player),
		locations:         make(map[playerID]gameID),
		games:             make(map[gameID]Room),
		registrationQueue: make(chan regrequest.Request),
		Lobby:             NewLobby(),
	}

	go actionConsumer(d)
	return d
}

func NewLobby() Lobby {
	l := Lobby{
		players: make(map[playerID]*nwmodel.Player),
		aChan:   make(chan nwmessage.Message, 100),
	}
	return l
}

// TODO handle errors for add/remove player
func (l *Lobby) AddPlayer(p *nwmodel.Player) error {
	l.players[p.ID] = p
	return nil
}

func (l *Lobby) RemovePlayer(p *nwmodel.Player) error {
	delete(l.players, p.ID)
	return nil
}

func (l *Lobby) GetPlayers() map[int]*nwmodel.Player {
	return l.players
}

func (l *Lobby) Recv(m nwmessage.Message) {
	l.aChan <- m
}

func (d *Dispatcher) Recv(m nwmessage.Message) {
	d.Lobby.Recv(m)
}

func (d *Dispatcher) handleRegRequest(r regrequest.Request) {
	fmt.Printf("Handling RegRequest: %v\n", r.Action)
	switch r.Action {
	case regrequest.Register:
		r.ResChan <- d.registerPlayer(r.Ws)
	case regrequest.Deregister:
		d.deregisterPlayer(r.Ws)
	}
	if r.ResChan != nil {
		close(r.ResChan)
	}
}

func (d *Dispatcher) registerPlayer(ws *websocket.Conn) *nwmodel.Player {

	p := nwmodel.NewPlayer(ws)
	d.Lobby.AddPlayer(p)
	d.players[ws] = p

	return p
}

func (d *Dispatcher) deregisterPlayer(ws *websocket.Conn) {
	p := d.players[ws]

	// fmt.Printf("Removing player id:%d\n", p.ID)

	if gameID, ok := d.locations[p.ID]; ok {
		d.games[gameID].RemovePlayer(p)
		delete(d.locations, p.ID)
	} else {
		// delete(d.Lobby.players, p.ID)
		d.Lobby.RemovePlayer(p)
	}

	p.Socket.Close()
}

// func (d *Dispatcher) scrubPlayerSocket(p *nwmodel.Player) {
// 	// p.outgoing <- Message{"error", "server", "!!Server Malfunction. Connection Terminated!!")}
// 	log.Printf("Scrubbing player: %v\n", p.ID)
// 	// d.removePlayer(p)
// 	// TODO REMOVE THE PLAYER

// 	// if the player is in a game, take him out of the game
// 	if gameID, ok := d.locations[p.ID]; ok {
// 		d.games[gameID].RemovePlayer(p)
// 		delete(d.locations, p.ID)
// 	} else {
// 		// delete(d.Lobby.players, p.ID)
// 		d.Lobby.RemovePlayer(p)
// 	}

// 	p.Socket.Close()
// }

func (d *Dispatcher) makeGame() {}

func (d *Dispatcher) destroyGame() {}

func (d *Dispatcher) joinRoom() {}

func (d *Dispatcher) leaveRoom() {}
