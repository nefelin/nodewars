package nwmodel

import (
	"fmt"
	"math/rand"
)

type nodeID = int

type node struct {
	ID          nodeID     `json:"id"` // keys and ids is redundant? TODO
	Connections []nodeID   `json:"connections"`
	Machines    []*machine `json:"machines"` // TODO why is this a list of pointerS?
	Feature     *machine   `json:"feature"`
	Remoteness  float64    //`json:"remoteness"`
	addressMap  map[string]*machine
}

// node methods -------------------------------------------------------------------------------

func (n *node) claimFreeMachine(p *Player) error {
	neutral := make([]int, 0)

	for i := range n.Machines {
		if n.Machines[i].TeamName == "" {
			neutral = append(neutral, i)
		}
	}

	if len(neutral) < 1 {
		return fmt.Errorf("Node %d contains no neutral machines to claim", n.ID)
	}

	target := neutral[rand.Intn(len(neutral))]

	n.Machines[target].dummyClaim(p.TeamName, "FULL")
	return nil

}

// coinVal calculates the coin produced per machine in a given node
func (n *node) coinVal(t teamName) float32 {
	return float32(n.Remoteness)
}

// coinProduction gives the total produced (for team t) in a given node
func (n *node) coinProduction(t teamName) float32 {
	var total float32
	coinPerMac := n.coinVal(t)
	for _, mac := range n.Machines {
		if mac.TeamName == t && mac.Powered {
			total += coinPerMac
		}
	}
	return total
}

func (n *node) initMachines() {
	n.Machines = make([]*machine, len(n.Connections))
	for i := range n.Connections {
		n.Machines[i] = newMachine()
		n.Machines[i].resetChallenge()
	}

	// n.Feature = newFeature(n.Feature.Type) // Preserve type in case map contains feature type info
	n.Feature.resetChallenge()

	n.initAddressMap()
}

func (n *node) initAddressMap() {
	featureAddress := newMacAddress(2)
	n.setMacAddress(featureAddress, n.Feature)

	for _, m := range n.Machines {

		newAddress := newMacAddress(2)
		_, ok := n.addressMap[newAddress]
		for ok {
			newAddress = newMacAddress(2)
			_, ok = n.addressMap[newAddress]
		}

		n.setMacAddress(newAddress, m)
	}

}

func (n *node) setMacAddress(address string, mac *machine) {
	// TODO error check name collisions
	mac.address = address
	n.addressMap[address] = mac

}

// addConnection is reciprocol
func (n *node) addConnection(m *node) {
	// if the connection already exists, ignore
	for _, nID := range n.Connections {
		if m.ID == nID {
			return
		}
	}

	if m.ID == n.ID {
		return
	}

	n.Connections = append(n.Connections, m.ID)
	m.Connections = append(m.Connections, n.ID)
}

func (n *node) remConnection(ni nodeID) {
	n.Connections = cutIntFromSlice(ni, n.Connections)
}

func cutIntFromSlice(p int, s []int) []int {
	for i, thisP := range s {
		if p == thisP {
			// swaps the last element with the found element and returns with the last element cut
			s[len(s)-1], s[i] = s[i], s[len(s)-1]
			return s[:len(s)-1]
		}
	}
	return s
}

func (n *node) hasMachineFor(t *team) bool {
	// t == nil means we don't care... used in calculating node eccentricity without rewriting dijkstras
	if t == nil {
		return true
	}

	// if a node has no machines, it allows routing for everyone
	// allows creation of neutral hubs
	if len(n.Machines) == 0 {
		return true
	}

	// if we control a powered machine here we can route through
	for _, mac := range n.Machines {
		if mac.TeamName != "" {
			if mac.TeamName == t.Name { // && mac.Powered {
				return true
			}
		}
	}

	// or if we control a powered feature here we can route through
	if n.Feature.TeamName == t.Name { // && n.Feature.Powered {
		return true
	}

	return false
}

func (n *node) machinesFor(t teamName) int {
	var count int

	for _, mac := range n.Machines {
		if mac.TeamName == t {
			count++
		}
	}

	if n.Feature.TeamName == t {
		count++
	}

	return count
}

func (n *node) supportsRouting(t teamName) bool {
	// t == nil means we don't care... used in calculating node eccentricity without rewriting dijkstras
	if t == "" {
		return true
	}

	// if a node has no machines, it allows routing for everyone
	// allows creation of neutral hubs
	if len(n.Machines) == 0 {
		return true
	}

	if n.machinesFor(t) < 1 {
		return false
	}

	return true
}

func (n *node) powerMachines(t teamName, onOff bool) {
	for _, mac := range n.Machines {
		if mac.TeamName == t {
			mac.Powered = onOff
		}
	}
}

// n.resetMachine should never be called directly. only from gm.removeModule
// func (n *node) resetMachine(slotIndex int) error {
// 	if slotIndex < 0 || slotIndex > len(n.Machines)-1 {
// 		return errors.New("No valid attachment")
// 	}

// 	machine := n.Machines[slotIndex]

// 	if machine.TeamName == "" {
// 		return errors.New("Machine is alread neutral")
// 	}

// 	// reset machine
// 	machine.reset()

// 	return nil
// }

// helper to generate macAddresses
const charBytes = "abcdefghijklmnopqrstuvwxyz0123456789"

func newMacAddress(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charBytes[rand.Intn(len(charBytes))]
	}
	return string(b)
}

// helper function for removing player from slice of players
func cutStrFromSlice(p string, s []string) []string {
	for i, thisP := range s {
		if p == thisP {
			// swaps the last element with the found element and returns with the last element cut
			s[len(s)-1], s[i] = s[i], s[len(s)-1]
			return s[:len(s)-1]
		}
	}
	// log.Printf("CutPlayer returning: %v", s)
	// log.Println("Player not found in slice")
	return s
}
