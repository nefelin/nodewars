package nwmodel

// Stringers ----------------------------------------------------------------------------------
import (
	"fmt"
	"strconv"
	"strings"
)

func (n node) String() string {
	return fmt.Sprintf("( <node> {ID: %v, Connections:%v, Machines:%v} )", n.ID, n.Connections, n.Machines)
}

// func (n node) modIDs() []modID {
// 	ids := make([]modID, 0)
// 	for _, slot := range n.Machines {
// 		if slot.TeamName != "" {
// 			ids = append(ids, slot.Module.id)
// 		}
// 	}
// 	return ids
// }

func (t team) String() string {
	var playerList []string
	for player := range t.players {
		playerList = append(playerList, string(player.GetName()))
	}
	return fmt.Sprintf("( <team> {Name: %v, Players:%v} )", t.Name, playerList)
}

func (p Player) String() string {
	return fmt.Sprintf("( <player> {Name: %v, team: %v} )", p.GetName(), p.TeamName)
}

func (r route) String() string {
	nodeCount := len(r.Nodes)
	nodeList := make([]string, nodeCount)

	for i, node := range r.Nodes {
		// this loop is a little funny because we are reversing the order of the node list
		// it's reverse ordered in the data structure but to be human readable we'd like
		// the list to read from source to target
		nodeList[nodeCount-i-1] = strconv.Itoa(node.ID)
	}

	return fmt.Sprintf("( <route> {Endpoint: %v, Through: %v} )", r.Endpoint.ID, strings.Join(nodeList, ", "))
}

func (n node) forMsg() string {

	macList := ""
	for i, mac := range n.Machines {
		macList += "\n" + strconv.Itoa(i) + ":" + mac.forMsg()
	}

	connectList := strings.Trim(strings.Join(strings.Split(fmt.Sprint(n.Connections), " "), ","), "[]")

	return fmt.Sprintf("NodeID: %v\nConnects To: %s\nFeature: \n%s\nMachines: %v", n.ID, connectList, n.Feature, macList)
}

func (m machine) forMsg() string {
	switch {
	case m.TeamName != "":
		return "(" + m.details() + ")"
	default:
		return "( -neutral- )"
	}
}

func (m machine) details() string {
	return fmt.Sprintf("[%s] [%s] [%s] [%d/%d]", m.TeamName, m.builder, m.language, m.Health, m.MaxHealth)
}

func (f feature) String() string {
	retStr := f.machine.forMsg()
	if f.Type != "" {
		retStr += fmt.Sprintf(", %s", f.Type)
	}
	return retStr
}

// func (m machine) forProbe() string {
// 	var header string
// 	switch {
// 	case m.Module != nil:
// 		header = "( " + m.Module.forMsg() + " )\n"
// 	default:
// 		header = "( -empty- )\n"
// 	}
// 	// task := "Task:\n" + m.challenge.Description
// 	return header //+ task

// }

func (r route) forMsg() string {
	nodeCount := len(r.Nodes)
	nodeList := make([]string, nodeCount)

	for i, node := range r.Nodes {
		// this loop is a little funny because we are reversing the order of the node list
		// it's reverse ordered in the data structure but to be human readable we'd like
		// the list to read from source to target
		nodeList[nodeCount-i-1] = strconv.Itoa(node.ID)
	}
	return fmt.Sprintf("(Endpoint: %v, Through: %v)", r.Endpoint.ID, strings.Join(nodeList, ", "))
}

// func (c CompileResult) String() string {
// 	ret := ""
// 	for k, v := range c.Graded {
// 		ret += fmt.Sprintf("(in: %v, out: %v)", k, v)
// 	}
// 	return ret
// }

// func (m machine) String() string {

// }
