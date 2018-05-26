package statemessage

import (
	"model/node"
)

type Message struct {
	*node.Map
	Alerts    []Alert     `json:"alerts"`
	PlayerLoc node.NodeID `json:"player_location"`
	*TrafficMap
	// Traffic map[string]trafficPacket `json:"traffic"` // maps edgeIds to lists of packets on that edge
}

type Alert struct {
	Actor    string      `json:"team"`
	Location node.NodeID `json:"location"`
}
