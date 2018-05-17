package nwmodel

import (
	"errors"
	"fmt"
	"nwmessage"
	"strconv"

	"validation"
)

var mapCmdList = map[string]playerCommand{
	"nm": cmdNewBlankMap,

	"nrm": cmdNewRandMap,

	"an": cmdAddNodes,

	"as": cmdAddString,
	// "rn": cmdRemoveNodes,
	// unlink
	// remove nodes

	"ln": cmdLinkNodes,

	//remoteness?????

	// ring redundant with line?

	"ap": cmdAddPoes,
	// "rp": cmdRemPoes,

	"bake": cmdBakeMap,

	// slot node_num, [criteria] what (takes challenge of this but not that, etc...)
}

func cmdAddNodes(p *Player, gm *GameModel, args []string) nwmessage.Message {
	var nodeCount int
	if len(args) != 2 {
		return nwmessage.PsError(errors.New("Need (2) args, a quantity to generate and a node to anchor to"))
		// args = append(args, "1")
	}

	nodeCount, err := strconv.Atoi(args[0])
	if err != nil {
		return nwmessage.PsError(err)
	}

	linkTo, err := strconv.Atoi(args[1])
	if err != nil {
		return nwmessage.PsError(err)
	}

	enter := gm.Map.addNodes(nodeCount)

	for _, node := range enter {
		node.addConnection(gm.Map.Nodes[linkTo])
	}

	// gm.broadcastGraphReset()
	gm.psBroadcastExcept(p, nwmessage.PsAlert("Nodes added to map"))
	gm.broadcastState()
	return nwmessage.PsSuccess(fmt.Sprintf("%d new nodes created", nodeCount))
}

func cmdAddString(p *Player, gm *GameModel, args []string) nwmessage.Message {
	var nodeCount int
	if len(args) == 0 {
		// return nwmessage.PsError(errors.New("Need either one or two arguments, received zero"))
		args = append(args, "1")
	}

	nodeCount, err := strconv.Atoi(args[0])
	if err != nil {
		return nwmessage.PsError(err)
	}

	linkTo, err := strconv.Atoi(args[1])
	if err != nil {
		return nwmessage.PsError(err)
	}

	enter := gm.Map.addNodes(nodeCount)

	for i, node := range enter {
		if i == len(enter)-1 {
			node.addConnection(gm.Map.Nodes[linkTo])
		} else {
			node.addConnection(enter[i+1])
		}
	}

	// if linkTo > -1 {
	// 	for _, node := range enter {
	// 		node.addConnection(gm.Map.Nodes[linkTo])
	// 	}
	// }

	// gm.broadcastGraphReset()
	gm.psBroadcastExcept(p, nwmessage.PsAlert("Map was reset"))
	gm.broadcastState()
	return nwmessage.PsSuccess(fmt.Sprintf("%d new nodes created", nodeCount))
}

func cmdAddPoes(p *Player, gm *GameModel, args []string) nwmessage.Message {
	targs := make([]int, len(args))

	for i := range targs {
		// log.Printf("i in loop: %d", i)
		targetNode, err := strconv.Atoi(args[i])
		// log.Printf("targetNode: %d", targetNode)
		if err != nil {
			return nwmessage.PsError(err)
		}
		if nodeCheck := gm.Map.getNode(targetNode); nodeCheck == nil {
			return nwmessage.PsError(fmt.Errorf("Node %d does not exist", targetNode))
		}
		targs[i] = targetNode
	}

	gm.Map.addPoes(targs...)

	// gm.broadcastGraphReset()
	gm.psBroadcastExcept(p, nwmessage.PsAlert("Map was reset"))
	gm.broadcastState()
	return nwmessage.PsSuccess(fmt.Sprintf("Added poes to nodes %v", targs))
}

func cmdBakeMap(p *Player, gm *GameModel, args []string) nwmessage.Message {
	gm.Map.initAllNodes()
	gm.addTeams(makeDummyTeams())
	// gm.broadcastGraphReset()
	gm.psBroadcastExcept(p, nwmessage.PsAlert("Map was baked and is ready to play"))
	gm.broadcastState()
	return nwmessage.PsSuccess(fmt.Sprintf("Map baked and ready to play"))
}

func cmdRemoveNodes(p *Player, gm *GameModel, args []string) nwmessage.Message {
	if len(args) == 0 {
		return nwmessage.PsError(errors.New("Need at least one node to remove"))
	}

	toRemove := make([]int, len(args))
	for i := range toRemove {
		targetNode, err := strconv.Atoi(args[i])
		if err != nil {
			return nwmessage.PsError(err)
		}

		if nodeCheck := gm.Map.getNode(targetNode); nodeCheck == nil {
			return nwmessage.PsError(fmt.Errorf("Node %d not found", targetNode))
		}

		toRemove[i] = targetNode
	}

	gm.Map.removeNodes(toRemove)

	// gm.broadcastGraphReset()
	gm.psBroadcastExcept(p, nwmessage.PsAlert("Map was reset"))
	gm.broadcastState()
	return nwmessage.PsSuccess(fmt.Sprintf("Removed nodes, ", toRemove))
}

func cmdLinkNodes(p *Player, gm *GameModel, args []string) nwmessage.Message {
	if len(args) != 2 {
		return nwmessage.PsError(errors.New("Need exactly two nodes to link"))
	}

	targ := make([]int, 2)
	for i := range targ {
		// log.Printf("i in loop: %d", i)
		targetNode, err := strconv.Atoi(args[i])
		// log.Printf("targetNode: %d", targetNode)
		if err != nil {
			return nwmessage.PsError(err)
		}
		if nodeCheck := gm.Map.getNode(targetNode); nodeCheck == nil {
			return nwmessage.PsError(fmt.Errorf("Node %d does not exist", targetNode))
		}

		targ[i] = targetNode
	}

	// log.Printf("linking node %d to node %d", targ[0], targ[1])
	gm.Map.Nodes[targ[0]].addConnection(gm.Map.Nodes[targ[1]])
	// toLink := make([]*node, len(args))
	// for i := range toLink {
	// 	targetNode, err := strconv.Atoi(args[i])
	// 	if err != nil {
	// 		return nwmessage.PsError(err)
	// 	}

	// 	if !gm.Map.nodeExists(targetNode) {
	// 		return nwmessage.PsError(fmt.Errorf("Node %d not found", targetNode))
	// 	}
	// 	toLink[i] = gm.Map.Nodes[targetNode]
	// }
	// log.Printf("toLink: %v", toLink)

	// for _, node := range toLink {
	// 	for otherNode := range toLink {
	// 		// log.Printf("linking node %d to node %d", node.ID, otherNode)
	// 		log.Printf("linking node %d to node %d", node.ID, gm.Map.Nodes[otherNode].ID)
	// 		node.addConnection(gm.Map.Nodes[otherNode])
	// 	}
	// }
	// gm.broadcastGraphReset()
	gm.psBroadcastExcept(p, nwmessage.PsAlert("Map was reset"))
	gm.broadcastState()
	return nwmessage.PsSuccess("")
}

func cmdNewBlankMap(p *Player, gm *GameModel, args []string) nwmessage.Message {
	if gm.running {
		return nwmessage.PsError(errors.New("Cannot alter map while game is running"))
	}

	if len(args) != 0 {
		return nwmessage.PsError(errors.New("Command does not accept arguments"))
	}

	// nodeIdcount should be irrelevant since its now tied to maps
	// nodeIDCount = 0
	newBlankMap := newNodeMap()
	gm.Map = &newBlankMap
	_ = gm.Map.addNodes(1)
	// gm.broadcastGraphReset()
	gm.psBroadcastExcept(p, nwmessage.PsAlert("Map was reset"))
	gm.broadcastState()
	return nwmessage.PsSuccess("Generating new blank map...")
}

func cmdNewRandMap(p *Player, gm *GameModel, args []string) nwmessage.Message {
	// TODO fix d3 to update...
	if gm.running {
		return nwmessage.PsError(errors.New("Cannot alter map while game is running"))
	}

	argc, err := validation.Validate([]string{"Int"}, args)
	if err != nil {
		return nwmessage.PsError(err)
	}
	nodeCount := argc[0].(int)

	// nodeIdcount should be irrelevant since its now tied to maps
	// nodeIDCount = 0
	gm.Map = newRandMap(nodeCount)
	gm.broadcastGraphReset()
	gm.psBroadcastExcept(p, nwmessage.PsAlert("Map was reset"))
	gm.broadcastState()
	return nwmessage.PsSuccess("Generating new random map...")
}
