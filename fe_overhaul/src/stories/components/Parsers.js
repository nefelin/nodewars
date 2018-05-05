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
	        this.context.setState({tinyTermIn: data.data})
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
	      		case 'editorState':
	      			this.context.setState({ aceContent: data.data })
	      			break
	      		case 'stdinState':
	      			this.context.setState({ stdin: data.data })
	      			break
	      	}
	    }
	}
}

class Outgoing {
	constructor(ws, context, debug){
		this.ws = ws
		this.context = context
		this.debug = debug || false
	}

	parseCmd(msg) {
		if (this.debug) console.log('wrapping player command in message object')
		this.ws.send(JSON.stringify({type: 'playerCmd', data: msg}))
	}
}

export {Incoming, Outgoing}