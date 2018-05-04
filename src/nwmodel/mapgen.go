package nwmodel

import (
	"math/rand"
	"time"
)

func newRandMap(n int) *nodeMap {
	rand.Seed(time.Now().UTC().UnixNano())
	nodeCount := n
	newMap := newNodeMap()

	// for i := 0; i < nodeCount; i++ {
	// 	newMap.addNodes(newMap.NewNode())
	// }
	newMap.addNodes(nodeCount)

	for i := 0; i < nodeCount; i++ {
		if i < nodeCount-1 {
			newMap.connectNodes(i, i+1)
		}

		for j := 0; j < rand.Intn(2); j++ {
			newMap.connectNodes(i, rand.Intn(nodeCount))
		}

	}

	newMap.initAllNodes()

	newMap.initPoes(2)
	// for len(newMap.POEs) < 2 {
	// 	newMap.addPoes(rand.Intn(nodeCount))
	// }
	return &newMap
}

func newDefaultMap() *nodeMap {
	newMap := newNodeMap()

	nodecount := 12

	// for i := 0; i < nodecount; i++ {
	// 	//Make new nodes
	// 	newMap.addNodes(newMap.NewNode())
	// }

	newMap.addNodes(nodecount)

	for i := 0; i < nodecount; i++ {
		//Make new edges
		targ1, targ2 := -1, -1

		if i < nodecount-3 {
			targ1 = i + 3
		}

		if i%3 < 2 {
			targ2 = i + 1
		}

		if targ1 != -1 {
			newMap.connectNodes(i, targ1)
		}

		if targ2 != -1 {
			newMap.connectNodes(i, targ2)
		}
	}

	// create module slots based on connectivity of node
	for _, node := range newMap.Nodes {
		node.initMachines()
	}

	newMap.addPoes(1, 10)

	return &newMap
}
