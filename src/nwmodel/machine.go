package nwmodel

import (
	"feature"
	"fmt"
	"math/rand"
	"sync"

	"nwmessage"
)

type machine struct {
	mu sync.RWMutex

	// accepts   challengeCriteria // store what challenges can fill this machine
	challenge Challenge

	Powered  bool   `json:"powered"`
	builder  string // `json:"creator"`
	TeamName string `json:"owner"`

	attachedPlayers map[*Player]bool

	address string // mac address in node where machine resides

	// solution string // store solution used to pass. could be useful for later mechanics
	Type feature.Type `json:"type"` // NA for non-features, none or other feature.Type for features

	language  string // `json:"languageId"`
	Health    int    `json:"health"`
	MaxHealth int    `json:"maxHealth"`
}

type challengeCriteria struct {
	IDs        []int64  // list of acceptable challenge ids
	Tags       []string // acceptable categories of challenge
	Difficulty [][]int  // acceptable difficulties, [5] = level five, [3,5] = 3,4, or 5
}

// init methods

func newMachine() *machine {
	return &machine{
		Powered:         true,
		attachedPlayers: make(map[*Player]bool),
	}
}

func newFeature() *machine {
	m := newMachine()
	m.setType(feature.None)
	return m
}

// machine methods -------------------------------------------------------------------------

// QUESTION: is technically reading state but an awkward place to employ a read lock
func (m *machine) detachAll(msg string) {
	for p := range m.attachedPlayers {
		p.macDetach()
		if msg != "" {
			p.Outgoing <- nwmessage.PsAlert(msg)
		}
	}
}

// mutex wrapper for easier debugging
func (m *machine) lock() {
	// fmt.Println("lock machine")
	m.mu.Lock()
}
func (m *machine) unlock() {
	// fmt.Println("unlock machine")
	m.mu.Unlock()
}
func (m *machine) rlock() {
	// fmt.Println("rlock machine")
	m.mu.RLock()
}
func (m *machine) runlock() {
	// fmt.Println("runlock machine")
	m.mu.RUnlock()
}

// Getters (need rlock)

func (m *machine) getType() feature.Type {
	m.rlock()
	defer m.runlock()
	return m.Type
}

func (m *machine) isNeutral() bool {
	m.rlock()
	defer m.runlock()

	if m.TeamName == "" {
		return true
	}
	return false
}

func (m *machine) isFeature() bool {
	m.rlock()
	defer m.runlock()

	if m.Type == nil {
		return false
	}
	return true
}

func (m *machine) belongsTo(teamName string) bool {
	m.rlock()
	defer m.runlock()

	if m.TeamName == teamName {
		return true
	}
	return false
}

// Setters (require Locks) -----------------------------------------------

func (m *machine) setType(t feature.Type) {
	m.lock()
	defer m.unlock()
	m.Type = t
}

func (m *machine) addPlayer(p *Player) {
	m.lock()
	defer m.unlock()

	m.attachedPlayers[p] = true
}

func (m *machine) remPlayer(p *Player) {
	m.lock()
	defer m.unlock()

	delete(m.attachedPlayers, p)

}

func (m *machine) reset() {
	m.lock()
	defer m.unlock()

	m.builder = ""
	m.TeamName = ""
	m.language = ""
	m.Powered = true

	m.detachAll(fmt.Sprintf("mac:%s is resetting, you have been detached", m.address))

	// if m.Type != nil { // reset feature type?
	// 	m.Type = feature.None
	// }

	m.Health = 0
	m.resetChallenge()
}

// resetChallenge should use m.accepts to get a challenge matching criteria TODO
func (m *machine) resetChallenge() {
	m.lock()
	defer m.unlock()

	m.challenge = getRandomChallenge()
	m.MaxHealth = len(m.challenge.Cases)

}

func (m *machine) claim(p *Player, r GradedResult) {
	m.lock()
	defer m.unlock()

	m.builder = p.name
	m.TeamName = p.TeamName
	m.language = p.language
	// m.Powered = true

	m.Health = r.passed()

}

// dummyClaim is used to claim a machine for a player without an execution result
func (m *machine) dummyClaim(teamName string, str string) {
	m.lock()
	defer m.unlock()

	// m.builder = p.name
	m.TeamName = teamName
	m.language = "python" // TODO find ore elegent solution for this
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
