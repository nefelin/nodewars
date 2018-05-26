package statemessage

import (
	"model/node"
	"strconv"
)

// traficMap and packet used to create more easily rendered statemessages
type TrafficMap struct {
	Traffic map[string][]packet `json:"traffic"`
}

type packet struct {
	Owner     string `json:"owner"`
	Direction string `json:"dir"`
}

func NewTrafficMap() *TrafficMap {
	return &TrafficMap{
		Traffic: make(map[string][]packet),
	}
}

// TODO route is stored reversed?
func (t *TrafficMap) AddRoute(r *node.Route, color string) {
	// fmt.Printf("asIds test: %v", r.asIds())
	ids := route2IDs(r)
	for i := 0; i < len(ids)-1; i++ {
		n1 := ids[i]
		n2 := ids[i+1]

		var dir, edge string

		// if we're connecting to poe, ignore
		if n1 == n2 {
			continue
		}

		if n1 > n2 {
			dir = "down"
			edge = strconv.Itoa(n1) + "e" + strconv.Itoa(n2)
		} else {
			dir = "up"
			edge = strconv.Itoa(n2) + "e" + strconv.Itoa(n1)
		}
		t.appendPacket(packet{color, dir}, edge)
	}
}

func (t *TrafficMap) appendPacket(p packet, edge string) {
	_, ok := t.Traffic[edge]
	if !ok {
		t.Traffic[edge] = []packet{p}
	} else {
		t.Traffic[edge] = append(t.Traffic[edge], p)
	}
}

// asIds reverses the order of the nodes and stores ids only
func route2IDs(r *node.Route) []node.NodeID {
	nodeCount := len(r.Nodes)
	list := make([]node.NodeID, nodeCount)

	for i := 0; i < nodeCount; i++ {
		list[i] = r.Nodes[nodeCount-1-i].ID
	}

	// fmt.Printf("Route: %v\n,Nodecount: %d\nList: %v\n", r, nodeCount, list)
	return list
}
