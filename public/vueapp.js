'use strict';
const versionNum = '0.0.0.1';
const versionTag = 'NodeWars:' + versionNum;
const confirmationPhrase = 'Welcome to NodeWars';

// codeCmds are playerCmds that need to submit the codebox to the messenger.
// codeBox content is only submitted with this type of command
const codeCmds = [
	"mk", "make", "makemod",
	"um", "un", "unmake",
	"rf", "ref", "refac"
]

// Minimal Vue Stuff
const vueApp = new Vue({
	el: '#app',

	data: {
			codeBox: 'Your code here',
			playerCmd: '',
			ws: null,
			graphInitialized: false,
			response: '',
			stdin: 'stdin test input'

		},

	methods: {
		sendMsg: function () {
			
			const cmd = this.playerCmd.split(" ")[0]
			let message

			// console.log('index of cmd', cmd, ":" ,codeCmds.indexOf(cmd));
			if (codeCmds.indexOf(cmd) != -1) {
				console.log("sending code along too")
				message = JSON.stringify({
					type: "playerCmd",
					data: this.playerCmd,
					code: this.codeBox
				})
			} else {
				// console.log("sending just command")
				message = JSON.stringify({
					type: "playerCmd",
					data: this.playerCmd
				})
			}

			console.log("seidning message:", message)

			this.playerCmd = ''
			this.ws.send(message)
		},

		autoLogin: function() {

			const stateReq = JSON.stringify({
				type: "stateRequest",
				data: ""
			});
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
			const msg = JSON.parse(e.data);
			switch (msg.sender) {
				case "pseudoServer":
					// const p = document.createElement("p")
					const t = document.createTextNode(msg.type + " " + msg.data)
					const out = document.querySelector("#psOutput")
					out.appendChild(t)
					out.appendChild(document.createElement("br"))
					
					// this.response = msg.type + " " + msg.data + "<br>" + this.response
					// inelegant trimming of whitespace when msg.type is blank TODO
					msg.type == "" ? console.log(msg.data)
								   : console.log(msg.type, msg.data)
					break
				case "server":
					if (msg.type == "gameState"){
						this.updateState(JSON.parse(msg.data))
						break
					}
					console.log("(server) ", msg.type, ":", msg.data)
					break

				default:
					console.log(e)
			}
			// switch (message.type) {
			// 	case "teamAssign":
			// 		this.teamColor = message.data
			// 		break
			// 	case "teamChat":
			// 		console.log("("+this.teamColor+")", message.sender + ":", message.data)
			// 		break
			// 	case "allChat":
			// 		console.log("(all)", message.sender + ":", message.data)
			// 		break
			// 	case "error":
			// 		console.log("> server error < \n", message.data)
			// 		break
			// 	case "gameState":
				
			// 		// console.log(message.data)
			// 		this.updateState(JSON.parse(message.data))
			// 		break
			// 	default:
			// 		console.log('unhandled server response:');
			// 		console.log("<", message)

			// }
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
