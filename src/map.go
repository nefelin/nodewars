package main

type Map struct {
	Nodes []*Node `json:"nodes"`
}

func NewDefaultMap() *Map {
	blueHome := NewNode("blue home")
	redHome := NewNode("red home")

	blueHome.addConnection(redHome)
	redHome.addConnection(blueHome)

	return &Map{[]*Node{blueHome, redHome}}
}
