package main

import (
	"fmt"
	"math/rand"

	"github.com/gorilla/websocket"
)

// Team ...
type Team struct {
	Name    string
	Color   string
	Players map[*Player]bool
	Channel chan Message
}

// Player ...
type Player struct {
	Name     string
	Team     *Team
	Socket   *websocket.Conn
	Outgoing chan Message
}

func NewTeam(name, color string) Team {
	return Team{name, color, make(map[*Player]bool), make(chan Message)}
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
	return fmt.Sprintf("<Team> (Name:%v, Color:%v, Players:%v)", t.Name, t.Color, playerList)
}

func (p Player) String() string {
	return fmt.Sprintf("Player Name: %v\nTeam: %v", p.Name, p.Team)
}

// func (p Player)

func makeDummyTeams() []Team {
	var teams []Team
	teams = append(teams, NewTeam("Blue", "blue"))
	teams = append(teams, NewTeam("Red", "red"))

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
