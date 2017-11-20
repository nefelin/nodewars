package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/gorilla/websocket"
)

type colorName = string

// Team ...
type Team struct {
	Name    colorName        // Names are only colors for now
	Players map[*Player]bool //THESE create circular JSON problem, decide on org scheme and fix, TODO
	MaxSize int
}

// Player ...
type Player struct {
	Name         string
	Team         *Team //THESE create circular JSON problem, decide on org scheme and fix, TODO
	PointOfEntry nodeID
	Socket       *websocket.Conn
	Outgoing     chan Message
}

// NewTeam creates a new team with color/name color
func NewTeam(n colorName) Team {
	return Team{n, make(map[*Player]bool), 2}
}

func (t Team) isFull() bool {
	if len(t.Players) < t.MaxSize {
		return false
	}
	return true
}

func (t *Team) broadcast(msg Message) {
	for player := range t.Players {
		player.Outgoing <- msg
	}
}

// Team-received Methods
func (t *Team) addPlayer(p *Player) {
	t.Players[p] = true
	p.Team = t

	// Tell client they've joined
	p.Outgoing <- Message{
		Type:   "teamAssign",
		Sender: "server",
		Data:   t.Name,
	}
}

func (t *Team) removePlayer(p *Player) {
	delete(t.Players, p)
	p.Team = nil
	// Notify client
	p.Outgoing <- Message{
		Type:   "teamUnassign",
		Sender: "server",
		Data:   t.Name,
	}
}

func (t Team) String() string {
	var playerList []string
	for player := range t.Players {
		playerList = append(playerList, player.Name)
	}
	return fmt.Sprintf("<Team> (Name: %v, Players:%v)", t.Name, playerList)
}

// Player-received methods

func (p *Player) joinTeam(t *Team) {
	if p.Team == nil {
		if !t.isFull() {
			t.addPlayer(p)
		} else {
			// tell player team is full, TODO centralize control messages
			p.Outgoing <- Message{"teamFull", "server", t.Name}
		}
	} else {
		p.Outgoing <- Message{"error", "server", "you are already a member of " + p.Team.Name}
	}
}

// right now player can only connect to adjacent nodes
func (p *Player) connectToNode(n nodeID) bool {
	log.Printf("Player %v attempting to connect to node %v from POE %v", p.Name, n, p.PointOfEntry)
	if gameMap.routeToNode(gameMap.Nodes[p.PointOfEntry], gameMap.Nodes[n]) != nil {
		return true
	}
	return false
}

func (p Player) String() string {
	return fmt.Sprintf("<Player> Name: %v, Team: %v", p.Name, p.Team)
}

func makeDummyTeams() map[string]Team {
	teams := make(map[colorName]Team)
	teams["red"] = NewTeam("red")
	teams["blue"] = NewTeam("blue")

	return teams
}

func registerPlayer(ws *websocket.Conn) *Player {
	newPlayer := Player{randStringBytes(5), nil, -1, ws, make(chan Message)}
	players[ws] = &newPlayer
	return &newPlayer
}

// For names for now :/
// const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const letterBytes = "adammarrygeorgejohnjeffstacylou"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
