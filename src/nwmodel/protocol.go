package nwmodel

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

// Handshake related constants
const versionNumber = "0.0.0.1"

// VersionTag is used for handshake and generally to identify the server type and version
const VersionTag = "NodeWars:" + versionNumber

// Global GameModel, not sure how to pass this to handleFunc but that would be better
var gm = NewDefaultModel()

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
	thisPlayer := gm.RegisterPlayer(ws)
	defer scrubPlayerSocket(thisPlayer)

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

// right now this just sends entire game nodewarsmodel for updates
// TODO move to a more imperative style of state updating, necessary?
// this should also take a player as an argument and take into account
// event visibility when composing state message
func calcStateMsgForPlayer() Message {
	gs := newGameState()

	// log.Println(gs)

	// for name, teamOb := range gm.Teams {
	// 	// make a blank roster of length players
	// 	gs.teamsOfPlayers[name] = make([]string, len(teamOb.Players))
	// 	i := 0
	// 	for p := range teamOb.Players {
	// 		gs.teamsOfPlayers[name][i] = p.Name
	// 		i++
	// 	}
	// }

	stateMsg, err := json.Marshal(gs)

	// log.Printf("stateMsg: %v", string(stateMsg))
	if err != nil {
		log.Println(err)
	}

	// log.Println(string(stateMsg))
	return Message{
		Type:   "gameState",
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
	msg.Sender = p.Name
	switch msg.Type {
	case "allChat":
		if p.Name == "" {
			p.outgoing <- Message{"error", "server", "you need a name to send messages"}
		} else {
			for player := range gm.Players {
				player.outgoing <- Message{"allChat", p.Name, msg.Data}
			}
		}

	case "teamChat":
		//HANDLE chat by unassigned player, maybe make an Observer team by default?
		if p.Team != nil {
			go p.Team.broadcast(*msg)
		} else {
			p.outgoing <- Message{"error", "server", "unable to teamChat without team assignment"}
		}

		// Attach sendersocket's name since its relevant for chats
	case "teamJoin":
		if p.Name == "" {
			p.outgoing <- Message{"error", "server", "you need a name to join a team"}
		} else {
			err := gm.assignPlayerToTeam(p, msg.Data)
			if err != nil {
				p.outgoing <- Message{"error", "server", fmt.Sprintln(err)}
			}
		}

		// team method handles messaging, fix TODO

	case "stateRequest":
		p.outgoing <- calcStateMsgForPlayer()

	case "setPOE":
		if p.Team == nil {
			p.outgoing <- Message{"error", "server", "you need a team to interact with the map"}
		} else {
			newPOE, err := strconv.Atoi(msg.Data)
			if err != nil {
				log.Printf("setPOE error: %v", err)
			} else {
				if res := gm.setPlayerPOE(p, newPOE); res {
					p.outgoing <- Message{"POEset", "server", msg.Data}
					// p.outgoing <- calcStateMsgForPlayer()
				} else {
					p.outgoing <- Message{"error", "server", "failed to set, '" + msg.Data + "', as POE. Either does not exist or player cannot switch POE"}
				}
			}
		}

	case "nodeConnect":
		if p.Team == nil {
			// if player has no team yet, complain
			p.outgoing <- Message{"error", "server", "You need a team to interact with the map"}
		} else {
			targetNode, err := strconv.Atoi(msg.Data)
			if err != nil {
				// if we have trouble converting msg to integer, complain
				log.Printf("connectToNode error: %v", err)
			} else {
				// if we're all good, try to connect the player to the node
				if targetNode > -1 && targetNode < nodeCount {
					if gm.tryConnectPlayerToNode(p, targetNode) {
						p.outgoing <- Message{"connectSuccess", "pseudoServer", msg.Data}
					} else {
						p.outgoing <- Message{"connectFail", "pseudoServer", msg.Data}
					}
				} else {
					p.outgoing <- Message{"error", "server", "node '" + msg.Data + "' does not exist"}
				}
			}
		}

	case "setPlayerName":
		err := gm.setPlayerName(p, msg.Data)
		if err != nil {
			p.outgoing <- Message{"error", "server", fmt.Sprintln(err)}
		} else {
			p.outgoing <- Message{"playerNameSet", "server", msg.Data}
		}

	default:
		p.outgoing <- Message{"error", "server", fmt.Sprintf("client sent uknown message type: %v", msg.Type)}
	}
}

// func sendWorldState(p *Player)

func scrubPlayerSocket(p *Player) {
	log.Printf("Scrubbing player: %v", p.Name)
	gm.RemovePlayer(p)
	p.socket.Close()
}
