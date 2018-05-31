class Incoming {
	constructor(props){
		this.context = props.context
		this.debug = props.debug || false
		// console.log('<Parsers.Incoming> context,', context)
	}

	open = () => {
		this.context.terminal.current.recv('Attempting to contact server...\n')
	}

	error = () => {
		this.context.terminal.current.recv('error: Unable to establish connection\n')	
	}

	close = () => {
		this.context.terminal.current.recv('\n\nerror: Server severed connection\n')
	}

	handle = (m) => {
		const data = JSON.parse(m.data)
	    // if (this.debug)
	    // 	if (data.type != "scoreState") console.log('<Parsers.Incoming> received message,', data)

	    switch (data.sender) {
	      case 'pseudoServer':
	        // if (this.debug) console.log('<Parsers.Incoming>', data)
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
	      		case 'graphResetFocus':
		      		// console.log('REEEEESetting graph focus')
	      			if (this.context.graph.current!=null)
		      			this.context.graph.current.resetFocus()
	      			break
	      		case 'graphFocus':
		      		// console.log('Setting graph focus')
	      			if (this.context.graph.current!=null)
		      			this.context.graph.current.setFocus(parseInt(data.data))
	      			break
	      		case 'challengeState':
	      			const challenge = JSON.parse(data.data)
	      			this.context.setState({ challenge })
	      			break
	      		case 'resultState':
	      			this.parseResultState(data.data)
	      			break
	      		case 'stdinState':
	      			this.context.setState({ stdin: data.data })
	      			break
	      		case 'editorState':
					this.context.setState({ aceContent: data.data })
					break
	      		case 'editorLangState':
	      			// console.log('editorLangState')
	      			const lang = data.data.charAt(0).toUpperCase() + data.data.substr(1)
	      			this.context.setState({ aceMode: lang }, this.context.buildAceMenu)
					break
				case 'langSupportState':
					// console.log('calling buildAceMenu')
					const supportedLanguages = JSON.parse(data.data)
					this.context.setState({ supportedLanguages }, this.context.buildAceMenu)
					break
	      	}
	    }
	}

	parseResultState(data) {
		data = JSON.parse(data)
		this.context.setState({
			testResults: data,
			compilerOutput: {type: data.message.type, message:data.message.data},
		})
	}
}

class Outgoing {
	constructor(props){
		this.ws = props.socket
		this.context = props.context
		this.debug = props.debug || false
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