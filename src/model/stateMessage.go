package model

type stateMessage struct {
	*nodeMap
	Alerts    []alert `json:"alerts"`
	PlayerLoc nodeID  `json:"player_location"`
	*trafficMap
	// Traffic map[string]trafficPacket `json:"traffic"` // maps edgeIds to lists of packets on that edge
}

type alert struct {
	Actor    string `json:"team"`
	Location nodeID `json:"location"`
}
