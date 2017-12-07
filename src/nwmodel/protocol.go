package nwmodel

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Handshake related constants
const versionNumber = "0.0.0.1"

// VersionTag is used for handshake and generally to identify the server type and version
const VersionTag = "NodeWars:" + versionNumber

// Global GameModel, not sure how to pass this to handleFunc but that would be better
var gm = NewDefaultModel()

var upgrader = websocket.Upgrader{}

// Ask about reduntant error messaging...
func doHandshake(ws *websocket.Conn) error {

	_, p, err := ws.ReadMessage()
	if err != nil {
		log.Printf("Could not read from socket: %v", err)
		return err
	}

	if string(p) == VersionTag {
		message := []byte("Welcome to NodeWars")
		if err := ws.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("Error while affirming handshake")
			return err
		}
	} else {
		errorMessage := []byte("Version mismation, Server is running " + VersionTag + ", closing connection.")
		if err := ws.WriteMessage(websocket.TextMessage, errorMessage); err != nil {
			log.Println("Error while aborting handshake")
			return err
		}
		return errors.New(string(errorMessage))
	}
	return nil
}

// HandleConnections is the point of entry for all websocket connections
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	log.Println("New player connected...")
	// Upgrade GET to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Close this socket when we're done

	// Attempt handshake
	err = doHandshake(ws)
	if err != nil {
		log.Printf("Handshake error: %v", err)
		ws.Close()
		return
	}

	// Assuming we're all good, register client
	// TODO reconsider this lifecycle, registering player without name has weird side effects
	// log.Println("Registering player...")
	thisPlayer := gm.RegisterPlayer(ws)
	defer scrubPlayerSocket(thisPlayer)

	// Spin up gorouting to monitor outgoing and send those messages to player.Socket
	// log.Println("Spinning up outgoing handler for player...")
	go outgoingRelay(thisPlayer)
	// cannot set language before outgoingRelay is running, will cause program hault
	thisPlayer.setLanguage("python")

	// send initial state
	// log.Println("Sending initial state message to player...")
	thisPlayer.outgoing <- calcStateMsgForPlayer(thisPlayer)
	// Handle socket stream
	for {
		var msg Message

		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}
		// log.Println("received player message")
		incomingHandler(&msg, thisPlayer)
	}
}

// right now this just sends entire game nodewarsmodel for updates
// TODO move to a more imperative style of state updating, necessary?
// this should also take a player as an argument and take into account
// event visibility when composing state message
func calcStateMsgForPlayer(p *Player) Message {
	stateMsg, err := json.Marshal(gm)

	// log.Printf("stateMsg: %v", string(stateMsg))
	if err != nil {
		log.Println(err)
	}

	// log.Println(string(stateMsg))
	return Message{
		Type:   "graphState",
		Sender: "server",
		Data:   string(stateMsg),
	}
}

// Should this do socket scrubbing on error? Is that redundant? TODO
func outgoingRelay(p *Player) {
	for {
		msg := <-p.outgoing
		if err := p.socket.WriteJSON(msg); err != nil {
			log.Printf("error dispatching message to %v", p.Name)
			scrubPlayerSocket(p)
			return
		}
	}
}

// Response are sometimes handled as imperatives, sometimes only effect state and
// are visible after entire stateMessage update. Pick a paradigm TODO
func incomingHandler(msg *Message, p *Player) {
	// Tie message with player name
	msg.Sender = string(p.Name)
	switch msg.Type {

	case "playerCmd":
		res := cmdHandler(msg, p)
		if res.Data != "" {
			p.outgoing <- res
		}

	case "stateRequest":
		p.outgoing <- calcStateMsgForPlayer(p)

	default:
		p.outgoing <- Message{"error", "server", fmt.Sprintf("client sent uknown message type: %v", msg.Type), ""}
	}
}

// func sendWorldState(p *Player)

func scrubPlayerSocket(p *Player) {
	// p.outgoing <- Message{"error", "server", "!!Server Malfunction. Connection Terminated!!")}
	log.Printf("Scrubbing player: %v", p.name())
	gm.RemovePlayer(p)
	p.socket.Close()
}
