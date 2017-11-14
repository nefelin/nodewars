package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"math/rand"
)

// Team ...
type Team struct {
	Name    string // Names are only colors for now
	Players map[*Player]bool
	Open    bool
}

// Player ...
type Player struct {
	Name     string
	Team     *Team
	Socket   *websocket.Conn
	Outgoing chan Message
}

// NewTeam creates a new team with color/name color
func NewTeam(color string) Team {
	return Team{color, make(map[*Player]bool), true}
}

func (t *Team) broadcast(msg Message) {
	fmt.Printf("message, %v\n", msg)
	for player := range t.Players {
		player.Outgoing <- msg
	}
}

// Does it matter if player joins team or team adds player?
func (t *Team) addPlayer(p *Player) {
	t.Players[p] = true
	p.Team = t

	// Tell player they've joined
	p.Outgoing <- Message{
		Type:   "teamAssign",
		Sender: "server",
		Data:   t.Name,
	}
}

func (t *Team) removePlayer(p *Player) {
	delete(t.Players, p)
	p.Team = nil
}

func (t Team) String() string {
	var playerList []string
	for player := range t.Players {
		playerList = append(playerList, player.Name)
	}
	return fmt.Sprintf("<Team> (Name: %v, Players:%v)", t.Name, playerList)
}

func (p Player) String() string {
	return fmt.Sprintf("Player Name: %v\nTeam: %v", p.Name, p.Team)
}

// func (p Player)

func makeDummyTeams() map[string]Team {
	teams := make(map[string]Team)
	teams["red"] = NewTeam("red")
	teams["blue"] = NewTeam("blue")

	return teams
}

func registerPlayer(ws *websocket.Conn) *Player {
	newPlayer := Player{randStringBytes(5), nil, ws, make(chan Message)}
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
