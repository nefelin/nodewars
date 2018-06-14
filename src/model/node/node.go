package node

import (
	"feature"
	"fmt"
	"math/rand"
	"model/machines"
)

type teamName = string
type NodeID = int
type MacAddress = string

type Node struct {
	ID          NodeID              `json:"id"` // keys and ids is redundant? TODO
	Connections []NodeID            `json:"connections"`
	Machines    []*machines.Machine `json:"machines"` // TODO why is this a list of pointerS?
	Feature     *machines.Machine   `json:"feature"`
	Remoteness  float64             //`json:"remoteness"`
	addressMap  map[MacAddress]*machines.Machine
}

// Node methods -------------------------------------------------------------------------------

// coinVal calculates the coin produced per machine in a given Node
func (n *Node) coinVal(t teamName) float32 {
	base := float32(n.Remoteness)

	if n.DominatedBy(t) {
		base = base * 2
	}

	if n.Feature.TeamName == t && n.Feature.Type == feature.Overclock {
		base = base * 2
	}

	return base
}

func (n *Node) DominatedBy(t teamName) bool { // does t control all non feature machines?
	for _, m := range n.Machines {
		if m.TeamName != t {
			return false
		}
	}
	return true
}

// coinProduction gives the total produced (for team t) in a given Node
func (n *Node) CoinProduction(t teamName) float32 {
	fmt.Println("<coinProduction>")
	var total float32
	coinPerMac := n.coinVal(t)
	for _, mac := range n.Machines {
		if mac.TeamName == t {
			mac.CoinVal = coinPerMac // bind value to machine for rendering
			fmt.Printf("<coinProduction> Setting machine CoinVal: %f\n", mac.CoinVal)
			if mac.Powered { // if the machine's powered, add to teams production
				total += coinPerMac
			}
		}
	}
	return total
}

func (n *Node) createMachinePerEdge() {
	n.Machines = make([]*machines.Machine, len(n.Connections))
	for i := range n.Connections {
		n.Machines[i] = machines.NewMachine()
	}
}

func (n *Node) createMachines(count int) {
	n.Machines = make([]*machines.Machine, count)
	for i := range n.Machines {
		n.Machines[i] = machines.NewMachine()
	}
}

func (n *Node) initMachines() {
	for _, m := range n.Machines {
		m.ResetChallenge()
	}

	// n.Feature = newFeature(n.Feature.Type) // Preserve type in case map contains feature type info
	n.Feature.ResetChallenge()
	n.initAddressMap()
}

func (n *Node) initAddressMap() {
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

func (n *Node) setMacAddress(address string, mac *machines.Machine) {
	// TODO error check name collisions
	mac.Address = address
	n.addressMap[address] = mac

}

// addConnection is reciprocol
func (n *Node) addConnection(m *Node) {
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

func (n *Node) remConnection(ni NodeID) {
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

func (n *Node) HasMachineFor(t teamName) bool { // includes feature
	// t == nil means we don't care... used in calculating Node eccentricity without rewriting dijkstras
	if t == "" {
		return true
	}

	// if a Node has no machines, it allows routing for everyone
	// allows creation of neutral hubs
	if len(n.Machines) == 0 {
		return true
	}

	// if we control a powered machine here we can route through
	for _, mac := range n.Machines {
		if mac.TeamName != "" {
			if mac.TeamName == t { // && mac.Powered {
				return true
			}
		}
	}

	// or if we control a powered feature here we can route through
	if n.Feature.TeamName == t { // && n.Feature.Powered {
		return true
	}

	return false
}

func (n *Node) MachinesFor(t teamName) int { // counts only non-feature machines
	var count int

	for _, mac := range n.Machines {
		if mac.TeamName == t {
			count++
		}
	}

	// if n.Feature.TeamName == t {
	// 	count++
	// }

	return count
}

func (n *Node) supportsRouting(t teamName) bool {
	// t == nil means we don't care... used in calculating Node eccentricity without rewriting dijkstras
	if t == "" {
		return true
	}

	// if a Node has no machines, it allows routing for everyone
	// allows creation of neutral hubs
	if len(n.Machines) == 0 {
		return true
	}

	if n.MachinesFor(t) < 1 && !n.Feature.BelongsTo(t) {
		return false
	}

	return true
}

func (n *Node) PowerMachines(t teamName, onOff bool) {
	for _, mac := range n.Machines {
		if mac.TeamName == t {
			mac.Powered = onOff
		}
	}

	if n.Feature.TeamName == t {
		n.Feature.Powered = onOff
	}
}

// func (n *Node) AddressList() []MacAddress {
// 	addList := make([]string, 0)
// 	for add := range n.addressMap {
// 		addList = append(addList, add)
// 	}
// 	sort.Strings(addList)
// 	return addList
// }

func (n *Node) Addresses() map[MacAddress]*machines.Machine {
	return n.addressMap
}

func (n *Node) MacAt(a MacAddress) *machines.Machine {
	// TODO ERROR check
	return n.addressMap[a]
}

func (n *Node) CanAttach(t teamName, macAddress string) error {
	if _, ok := n.addressMap[macAddress]; !ok {
		return fmt.Errorf("Invalid address, '%s'", macAddress)
	}

	// If for some reason player is forbidden from attaching (ie encryption)...

	return nil
}

// n.resetMachine should never be called directly. only from gm.removeModule
// func (n *Node) resetMachine(slotIndex int) error {
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

func (n Node) String() string {
	return fmt.Sprintf("( <node> {ID: %v, Connections:%v, Machines:%v} )", n.ID, n.Connections, n.Machines)
}
