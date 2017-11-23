package main

func newFuncDesc() funcDesc {
	return funcDesc{
		name:    "",
		inputs:  make([]argDesc, 0),
		outputs: make([]argType, 0),
	}
}

func (fd *funcDesc) setName(n string) {
	fd.name = n
}

func (fd *funcDesc) addInput(ad argDesc) {
	fd.inputs = append(fd.inputs, ad)
}

func (fd *funcDesc) addOutput(at argType) {
	fd.outputs = append(fd.outputs, at)
}

func newTestDesc() testDesc {
	return testDesc{
		id:        "",
		desc:      "",
		protoFunc: newFuncDesc(),
		inputs:    make([]ioExpect, 0),
		outputs:   make([]ioExpect, 0),
	}
}

func (td *testDesc) setName(n string) {
	td.id = n
}

func (td *testDesc) setDesc(n string) {
	td.desc = n
}

func (td *testDesc) addTestPair(i, o ioExpect) {
	td.inputs = append(td.inputs, i)
	td.outputs = append(td.outputs, o)
}

func main() {
	// our proto = "def helloWorld():"
	// userCode := "return \"hello world\""

	test := newTestDesc()

	test.protoFunc.setName("helloWorld")
	test.protoFunc.addOutput(argType{"bull", "string"})
	test.protoFunc.addOutput("string")

	test.addTestPair(ioExpect{"unimportant", "string"}, ioExpect{"Hello World", "string"})
}

/*
in python file should look like


# user submission
def helloWorld(bull):
	return "hello world"

# test processing
results = []

inputs = [values already type converted]
outputs = [values already type converted]

def convertType(value, type):
	#case statement returning the value cast as the type we want?

for i of range(len(inputs)):
	retVal = %protoFunc.name%(inputs[i])
	if retVal == outputs[i]:
		results.append(True)
	else:
		results.append(False)

JSON.output(results)

*/

// name - retOdd
// inputs[0] = name a, type int
// inputs[0] = name b, type int

// // takes two ints, should return the odd value, or a+b if both are odd, or 0 if neither is.
// func retOdd(a, b int) int {

// }
