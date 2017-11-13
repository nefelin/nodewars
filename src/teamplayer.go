package main

import (
	"fmt"
	"github.com/gorilla/websocket"
)

// Team ...
type Team struct {
	Name    string
	Color   string
	Players map[*Player]bool
	Channel chan string
}

func (t Team) broadcast(m Message) {
	for player := range t.Players {
		fmt.Println(player.Name)
	}
}

// Does it matter if player joins team or team adds player?
func (t Team) addPlayer(p *Player) {
	t.Players[p] = true
	p.Team = &t
}

func (t Team) String() string {
	// TODO implement playlist
	return fmt.Sprintf("Team name: %v, Color: %v", t.Name, t.Color)
}

// Player ...
type Player struct {
	Name   string
	Team   *Team
	Socket *websocket.Conn
}

func (p Player) String() string {
	return fmt.Sprintf("Player Name: %v\nTeam: %v", p.Name, p.Team)
}

// func (p Player)

func makeDummyTeams() []Team {
	var teams []Team
	teams = append(teams, Team{"Blue", "blue", nil, make(chan string)})
	teams = append(teams, Team{"Red", "red", nil, make(chan string)})

	return teams
}

func registerPlayer(ws *websocket.Conn) *Player {
	newPlayer := Player{"", nil, ws}
	players[ws] = &newPlayer
	return &newPlayer
}
