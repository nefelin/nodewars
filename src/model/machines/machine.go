package machines

import (
	"challenges"
	"feature"
	"math/rand"
	"model/player"
	"sync"
)

type Machine struct {
	sync.Mutex
	// accepts   challengeCriteria // store what challenges can fill this Machine
	Challenge challenges.Challenge

	Powered  bool    `json:"powered"`
	Builder  string  // `json:"creator"`
	TeamName string  `json:"owner"`
	CoinVal  float32 `json:"coinval"`

	Address string // mac address in node where Machine resides

	Type feature.Type `json:"type"` // NA for non-features, none or other feature.Type for features

	Solution challenges.Solution
	// Health    int `json:"health"`
	MaxHealth int `json:"maxHealth"`
}

type challengeCriteria struct {
	IDs        []int64  // list of acceptable challenge ids
	Tags       []string // acceptable categories of challenge
	Difficulty [][]int  // acceptable difficulties, [5] = level five, [3,5] = 3,4, or 5
}

// init methods

func NewMachine() *Machine {
	return &Machine{
		Powered: true,
	}
}

func NewFeature() *Machine {
	m := NewMachine()
	m.Type = feature.None
	return m
}

// Machine methods -------------------------------------------------------------------------
func (m *Machine) Health() int {
	if !m.Powered {
		return 1
	}

	return m.Solution.Strength
}

// resetChallenge should use m.accepts to get a challenge matching criteria TODO
func (m *Machine) ResetChallenge() {
	m.Challenge = challenges.GetRandomChallenge()
	m.MaxHealth = len(m.Challenge.Cases)
}

func (m *Machine) IsNeutral() bool {
	if m.TeamName == "" {
		return true
	}
	return false
}

func (m *Machine) IsGateway() bool {
	// fmt.Printf("Machine Type: %v", m.Type)
	// fmt.Printf("Feature NA: %v", feature.NA)
	// fmt.Printf("Equal: %v", m.Type == feature.NA)
	if m.Type == nil {
		return false
	}
	return true
}

func (m *Machine) BelongsTo(teamName string) bool {
	if m.TeamName == teamName {
		return true
	}
	return false
}

func (m *Machine) Reset() {
	m.TeamName = ""
	m.Powered = true

	m.Solution = challenges.Solution{}

	// if m.Type != nil { // reset feature type?
	// 	m.Type = feature.None
	// }

	m.ResetChallenge()
}

func (m *Machine) Claim(team string, s challenges.Solution) {
	m.TeamName = team
	m.Solution = s
}

// dummyClaim is used to claim a Machine for a player without an execution result
func (m *Machine) DummyClaim(teamName string, str string) {
	// m.Builder = p.name
	m.TeamName = teamName

	switch str {
	case "FULL":
		m.Solution.Strength = m.MaxHealth
	case "RAND":
		m.Solution.Strength = rand.Intn(m.MaxHealth) + 1
	case "MIN":
		m.Solution.Strength = 1
	}
}

func (m *Machine) AcceptsLanguageFrom(p *player.Player, lang string) bool {
	if m.IsNeutral() {
		return true
	}
	if m.BelongsTo(p.TeamName) {
		return true
	}
	if m.Solution.IsDummy {
		return true
	}
	if m.Solution.Language == lang {
		return true
	}
	return false
}
