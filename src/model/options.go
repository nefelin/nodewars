package model

import (
	"challenges"
	"model/node"
)

type gameOptions struct {
	name        string
	timelimit   int
	languages   map[string]challenges.Language
	minPlayers  int
	maxPlayers  int
	langLock    bool   // do players need to use same language as enemy module to solve
	autoAssign  bool   // assign players to a team automatically
	password    string // none means game is open
	mapSize     int
	mapGen      func(int) (*node.Map, error)
	coinGoal    float32
	defaultLang string
	// teamCount int // how many teams
	// FFA bool // each player is their own team, deal with colors :/
}

func newDefaultOptions() gameOptions {
	return gameOptions{
		name:        "",
		timelimit:   0,
		languages:   challenges.GetLanguages(),
		minPlayers:  0,
		maxPlayers:  0,
		langLock:    true,
		autoAssign:  false,
		password:    "",
		mapSize:     12,
		mapGen:      node.CutTestMap,
		coinGoal:    10000,
		defaultLang: "",
	}
}

// new model options

// func CoinGoal(goal int) func(*GameModel) error {
// 	return func(m *GameModel) error {
// 		return m.setCoinGoal(goal)
// 	}
// }

// func Languages(langs ...challenge.Language) func(*GameModel) error {
// 	return func(m *GameModel) error {
// 		return m.setLanguages(langs)
// 	}
// }

// func TimeLimit(limit int) func(*GameModel) error {
// 	return func(m *GameModel) error {
// 		return m.setTimeLimit(limit)
// 	}
// }

// func Name(name string) func(*GameModel) error {
// 	return func(m *GameModel) error {
// 		m.options.name = name
// 		return nil
// 	}
// }

// func MapMode(m *GameModel) error {
// 	return m.setMode(modes.Map)
// }
