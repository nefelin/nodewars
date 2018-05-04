package nwmodel

import "strconv"

type route struct {
	Endpoint *node   `json:"endpoint"`
	Nodes    []*node `json:"nodes"`
}

// route methods --------------------------------------------
func (r route) containsNode(n *node) (int, bool) {
	for i, node := range r.Nodes {
		if n == node {
			return i, true
		}
	}
	return 0, false
}

func (r route) length() int {
	return len(r.Nodes)
}

// traficMap and packet used to create more easily rendered statemessages
type trafficMap struct {
	Traffic map[string][]packet `json:"traffic"`
}

type packet struct {
	Owner     string `json:"owner"`
	Direction string `json:"dir"`
}

func newTrafficMap() *trafficMap {
	return &trafficMap{
		Traffic: make(map[string][]packet),
	}
}

// TODO route is stored reversed?
func (t *trafficMap) addRoute(r *route, color string) {
	for i, n := range r.Nodes {
		n1 := n.ID
		var n2 nodeID

		var dir, edge string

		if i == len(r.Nodes)-1 {
			n2 = r.Endpoint.ID
		} else {
			n2 = r.Nodes[i+1].ID
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

func (t *trafficMap) appendPacket(p packet, edge string) {
	_, ok := t.Traffic[edge]
	if !ok {
		t.Traffic[edge] = []packet{p}
	} else {
		t.Traffic[edge] = append(t.Traffic[edge], p)
	}
}
