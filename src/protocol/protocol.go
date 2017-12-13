package protocol

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"nwmessage"
	"nwmodel"
	"strconv"

	"github.com/gorilla/websocket"
)

// Handshake related constants
const versionNumber = "0.0.0.1"

// VersionTag is used for handshake and generally to identify the server type and version
const VersionTag = "NodeWars:" + versionNumber

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
func HandleConnections(w http.ResponseWriter, r *http.Request, d *Dispatcher) {
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
	thisPlayer := d.registerPlayer(ws)
	defer d.scrubPlayerSocket(thisPlayer)

	// Spin up gorouting to monitor outgoing and send those messages to player.Socket
	// log.Println("Spinning up outgoing handler for player...")
	go outgoingRelay(thisPlayer)
	thisPlayer.Outgoing <- nwmessage.PromptState("lobby>")

	// Handle socket stream
	for {
		var msg nwmessage.Message

		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}
		// log.Println("received player message")
		incomingHandler(d, msg, thisPlayer)
	}
}

// Should this do socket scrubbing on error? Is that redundant? TODO
func outgoingRelay(p *nwmodel.Player) {
	for {
		msg := <-p.Outgoing
		if err := p.Socket.WriteJSON(msg); err != nil {
			log.Printf("error dispatching message to %v", p.Name)
			// scrubPlayerSocket(p)
			return
		}
	}
}

// Response are sometimes handled as imperatives, sometimes only effect state and
// are visible after entire stateMessage update. Pick a paradigm TODO
func incomingHandler(d *Dispatcher, msg nwmessage.Message, p *nwmodel.Player) {
	// Tie message with player name
	msg.Sender = strconv.Itoa(p.ID)
	switch msg.Type {

	case "playerCmd":
		d.Recv(msg)
		// d.getRoom(p.ID).recv(msg)
		// res := cmdHandler(msg, p)
		// if res.Data != "" {
		// 	p.outgoing <- res
		// }

	default:
		p.Outgoing <- nwmessage.Message{"error", "server", fmt.Sprintf("client sent uknown message type: %v", msg.Type), ""}
	}
}
