'use strict';
const versionNum = '0.0.0.1';
const versionTag = 'NodeWars:' + versionNum;
const confirmationPhrase = 'Welcome to NodeWars';



// Minimal Vue Stuff
new Vue({
	el: '#app',

	data: function() {
		return {
			message: '',
			msgType: 'teamJoin',
			ws: null,
			teamColor: 'white',
			graphInitialized: false
		}
	},

	methods: {
		sendMsg: function () {
			const message = JSON.stringify({
				type: this.msgType,
				data: this.message
			})
			this.message = ''
			// console.log("sending", message, "to socket:", ws)
			this.ws.send(message)
		},

		autoLogin: function() {
			const setName = JSON.stringify({
				type: "setPlayerName",
				data: "dumdum"
			});

			const teamJoin = JSON.stringify({
				type: "teamJoin",
				data: "red"
			});

			const setPOE = JSON.stringify({
				type: "setPOE",
				data: "7"
			});

			const stateReq = JSON.stringify({
				type: "stateRequest",
				data: ""
			});

			this.ws.send(setName)
			this.ws.send(teamJoin)
			this.ws.send(setPOE)
			this.ws.send(stateReq)

		},

		handshake: function() {
			console.log('>', versionTag)
			this.ws.send(versionTag)
		},

		reFocus: function() {
			msgBox.focus()
		},

		updateState: function(newState) {
			console.log('updateState called with state:', newState)
			// console.log("nodeMap before arrayifyNodeMap:", newState.nodeMap)
			// let nodeMap = arrayifyNodeMap(newState.nodeMap)
			let nodeMap = newState.map
			// console.log("nodeMap after arrayifyNodeMap:", newState.nodeMap)
			if (!this.graphInitialized) {
				initGraph(nodeMap)
				this.graphInitialized = true;
			} 
			updateGraph(newState)
			
				
			// update other hud elements
			
		},

		parseServerMessages: function(e) {
			const message = JSON.parse(e.data);
			switch (message.type) {
				case "teamAssign":
					this.teamColor = message.data
					break
				case "teamChat":
					console.log("("+this.teamColor+")", message.sender + ":", message.data)
					break
				case "allChat":
					console.log("(all)", message.sender + ":", message.data)
					break
				case "error":
					console.log("> server error < \n", message.data)
					break
				case "gameState":
				
					// console.log(message.data)
					this.updateState(JSON.parse(message.data))
					break
				default:
					console.log('unhandled server response:');
					console.log("<", message)

			}
		}
	},

	created: function () {
		svgInit();
		
		this.ws = new WebSocket('ws://' + window.location.host + '/ws');

		// Send handshake once socket is open
		this.ws.addEventListener('open', () => {
			this.handshake()
		});

		// Debugging purposes
		this.ws.addEventListener('message', (e) => {
			// console.log('<', e.data)
		});

		// Confirm receipt of correct version number (one time listener)
		this.ws.addEventListener('message', (e)=>{
			if (e.data == confirmationPhrase)
				console.log('> Handshake succesful <')

				// turn on normal message parsing
				this.ws.addEventListener('message', this.parseServerMessages);
				this.ws.addEventListener('close', (d) => {console.log("server severed connection:", d)});
				this.autoLogin()
				return
			console.log('Server said:', e.data)
			throw "Error: failed to negotiate handshake with server"
			this.ws.close()
		}, { once: true });
	}
})
