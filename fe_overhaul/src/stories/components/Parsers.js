class Incoming {
	constructor(context, debug){
		// this.ws = ws
		this.debug = debug || false
		this.context = context
		console.log('<Parsers.Incoming> context,', context)
	}

	handle = (m) => {
		const data = JSON.parse(m.data)
	    if (this.debug) console.log('<Parsers.Incoming> received message,', data)

	    switch (data.sender) {
	      case 'pseudoServer':
	        if (this.debug) console.log('<Parsers.Incoming> routing to terminal')
	        this.context.setState({tinyTermIn: data.data})
	        break
	      default:
	      	switch (data.type){
	      		case 'teamState':
	      			this.context.setState({team: data.data})
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