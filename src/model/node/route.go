package node

import (
	"fmt"
	"strconv"
	"strings"
)

type Route struct {
	Nodes []*Node `json:"nodes"`
}

// route methods --------------------------------------------
func (r Route) Endpoint() *Node {
	if len(r.Nodes) < 1 {
		return nil
	}

	return r.Nodes[0]
}

func (r Route) RunsThrough(n *Node) bool { // runsThrough does not check the endpoint
	for i := 0; i < len(r.Nodes)-1; i++ { // don't check last Node
		Node := r.Nodes[i]
		if n == Node {
			return true
		}
	}
	return false
}

func (r Route) Length() int {
	return len(r.Nodes)
}

func (r Route) String() string {
	nodeCount := len(r.Nodes)
	nodeList := make([]string, nodeCount)

	for i, node := range r.Nodes {
		// this loop is a little funny because we are reversing the order of the node list
		// it's reverse ordered in the data structure but to be human readable we'd like
		// the list to read from source to target
		nodeList[nodeCount-i-1] = strconv.Itoa(node.ID)
	}

	return fmt.Sprintf("( <route> {Endpoint: %v, Through: %v} )", r.Endpoint().ID, strings.Join(nodeList, ", "))
}

// func (r route) forMsg() string {
// 	nodeCount := len(r.Nodes)
// 	nodeList := make([]string, nodeCount)

// 	for i, node := range r.Nodes {
// 		// this loop is a little funny because we are reversing the order of the node list
// 		// it's reverse ordered in the data structure but to be human readable we'd like
// 		// the list to read from source to target
// 		nodeList[nodeCount-i-1] = strconv.Itoa(node.ID)
// 	}
// 	return fmt.Sprintf("(Endpoint: %v, Through: %v)", r.Endpoint().ID, strings.Join(nodeList, ", "))
// }
