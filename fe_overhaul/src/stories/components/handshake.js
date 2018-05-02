
const versionNum = '1.0.0';
const versionTag = 'NodeWars:' + versionNum;
const confirmationPhrase = 'Welcome to NodeWars';


class NWSocket {
	constructor(parser) {
		let ws_protocol = 'ws://'
		if (location.protocol == 'https:')
			ws_protocol = 'wss://'

		// const ws = new WebSocket(ws_protocol + window.location.host + '/ws');
		const ws = new WebSocket(ws_protocol + 'localhost:8080' + '/ws');


		function handshake() {
			console.log('>', versionTag)
			ws.send(versionTag)
		}

		ws.addEventListener('open', () => {
					handshake()
				});

		ws.addEventListener('message', (e)=>{
			if (e.data == confirmationPhrase) {
				console.log('> Handshake succesful <')

				// turn on normal message parsing
				ws.addEventListener('message', parser);

				// handle server terminating connection
				ws.addEventListener('close', (d) => {console.log("server severed connection:", d)});
				return
			}
			console.log('Server said:', e.data)
			throw "Error: failed to negotiate handshake with server"
			ws.close()
		}, { once: true });
	}

	send(msg) {
		console.log('donothing that should send msg,',msg)
	}
}


function parser(m) {
	console.log('received message,', m)
}

export default NWSocket