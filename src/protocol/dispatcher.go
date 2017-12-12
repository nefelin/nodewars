package protocol

import (
	"log"
	"nwmodel"
	"strconv"

	"github.com/gorilla/websocket"
)

type playerID = int
type gameID = string

// Dispatcher ...
type Dispatcher struct {
	locations map[playerID]gameID
	games     map[gameID]Room
	Lobby
}

type Lobby struct {
	players map[playerID]*nwmodel.Player
	aChan   chan nwmodel.Message
}

// Room ...
type Room interface {
	recv(msg nwmodel.Message)
	// getStateForP(pID int)
}

func NewDispatcher() *Dispatcher {
	d := &Dispatcher{
		locations: make(map[playerID]gameID),
		games:     make(map[gameID]Room),
		Lobby:     NewLobby(),
	}

	return d
}

func NewLobby() Lobby {
	l := Lobby{
		players: make(map[playerID]*nwmodel.Player),
		aChan:   make(chan nwmodel.Message, 100),
	}
	go actionConsumer(l)
	return l
}

func (d *Dispatcher) recv(m nwmodel.Message) {
	pID, _ := strconv.Atoi(m.Sender)
	gameID, ok := d.locations[pID]

	if !ok || gameID == "lobby" {
		d.Lobby.recv(m)
	} else {
		d.games[gameID].recv(m)
	}
}

func (l *Lobby) recv(m nwmodel.Message) {
	l.aChan <- m
}

func actionConsumer(l Lobby) {
	for {
		msg := <-l.aChan
		id, _ := strconv.Atoi(msg.Sender)

		msg.Type = "alertFlash"
		msg.Sender = "server"
		msg.Data = "blue"
		// gm := nwmodel.NewDefaultModel()
		// state := gm.CalcState(nil)

		l.players[id].Outgoing <- msg
		// nwmodel.Message{
		// 	Type:   "graphState",
		// 	Sender: "server",
		// 	Data:   state,
		// }
	}
}

// func (d *Dispatcher) getRoom(pID int) Room {
// 	return d.games[d.players[pID].ID]
// }

func (d *Dispatcher) registerPlayer(ws *websocket.Conn) *nwmodel.Player {
	p := nwmodel.NewPlayer(ws)
	d.Lobby.players[p.ID] = p

	return p
}

func (d *Dispatcher) removePlayer() {}

func (d *Dispatcher) scrubPlayerSocket(p *nwmodel.Player) {
	// p.outgoing <- Message{"error", "server", "!!Server Malfunction. Connection Terminated!!")}
	log.Printf("Scrubbing player: %v", p.ID)
	// d.removePlayer(p)
	// TODO REMOVE THE PLAYER
	p.Socket.Close()
}

func (d *Dispatcher) makeGame() {}

func (d *Dispatcher) destroyGame() {}

func (d *Dispatcher) joinRoom() {}

func (d *Dispatcher) leaveRoom() {}
