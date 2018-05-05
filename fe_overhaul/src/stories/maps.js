// export const BlankMap = JSON.parse({"map":{"nodes":[{"id":0,"connections":[1],"machines":[{"type":"","module":null,"powered":false}],"remoteness":1},{"id":1,"connections":[0,2],"machines":[{"type":"","module":null,"powered":false},{"type":"","module":null,"powered":false}],"remoteness":0.8333333333333334},{"id":2,"connections":[1,3],"machines":[{"type":"","module":null,"powered":false},{"type":"","module":null,"powered":false}],"remoteness":0.6666666666666666},{"id":3,"connections":[2,4,9],"machines":[{"type":"","module":null,"powered":false},{"type":"","module":null,"powered":false},{"type":"","module":null,"powered":false}],"remoteness":0.5},{"id":4,"connections":[3,5],"machines":[{"type":"","module":null,"powered":false},{"type":"","module":null,"powered":false}],"remoteness":0.6666666666666666},{"id":5,"connections":[4,6],"machines":[{"type":"","module":null,"powered":false},{"type":"","module":null,"powered":false}],"remoteness":0.8333333333333334},{"id":6,"connections":[5,7],"machines":[{"type":"","module":null,"powered":false},{"type":"","module":null,"powered":false}],"remoteness":1},{"id":7,"connections":[6,8],"machines":[{"type":"","module":null,"powered":false},{"type":"","module":null,"powered":false}],"remoteness":1},{"id":8,"connections":[7,9],"machines":[{"type":"","module":null,"powered":false},{"type":"","module":null,"powered":false}],"remoteness":0.8333333333333334},{"id":9,"connections":[3,8],"machines":[{"type":"","module":null,"powered":false},{"type":"","module":null,"powered":false}],"remoteness":0.6666666666666666}],"poes":{"0":true,"6":true,"7":true}},"teams":{"blue":{"name":"blue","procPow":0,"vicPoints":0},"red":{"name":"red","procPow":0,"vicPoints":0}},"players":{"1":{"id":1,"team":"","route":null,"ChatMode":false}},"poes":{},"pointGoal":1000})
export const SingleNode = {
							"nodes":[
								{"id":0,"connections":[],"machines":[
									{ "owner":"blue", "type":"", "powered":false },
									{ "owner":"blue", "type":"", "powered":false },
									{ "owner":"", "type":"", "powered":true },
									{ "owner":"", "type":"", "powered":true },
									{ "owner":"red", "type":"", "powered":true },
								],
								"player_here":false,"feature":{"type":"cloak","owner":"red", "powered": true},"remoteness":1},
							],
							"alerts": {}, "traffic": {}
						  }

export const SingleNodeChange = {
							"nodes":[
								{"id":0,"connections":[],"machines":[
									{ "owner":"blue", "type":"", "powered":true },
									{ "owner":"blue", "type":"", "powered":true },
									{ "owner":"blue", "type":"", "powered":true },
									{ "owner":"red", "type":"", "powered":false },
									{ "owner":"red", "type":"", "powered":false },
								],
								"player_here":false,"feature":{"owner":"red","type":"poe","powered":false},"remoteness":1},
							],
							"alerts": {}, "traffic": {}
						  }

export const SingleNodeNoFeature = {
							"nodes":[
									{"id":0,"connections":[],"machines":[
										{ "owner":"blue", "type":"", "powered":false },
										{ "owner":"blue", "type":"", "powered":false },
										{ "owner":"blue", "type":"", "powered":false },
										{ "owner":"red", "type":"", "powered":true },
										{ "owner":"red", "type":"", "powered":false },
									],
									"player_here":false,"feature":{"type":"","owner":"", "powered":true},"remoteness":1},
								],
								"alerts": {}, "traffic": {}
							  }

export const DoubleNode = {
							"nodes":[
								{"id":0,"connections":[1],"machines":[
									{ "owner":"blue", "type":"", "powered":true },
									{ "owner":"", "type":"", "powered":true },
								],
								"player_here":false,"feature":{"type":"","owner":"", "powered": true},"remoteness":1},
								{"id":1,"connections":[0],"machines":[
									{ "owner":"blue", "type":"", "powered":true },
									{ "owner":"", "type":"", "powered":true },
								],
								"player_here":false,"feature":{"type":"","owner":"", "powered": true},"remoteness":1},
							],
							"alerts": {},
							"traffic": {
									"1e0": [
										{"owner":"red", "dir":"up"},
									],
								},
						  }
export const DoubleNodeChange = {
							"nodes":[
								{"id":0,"connections":[1],"machines":[
									{ "owner":"blue", "type":"", "powered":true },
									{ "owner":"", "type":"", "powered":true },
								],
								"player_here":false,"feature":{"type":"","owner":"", "powered": true},"remoteness":1},
								{"id":1,"connections":[0],"machines":[
									{ "owner":"blue", "type":"", "powered":true },
									{ "owner":"", "type":"", "powered":true },
								],
								"player_here":false,"feature":{"type":"","owner":"", "powered": true},"remoteness":1},
							],
							"alerts": {},
							"traffic": {
									"1e0": [
										{"owner":"red", "dir":"up"},
										{"owner":"blue", "dir":"down"},
									],
								},
						  }


export const SimpleMap = {"nodes":[
									{"id":0,"connections":[1,2,4],"machines":[
										{ "owner":"", "type":"", "powered":true },
										{ "owner":"", "type":"", "powered":true },
										{ "owner":"", "type":"", "powered":true },
										{ "owner":"blue", "type":"", "powered":true },
										{ "owner":"", "type":"", "powered":true },
									],
									"feature":{ "owner":"blue", "type":"cloak", "powered":false },"remoteness":.5},

									{"id":1,"connections":[0,3,2,6],"machines":[
										{"type":"","owner":"","powered":true},
										{"type":"","owner":"","powered":true}
									],
									"feature":{"type":"poe","owner":"red", "powered":true},"remoteness":1},
									{"id":2,"connections":[0,1,6],"machines":[
										{ "owner":"red", "type":"", "powered":true },
										{ "owner":"red", "type":"", "powered":true },
										{ "owner":"blue", "type":"", "powered":false },
										{ "owner":"red", "type":"", "powered":true },
										{ "owner":"", "type":"", "powered":true },
									],
									"feature":{ "owner":"red", "type":"firewall", "powered":false },"remoteness":1},
									{"id":3,"connections":[1,4],"machines":[
										{"type":"","owner":"","powered":true},
										{"type":"","owner":"","powered":true},
										{"type":"","owner":"","powered":true},
									],
									"feature":{"type":"overclock","owner":"red", "powered":true },"remoteness":0.8333333333333334},
									{"id":4,"connections":[3,5,0],"machines":[
										{ "owner":"blue", "type":"", "powered":true },
										{ "owner":"blue", "type":"", "powered":true },
										{ "owner":"", "type":"", "powered":true },
										{ "owner":"blue", "type":"", "powered":true },
										{ "owner":"", "type":"", "powered":true },
									],
									"feature":{ "owner":"", "type":"", "powered":true },"remoteness":1},
									{"id":5,"connections":[4],"machines":[
										{ "owner":"", "type":"", "powered":true },
										{ "owner":"red", "type":"", "powered":true },
										{ "owner":"blue", "type":"", "powered":true },
									],"feature":{"type":"poe","owner":"blue", "powered":true },"remoteness":1},
									{"id":6,"connections":[1,2],"machines":[
										{"type":"","module":null,"powered":false},
										{ "owner":"blue", "type":"", "powered":true },
										{ "owner":"blue", "type":"", "powered":true },
										{ "owner":"", "type":"", "powered":true },
										{ "owner":"blue", "type":"", "powered":true },
										{ "owner":"", "type":"", "powered":true },
										{ "owner":"red", "type":"", "powered":true },
										{ "owner":"red", "type":"", "powered":true },
										{ "owner":"blue", "type":"", "powered":false },
										{ "owner":"red", "type":"", "powered":true },
										{ "owner":"", "type":"", "powered":true },
									],
									"feature":{ "type":"", "owner":"", "powered":true },"remoteness":1},
								], "player_location": 0,
								"alerts": {
									1: "blue",
									4: "red",
									0: "blue",
								},
								"traffic": {
									"6e1": [
										{"owner":"red", "dir":"up"},
										{"owner":"red", "dir":"up"},
										{"owner":"red", "dir":"up"},
									],
									"5e4": [
										{"owner":"blue", "dir":"down"},
									],
									"4e0": [
										{"owner":"blue", "dir":"down"},
									],
									"2e0": [
										{"owner":"blue", "dir":"up"},
										{"owner":"red", "dir":"down"}
									],
									"2e1": [
										{"owner":"red", "dir":"up"},
										{"owner":"red", "dir":"up"}
									]
								}
								}

export const SimpleMapChange = {"nodes":[
									{"id":0,"connections":[1,2,4],"machines":[
										{ "owner":"", "type":"", "powered":true },
										{ "owner":"", "type":"", "powered":true },
										{ "owner":"", "type":"", "powered":true },
										{ "owner":"blue", "type":"", "powered":true },
										{ "owner":"", "type":"", "powered":true },
									],
									"feature":{ "owner":"blue", "type":"cloak", "powered":true },"remoteness":.5},
									{"id":1,"connections":[0,3,2,6],"machines":[{"type":"","module":null,"powered":false}],"feature":{"type":"poe","owner":"red"},"remoteness":1},
									{"id":2,"connections":[0,1,6],"machines":[
										{ "owner":"red", "type":"", "powered":true },
										{ "owner":"blue", "type":"", "powered":true },
										{ "owner":"blue", "type":"", "powered":false },
										{ "owner":"red", "type":"", "powered":true },
										{ "owner":"", "type":"", "powered":true },
									],
									"feature":{"owner":"red", "type":"firewall", "powered":true},"remoteness":1},
									{"id":3,"connections":[1,4],"machines":[{"type":"","module":null,"powered":false},{"type":"","module":null,"powered":false}],"feature":{"type":"overclock","owner":"red"},"remoteness":0.8333333333333334},
									{"id":4,"connections":[3,5,0],"machines":[
										{ "owner":"blue", "type":"", "powered":true },
										{ "owner":"blue", "type":"", "powered":true },
										{ "owner":"red", "type":"", "powered":true },
										{ "owner":"blue", "type":"", "powered":true },
										{ "owner":"", "type":"", "powered":true },
									],
									"feature":{ "owner":"", "type":"", "powered":true},"remoteness":1},
									{"id":5,"connections":[4],"machines":[
										{ "owner":"", "type":"", "powered":true },
										{ "owner":"red", "type":"", "powered":true },
										{ "owner":"blue", "type":"", "powered":true },
									],"feature":{ "owner":"blue", "type":"poe", "powered": true },"remoteness":1},
									{"id":6,"connections":[1,2],"machines":[
										{"type":"","module":null,"powered":false},
										{ "owner":"blue", "type":"", "powered":true },
										{ "owner":"blue", "type":"", "powered":true },
										{ "owner":"", "type":"", "powered":true },
										{ "owner":"blue", "type":"", "powered":true },
										{ "owner":"", "type":"", "powered":true },
										{ "owner":"red", "type":"", "powered":true },
										{ "owner":"red", "type":"", "powered":true },
										{ "owner":"blue", "type":"", "powered":false },
										{ "owner":"red", "type":"", "powered":true },
										{ "owner":"", "type":"", "powered":true },
									],
									"feature":{ "type":"", "owner":"", "powered":true },"remoteness":1},
								], "player_location": 0,
								"alerts": {
									1: "blue",
									4: "red",
									0: "blue",
								},
								"traffic": {
									"6e1": [
										{"owner":"red", "dir":"up"},
										{"owner":"red", "dir":"up"},
										{"owner":"red", "dir":"up"},
									],
									"5e4": [
										{"owner":"blue", "dir":"down"},
									],
									"4e0": [
										{"owner":"blue", "dir":"down"},
									],
									"2e0": [
										{"owner":"blue", "dir":"up"},
										{"owner":"red", "dir":"down"}
									],
									"2e1": [
										{"owner":"red", "dir":"up"},
										{"owner":"red", "dir":"up"}
									]
								}
								}
								
// export const TwoNodeTraffic = {
// 								"nodes":[
// 									{"id":0,"connections":[1],"machines":[
// 										{ "owner":"", "type":"", "powered":true },
// 										{ "owner":"", "type":"", "powered":true },
// 										{ "owner":"", "type":"", "powered":true },
// 										{ "owner":"blue", "type":"", "powered":true },
// 										{ "owner":"", "type":"", "powered":true },
// 									],
// 									"player_here":true,"feature":{"type":"poe","owner":"blue"},"remoteness":1},
									
// 									{"id":1,"connections":[0,2],"machines":[
// 										{ "owner":"red", "type":"", "powered":true },
// 										{ "owner":"red", "type":"", "powered":true },
// 										{ "owner":"blue", "type":"", "powered":false },
// 										{ "owner":"red", "type":"", "powered":true },
// 										{ "owner":"", "type":"", "powered":true },
// 									],
// 									"feature":{
// 										"type":"firewall",
// 										"owner":"red"
// 									},
// 									"remoteness":1},

// 									{"id":2,"connections":[1],"machines":[
// 										{ "owner":"red", "type":"", "powered":true },
// 										{ "owner":"red", "type":"", "powered":true },
// 										{ "owner":"blue", "type":"", "powered":false },
// 										{ "owner":"red", "type":"", "powered":true },
// 										{ "owner":"", "type":"", "powered":true },
// 									],
// 									"feature":{
// 										"type":"firewall",
// 										"owner":"red"
// 									},
// 									"remoteness":1},
// 								],
// 								"alerts": {}, "traffic": {
// 									"1e0": [
// 										{"owner":"blue", "dir":"down"},
// 										// {"owner":"red", "dir":"up"},

// 									],
// 									"2e1": [
// 										{"owner":"red", "dir":"down"}
// 									],
// 								}
// 							  }

// export const TwoNodeTrafficTwo = {
// 								"nodes":[
// 									{"id":0,"connections":[1],"machines":[
// 										{ "owner":"", "type":"", "powered":true },
// 										{ "owner":"", "type":"", "powered":true },
// 										{ "owner":"", "type":"", "powered":true },
// 										{ "owner":"blue", "type":"", "powered":true },
// 										{ "owner":"", "type":"", "powered":true },
// 									],
// 									"feature":{"type":"poe","owner":"blue"},"remoteness":1},
									
// 									{"id":1,"connections":[0,2],"machines":[
// 										{ "owner":"red", "type":"", "powered":true },
// 										{ "owner":"red", "type":"", "powered":true },
// 										{ "owner":"blue", "type":"", "powered":false },
// 										{ "owner":"red", "type":"", "powered":true },
// 										{ "owner":"", "type":"", "powered":true },
// 									],
// 									"feature":{
// 										"type":"firewall",
// 										"owner":"red"
// 									},
// 									"remoteness":1},

// 									{"id":2,"connections":[1],"machines":[
// 										{ "owner":"red", "type":"", "powered":true },
// 										{ "owner":"red", "type":"", "powered":true },
// 										{ "owner":"blue", "type":"", "powered":false },
// 										{ "owner":"red", "type":"", "powered":true },
// 										{ "owner":"", "type":"", "powered":true },
// 									],
// 									"feature":{
// 										"type":"firewall",
// 										"owner":"red"
// 									},
// 									"player_here":true, "remoteness":1},
// 								],
// 								"alerts": {}, "traffic": {
// 									"1e0": [
// 										{"owner":"red", "dir":"up"},
// 										{"owner":"red", "dir":"up"},
// 										{"owner":"red", "dir":"up"},
// 										{"owner":"red", "dir":"up"},
// 										{"owner":"red", "dir":"up"},
// 										{"owner":"red", "dir":"up"},
// 										{"owner":"red", "dir":"up"},
// 										{"owner":"red", "dir":"up"},
// 										{"owner":"red", "dir":"up"},
// 									],
// 									"2e1": [
// 										{"owner":"red", "dir":"down"}
// 									],
// 								}
// 							  }

// export const TwoNodeNoTraffic = {
// 								"nodes":[
// 									{"id":0,"connections":[1],"machines":[
// 										{ "owner":"", "type":"", "powered":true },
// 										{ "owner":"", "type":"", "powered":true },
// 										{ "owner":"", "type":"", "powered":true },
// 										{ "owner":"blue", "type":"", "powered":true },
// 										{ "owner":"", "type":"", "powered":true },
// 									],
// 									"feature":{"type":"poe","owner":"blue"},"remoteness":1},
									
// 									{"id":1,"connections":[0,2],"machines":[
// 										{ "owner":"red", "type":"", "powered":true },
// 										{ "owner":"red", "type":"", "powered":true },
// 										{ "owner":"blue", "type":"", "powered":false },
// 										{ "owner":"red", "type":"", "powered":true },
// 										{ "owner":"", "type":"", "powered":true },
// 									],
// 									"feature":{
// 										"type":"firewall",
// 										"owner":"red"
// 									},
// 									"remoteness":1},

// 									{"id":2,"connections":[1],"machines":[
// 										{ "owner":"red", "type":"", "powered":true },
// 										{ "owner":"red", "type":"", "powered":true },
// 										{ "owner":"blue", "type":"", "powered":false },
// 										{ "owner":"red", "type":"", "powered":true },
// 										{ "owner":"", "type":"", "powered":true },
// 									],
// 									"feature":{
// 										"type":"firewall",
// 										"owner":"red"
// 									},
// 									"remoteness":1},
// 								],
// 								"alerts": {}, "traffic": {}
// 							  }