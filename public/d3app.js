const width = 640,
	  height = 480,
	  nodeBaseRadius = 16,
	  nodeRadiusMultiplier = nodeBaseRadius/4,
	  strokeWidth = 2

let nodeGroups = null,
	nodeMains = null,
	nodeLabels = null,
	nodeModules = null,
	nodeTraffics = null,
	link = null,
	svg = null,
	layout = null,
	simulation = null

// const colors = d3.scale.category10();

function reveal() {
	// TODO don't do this after first reveal
	// console.log('force layout settled');
	svg.selectAll(".nodes, line").transition()
	   .duration(400)
	   .style("opacity", 1)

}

function svgInit() {
	svg = d3.select('#graph')
		    .style("border", "1px solid black")
		    .attr('width', width)
		    .attr('height', height)
		    .append('svg')
		    .attr('width', width)
		    .attr('height', height) // TODO dynamically pull values, particularly on resize
}

// Arrayify turns objects into lists that d3 likes
// function arrayifyNodeMap (nodeMap) {
// 	newNodes = []
// 	newEdges = []
// 	for (key of Object.keys(nodeMap.nodes)){
// 		// console.log(nodeMap.nodes[key])
// 		newNodes.push(nodeMap.nodes[key])
// 	}
// 	for (key of Object.keys(nodeMap.edges)){
// 		// console.log(nodeMap.nodes[key])
// 		newEdges.push(nodeMap.edges[key])
// 	}

// 	return {nodes: newNodes, edges: newEdges}
// }



function updateGraph (gameState) {
	edgify(gameState.map)
	attachPOEs(gameState)
	attachCoords(gameState.map)
	const nodeMap = gameState.map
	console.log('updateGraph')

	// console.log("state:", gameState)
	console.log("nodeMap:", nodeMap)

	console.log("ng data before", nodeGroups.select("circle").data())	

	nodeGroups.data(nodeMap.nodes)
		.select("circle") // select instead of selectAll auto binds to parent data
		// .attr("class", d => {
		// 	if (d.poe.length > 0) {
		//  		if (d.poe[0].team != null)
		//  			return "POE"
		//  	}
		 	
		// })
		.style("fill", function(d) {
		 	if (d.poes.length > 0) {
		 		if (d.poes[0].team != null)
		 			return d.poes[0].team.name
		 	}
		 	return "white"
		 })
		// .classed("player-connected", d => d.connectedPlayers.length>0)
		// .classed("traffic", d => d.traffic.length>0)




	link.data(nodeMap.edges)

	console.log("ng data after", nodeGroups.data())	

	// // Since we are updating the data with new objects,
	// // we need to point our simulation at the new objects 
	// // to ensure continued tracking:
	simulation.nodes(nodeMap.nodes)
	simulation.force("link")
            .links(nodeMap.edges);

	// // update node and edge traffic
		 
	// // update player POEs and ongoing connections

	// // update module contents
}

function attachCoords(nodeMap) {
	const data = nodeGroups.data()

	// console.log('nodeMap.nodes before:', nodeMap.nodes)
	for (let i=0; i<data.length; i++) {
		nodeMap.nodes[i].x = data[i].x
		nodeMap.nodes[i].y = data[i].y
	}
	// console.log('nodeMap.nodes after:', nodeMap.nodes)
}


// Combine all functions that iterate over all nodes into single function to reduce runtime TODO
function attachPOEs(gameState) {
	for (let node of gameState.map.nodes) {
		node.poes = []
	}

	// console.log("pre attachPoes", gameState.map.nodes)

	for (let playerID of Object.keys(gameState.poes)) {
		const poeID = gameState.poes[playerID].id
		gameState.map.nodes[poeID].poes.push(gameState.players[playerID])
	}
	// console.log("post attachPoes", gameState.map.nodes)
}

function edgify(nodeMap) {

	// console.log('edgifying', nodeMap)
	seenEdges = {}
	nodeMap.edges = []

	for (let i in nodeMap.nodes){

		i = parseInt(i)
		for (let connectionID of nodeMap.nodes[i].connections) {
			let edgeID = ""
			if (i > connectionID)
				edgeID = i + "e" + connectionID
			else
				edgeID = connectionID + "e" + i

			if (!seenEdges[edgeID])
				seenEdges[edgeID] = {id:edgeID, source:i, target:connectionID}
		}
	}

	for (let edgeID of Object.keys(seenEdges)) {
		nodeMap.edges.push(seenEdges[edgeID])
	}
	// console.log("edgify produced:",nodeMap)


}

// initGraph creates the actual dom elements and provides necessary class tags etc
// as well as setting up rules for interactivity (i.e. zoom, drag, etc)
// updateGraph is responsible for actually mapping game data to individual nodes/edges
function initGraph (nodeMap) {
		console.log('Initializing Graph...');

		edgify(nodeMap)

		simulation = d3.forceSimulation()
            .force("link", d3.forceLink().distance(nodeBaseRadius*4))
            .force("collide",d3.forceCollide( function(d){ return d.r + 8 }) )
            .force("charge", d3.forceManyBody().strength(-300))
            .force("center", d3.forceCenter(width / 2, height / 2))

    
        link = svg.selectAll('.edge')
        	.data(nodeMap.edges)
            .enter()
            .append("line")
            .attr("class", "edge")
            .attr("stroke", "black")
            .attr("stroke-width", strokeWidth)

        nodeGroups = svg.selectAll(".node")
        	.data(nodeMap.nodes)
        	.enter()
        	.append("g")
            .attr("class", "node-group")

        nodeMains = nodeGroups.append("circle")
        	.attr("r", d=>nodeBaseRadius+d.connections.length*nodeRadiusMultiplier )
        	.style("stroke", "black")
        	.attr("stroke-width", strokeWidth)
        	.attr("class", "node-main")

        // nodeTraffics = nodeGroups.append("circle")
        // 	.attr("r", d=>3)
        // 	.style("stroke", "black")
        // 	.style("fill", "orange")
        // 	.attr("stroke-width", 1)
        // 	.attr("class", "node-module")
        // 	.attr("cx", function(d) {
        // 		const mainsSiblingRadius = this.parentNode.childNodes[0].r.baseVal.value
        // 		return mainsSiblingRadius + this.r.baseVal.value
        // 	})
	       //  .attr("cy", 0)

        nodeLabels = nodeGroups.append("text")
	       .attr("dx", -nodeBaseRadius*.6)
	       .attr("dy", -nodeBaseRadius*2.8)
	       .attr("class", "node-label")
	       .text(d=>"ID: " + d.id 
	       // 			// "\nConnected Players: " + d.connectedPlayers +
	       // 			"\nPOE: " + d.poe
	       // 			// "\nModules: " + d.modules + 
	       // 			// "\nTraffic: " + d.traffic
	       			)

        
        var ticked = function() {
            link
                .attr("x1", function(d) { return d.source.x; })
                .attr("y1", function(d) { return d.source.y; })
                .attr("x2", function(d) { return d.target.x; })
                .attr("y2", function(d) { return d.target.y; });
    
    		nodeGroups.attr("transform", d => { return "translate(" + d.x + "," + d.y + ")"; })
        }

        simulation
            .nodes(nodeMap.nodes)
            .on("tick", ticked);
    
        simulation.force("link")
            .links(nodeMap.edges);  

	console.log("Calculation Layout...")

}
