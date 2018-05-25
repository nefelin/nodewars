package machines

import (
	"challenges"
	"feature"
	"math/rand"
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

	// solution string // store solution used to pass. could be useful for later mechanics
	Type feature.Type `json:"type"` // NA for non-features, none or other feature.Type for features

	Language  string // `json:"languageId"`
	Health    int    `json:"health"`
	MaxHealth int    `json:"maxHealth"`
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

func (m *Machine) IsFeature() bool {
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
	m.Builder = ""
	m.TeamName = ""
	m.Language = ""
	m.Powered = true

	// if m.Type != nil { // reset feature type?
	// 	m.Type = feature.None
	// }

	m.Health = 0
	m.ResetChallenge()
}

func (m *Machine) Claim(lang, Builder, team string, r challenges.GradedResult) {
	m.Language = lang
	m.Builder = Builder
	m.TeamName = team

	// m.Powered = true

	m.Health = r.Passed()
	// m.MaxHealth = len(r.Graded)
}

// dummyClaim is used to claim a Machine for a player without an execution result
func (m *Machine) DummyClaim(teamName string, str string) {
	// m.Builder = p.name
	m.TeamName = teamName
	m.Language = "python" // TODO find ore elegent solution for this
	// m.Powered = true

	switch str {
	case "FULL":
		m.Health = m.MaxHealth
	case "RAND":
		m.Health = rand.Intn(m.MaxHealth) + 1
	case "MIN":
		m.Health = 1
	}
}
