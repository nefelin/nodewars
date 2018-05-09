class Incoming {
	constructor(context, debug){
		// this.ws = ws
		this.debug = debug || false
		this.context = context
		console.log('<Parsers.Incoming> context,', context)
	}

	handle = (m) => {
		const data = JSON.parse(m.data)
	    if (this.debug)
	    	if (data.type != "scoreState") console.log('<Parsers.Incoming> received message,', data)

	    switch (data.sender) {
	      case 'pseudoServer':
	        // if (this.debug) console.log('<Parsers.Incoming> routing to terminal')
	        this.context.terminal.current.recv(data.data)
	        break

	      default:
	      	switch (data.type){
	      		case 'teamState':
	      			this.context.setState({team: data.data})
	      			break
	      		case 'graphState':
		      		if (this.context.graph.current!=null)
		      			this.context.graph.current.update(JSON.parse(data.data))
	      			break
	      		case 'graphReset':
	      			if (this.context.graph.current!=null)
		      			this.context.graph.current.reset()
	      			break
	      		case 'resultState':
	      			this.parseResultState(data.data)
	      			break
	      		case 'editorState':
	      			this.context.setState({ aceContent: data.data })
	      			break
	      		case 'stdinState':
	      			this.context.setState({ stdin: data.data })
	      			break
	      	}
	    }
	}

	parseResultState(data) {
		console.log("RESULT", data)
		data = JSON.parse(data)
		

		this.context.setState({
			testResults: data,
			compilerOutput: {type: data.message.type, message:data.message.data},
		})
	}
}

class Outgoing {
	constructor(ws, context, debug){
		this.ws = ws
		this.context = context
		this.debug = debug || false
	}

	parseCmd(cmd) {
		// if (this.debug) console.log('wrapping player command in message object')
		this.ws.send(JSON.stringify({type: 'playerCmd', data: cmd}))
	}

	componentState(component, data) {
		// if (this.debug) console.log('sending', component, 'state:', data)
		this.ws.send(JSON.stringify({type: component + 'State', data: data}))
	}
}

export {Incoming, Outgoing}