package protocol

import (
	"log"
	"nwmessage"
	"nwmodel"

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
	aChan   chan nwmessage.Message
}

// Room ...
type Room interface {
	Recv(msg nwmessage.Message)
	AddPlayer(p *nwmodel.Player) error
	RemovePlayer(p *nwmodel.Player) error
	// getStateForP(pID int)
}

func NewDispatcher() *Dispatcher {
	d := &Dispatcher{
		locations: make(map[playerID]gameID),
		games:     make(map[gameID]Room),
		Lobby:     NewLobby(),
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

func (d *Dispatcher) Recv(m nwmessage.Message) {
	// pID, _ := strconv.Atoi(m.Sender)
	// gameID, ok := d.locations[pID]

	// if !ok || gameID == "lobby" {
	d.Lobby.aChan <- m
	// } else {

	// }
}

// func (l *Lobby) recv(m nwmessage.Message) {
// 	l.aChan <- m
// }

// func actionConsumer(l Lobby) {
// 	for {
// 		msg := <-l.aChan
// 		id, _ := strconv.Atoi(msg.Sender)

// 		msg.Type = "allChat"
// 		msg.Sender = "pseudoServer"
// 		msg.Data = "blue"
// 		// gm := nwmodel.NewDefaultModel()
// 		// state := gm.CalcState(nil)

// 		l.players[id].Outgoing <- msg
// 		l.players[id].Outgoing <- nwmessage.PromptStateMsg("lobby>")
// 		// nwmodel.Message{
// 		// 	Type:   "graphState",
// 		// 	Sender: "server",
// 		// 	Data:   state,
// 		// }
// 	}
// }

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

	// if the player is in a game, take him out of the game
	if gameID, ok := d.locations[p.ID]; ok {
		d.games[gameID].RemovePlayer(p)
		delete(d.locations, p.ID)
	}

	p.Socket.Close()
}

func (d *Dispatcher) makeGame() {}

func (d *Dispatcher) destroyGame() {}

func (d *Dispatcher) joinRoom() {}

func (d *Dispatcher) leaveRoom() {}
