package model

import (
	"fmt"
)

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
		playerList = append(playerList, string(player.Name()))
	}
	return fmt.Sprintf("( <team> {Name: %v, Players:%v} )", t.Name, playerList)
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
