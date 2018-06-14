package node

import (
	"fmt"
	"math/rand"
	"time"
)

// map generators should return errors TODO

func MainMap(n int) (*Map, error) {

	m := NewMap()

	ring := m.addNodes(8)
	ringConnect(m, ring)

	center := m.addNodes(1)
	m.connectNodes(ring[2].ID, center[0].ID)
	m.connectNodes(ring[6].ID, center[0].ID)

	poes := m.addNodes(2)
	m.connectNodes(ring[0].ID, poes[0].ID)
	m.connectNodes(ring[4].ID, poes[1].ID)

	m.initAllNodes()
	m.addPoes(poes[0].ID, poes[1].ID)
	return m, nil
}

func GridMap(n int) (*Map, error) {
	m := NewMap()

	rows := 3
	cols := 4

	total := rows * cols
	m.addNodes(total)
	for i := 0; i < total; i++ {
		if i != 0 && i%cols != 0 {
			m.connectNodes(i, i-1)
		}
		if i < total-4 {
			m.connectNodes(i, i+4)
		}
	}
	m.initAllNodes()

	m.addPoes(0, total-1)
	return m, nil
}

func ClusterMap(n int) (*Map, error) {
	clusterCount := 4
	minNode := 12
	if n < minNode {
		return nil, fmt.Errorf("Ring map requires at least %d nodes, only got %d", minNode, n)
	}
	rand.Seed(time.Now().UTC().UnixNano())

	clusterSize := n / clusterCount
	clusters := make([][]*Node, clusterCount)

	m := NewMap()
	for i := range clusters {
		clusters[i] = m.addNodes(clusterSize)
		ringConnect(m, clusters[i])
		poeLoc := 1

		m.addPoes(clusters[i][poeLoc].ID)

		if i != 0 {
			m.connectNodes(clusters[i][0].ID, clusters[i-1][0].ID)
		}

		if i == len(clusters)-1 {
			m.connectNodes(clusters[i][0].ID, clusters[0][0].ID)
		}
	}
	m.initAllNodes()
	return m, nil
}

// func RingMap(n int) (*Map, error) {
// 	n = 10
// 	minNode := 10

// 	if n < minNode {
// 		return nil, fmt.Errorf("Ring map requires at least %d nodes, only got %d", minNode, n)
// 	}
// 	rand.Seed(time.Now().UTC().UnixNano())

// 	m := NewMap()

// 	outer := m.addNodes(math.Floor(n*2/3))
// 	center :=
// 	inner := n - outer - inner

// 	i := 0
// 	for outer > 0 {
// 		outer--
// 		m.connectNodes(i, i+1)
// 		i++
// 	}

// 	for inner >0 {
// 		inner --
// 		m.connectNodes(i, i+1)
// 		i++
// 	}

// 	return &m, nil
// }

func NewRandMap(n int) (*Map, error) {
	rand.Seed(time.Now().UTC().UnixNano())
	nodeCount := n
	newMap := NewMap()

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
	return newMap, nil
}

func newDefaultMap(n int) (*Map, error) {
	newMap := NewMap()

	nodecount := n

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

	return newMap, nil
}

func ringConnect(m *Map, nodes []*Node) {
	for i, n := range nodes {
		if i < len(nodes)-1 {
			m.connectNodes(n.ID, nodes[i+1].ID)
		} else {
			m.connectNodes(n.ID, nodes[0].ID)
		}
	}
}
