package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
)

// Handshake related constants
const versionNumber = "0.0.0.1"
const versionTag = "NodeWars:" + versionNumber

var players = make(map[*websocket.Conn]*Player) // connected players
var broadcast = make(chan Message)
var teams = makeDummyTeams()
var gameMap = NewDefaultMap()

var upgrader = websocket.Upgrader{}

// Message is our basic message struct
type Message struct {
	Type   string `json:"type"`
	Sender string `json:"sender"`
	Data   string `json:"data"`
}

// Ask about reduntant error messaging...
func doHandshake(ws *websocket.Conn) error {

	_, p, err := ws.ReadMessage()
	if err != nil {
		log.Printf("Could not read from socket: %v", err)
		return err
	}

	if string(p) == versionTag {
		message := []byte("Welcome to NodeWars")
		if err := ws.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("Error while affirming handshake")
			return err
		}
	} else {
		errorMessage := []byte("Version mismation, Server is running " + versionTag + ", closing connection.")
		if err := ws.WriteMessage(websocket.TextMessage, errorMessage); err != nil {
			log.Println("Error while aborting handshake")
			return err
		}
		return errors.New(string(errorMessage))
	}
	return nil
}

func handleConnections(w http.ResponseWriter, r *http.Request) {

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
	thisPlayer := registerPlayer(ws)
	defer scrubPlayerSocket(ws)

	// Spin up gorouting to monitor outgoing and send those messages to player.Socket
	go outgoingRelay(thisPlayer)

	// Handle socket stream
	for {
		var msg Message

		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}
		incomingHandler(&msg, thisPlayer)
	}
}

func calcStateMsgForPlayer() Message {
	currentState := GameState{*gameMap, make([]GameEvent, 0)}
	stateMsg, _ := json.Marshal(currentState)

	return Message{
		Type:   "gameState",
		Sender: "server",
		Data:   string(stateMsg), // string(stateMsg),
	}
}

// Should this do socket scrubbing on error? Is that redundant? TODO
func outgoingRelay(p *Player) {
	for {
		msg := <-p.Outgoing
		if err := p.Socket.WriteJSON(msg); err != nil {
			log.Printf("error dispatching message to %v", p.Name)
			return
		}
	}
}

func incomingHandler(msg *Message, p *Player) {
	// Tie message with player name
	msg.Sender = p.Name
	switch msg.Type {
	case "allChat":
		for _, player := range players {
			player.Outgoing <- Message{"allChat", p.Name, msg.Data}
		}

	case "teamChat":
		//HANDLE chat by unassigned player, maybe make an Observer team by default?
		if p.Team != nil {
			go p.Team.broadcast(*msg)
		} else {
			p.Outgoing <- Message{"error", "server", "unable to teamChat without team assignment"}
		}

		// Attach sendersocket's name since its relevant for chats
	case "teamJoin":
		if team, ok := teams[msg.Data]; ok {
			p.joinTeam(&team)
		} else {
			p.Outgoing <- Message{"error", "server", "team '" + msg.Data + "' does not exist"}
		}

	case "stateRequest":
		p.Outgoing <- calcStateMsgForPlayer()

	default:
		p.Outgoing <- Message{"error", "server", "uknown message type"}
	}
}

// func sendWorldState(p *Player)

func scrubPlayerSocket(ws *websocket.Conn) {
	players[ws].Team.removePlayer(players[ws])
	delete(players, ws)
	ws.Close()
}

func main() {

	// So it doesn't complain about fmt
	fmt.Println("Starting " + versionTag + " server...")

	// Set up log file
	f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	//Close log file when we're done
	defer f.Close()
	//set output of logs to f
	log.SetOutput(f)

	// Start Webserver
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", handleConnections)

	log.Println("Starting server on port 8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
