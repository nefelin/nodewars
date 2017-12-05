'use strict'

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
	// TODO combine all this stuff into 
	makeEdges(gameState.map)
	attachPOEs(gameState)
	attachRoutes(gameState)
	attachCoords(gameState.map)
	arrayifyModules(gameState.map)
	const nodeMap = gameState.map
	console.log('updateGraph')

	// console.log("state:", gameState)
	// console.log("nodeMap:", nodeMap)

	var t = d3.transition()
      .duration(750)

	// rebind node-group data
	nodeGroups.data(nodeMap.nodes)

	// update node-mains 
	nodeGroups.select(".node-main")
		.classed("node-poe", d => d.poes.length>0)
		.classed("player-connected", d => d.connectedPlayers.length>0)
		.classed("traffic", d => d.traffic.length>0)
		.transition(t)
		.style("fill", d => {
			if (d.poes.length>0)
				return d.poes[0].team.name
			return "white"
		})

	// update modules 
	nodeGroups.each(function(d) {
       	const modules = d3.select(this).selectAll(".node-module") // select instead of selectAll auto binds to parent data
			.data(d.modList)


		const parentRadius = d3.select(this.parentNode).select(".node-main").attr("r");
		const nodeRadius = parentRadius

		// const nodeRadius = (nodeBaseRadius+d.connections.length*nodeRadiusMultiplier)
		const modRad = nodeBaseRadius/2.5
		const spacing = modRad*3
		const angleInc = 60*0.017453; //convert to radian

		modules.exit()
			.transition(t)
			.attr("r", 0)
			.remove()

		modules.enter()
			   .append("circle")
			   .attr("class", "node-module")
   			   .style("fill", d => d.team.name)
   			   .style("fill-opacity", d => {
   			   		console.log("making module of fill at:" , d);
   			   		return d.health/d.maxHealth
   			   	})
	           .style("stroke", "black")
	           .style("stroke-width", 2)
	           .attr("opacity", 0)
	           .attr("r", d => modRad)
	           // .attr("cy", -nodeRadius/2)
	           .transition(t)
	           .style("opacity", 1)
	           .attr("cx", (d,i) => nodeRadius/2 * Math.cos(-1.5708+angleInc*i))
	           .attr("cy", (d,i) => nodeRadius/2 * Math.sin(-1.5708+angleInc*i))

	    // update all modules
	    modules.transition(t)
			.style("fill", d => d.team.name)
		    .style("fill-opacity", d => {
		   		console.log("making module of fill at:" , d);
		   		return d.health/d.maxHealth
		    })

	}) 

	link.data(nodeMap.edges)
		.classed("traffic", d => d.traffic.length>0)

	// console.log("ng data after", nodeGroups.data())	

	// // Since we are updating the data with new objects,
	// // we need to point our simulation at the new objects 
	// // to ensure continued tracking:
	simulation.nodes(nodeMap.nodes)
	simulation.force("link")
            .links(nodeMap.edges)
}

function arrayifyModules(nodeMap) {
	for (let node of nodeMap.nodes) {
		node.modList = []
		for (let modID of Object.keys(node.modules)){
			node.modList.push(node.modules[modID])
		}
	}
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

function getEdgeIn(edgeID, edgeList) {
	for (let edge of edgeList) {
		// console.log("getEdge is comparing",edge.id, "and", edgeID)
		if (edge.id == edgeID)
			return edge
	}
	return null 
}

function attachRoutes(gameState) {
	for (let node of gameState.map.nodes) {
		node.connectedPlayers = []
		node.traffic = []
	}

	for (let edge of gameState.map.edges) {
		edge.traffic = []
	}

	// console.log("pre attachRoutes", gameState.map)

	for (let playerID of Object.keys(gameState.players)) {
		const player = gameState.players[playerID]
		const route = player.route

		if (route) {
			// iterate in reverse since routes are reverse ordered
			for (let i = route.nodes.length-1; i > -1; i--) {
				// // attach traffic to nodes
				const thisNode = route.nodes[i]
				gameState.map.nodes[thisNode.id].traffic.push(player)

				//attach traffic to edges if we're not connecting to poe
				// if we're in the middle of the route push traffic to connector
				if (route.nodes[0].id != route.endpoint.id){
					let thisEdgeID
					if (i > 0) {
						thisEdgeID = makeEdgeID(route.nodes[i].id, route.nodes[i-1].id)
					} else {
						// otherwise add traffic between last node and target
						thisEdgeID = makeEdgeID(route.nodes[0].id, route.endpoint.id)
					}

					getEdgeIn(thisEdgeID, gameState.map.edges).traffic.push(player)
				}
			}
			// attach endpoints
			gameState.map.nodes[route.endpoint.id].connectedPlayers.push(player)
		}
	}
	// console.log("post attachRoutes", gameState.map)
}

function makeEdgeID(id1, id2) {
	let edgeID
	if (id1 > id2)
		edgeID = id1 + "e" + id2
	else
		edgeID = id2 + "e" + id1
	return edgeID
}

function makeEdges(nodeMap) {

	// console.log('edgifying', nodeMap)
	const seenEdges = {}
	nodeMap.edges = []

	for (let i in nodeMap.nodes){

		i = parseInt(i)
		for (let connectionID of nodeMap.nodes[i].connections) {
			let edgeID = makeEdgeID(i, connectionID)

			if (!seenEdges[edgeID])
				seenEdges[edgeID] = {id:edgeID, source:i, target:connectionID}
		}
	}

	for (let edgeID of Object.keys(seenEdges)) {
		nodeMap.edges.push(seenEdges[edgeID])
	}
	// console.log("makeEdges produced:",nodeMap)


}

// initGraph creates the actual dom elements and provides necessary class tags etc
// as well as setting up rules for interactivity (i.e. zoom, drag, etc)
// updateGraph is responsible for actually mapping game data to individual nodes/edges
function initGraph (nodeMap) {
		console.log('Initializing Graph...');

		makeEdges(nodeMap)

		simulation = d3.forceSimulation()
            .force("link", d3.forceLink().distance(nodeBaseRadius*4))
            .force("collide",d3.forceCollide( function(d){ return d.r + 8 }) )
            .force("charge", d3.forceManyBody().strength(-300))
            .force("center", d3.forceCenter(width / 2, height / 2))

    
        link = svg.selectAll('.edge')
        	.data(nodeMap.edges)
            .enter()
            .append("line")
            .attr("class", d => "edge edgeID-" + d.id)
            .attr("stroke", "black")
            .attr("stroke-width", strokeWidth)

        nodeGroups = svg.selectAll(".node-group")
        	.data(nodeMap.nodes)
        	.enter()
        	.append("g")
            .attr("class", "node-group")

        nodeMains = nodeGroups.each(function(d) {
        	// add backing to hide main's transparency when animating POEs
        	d3.select(this).append("circle")
        	  .attr("r", d=>nodeBaseRadius+d.connections.length*nodeRadiusMultiplier)
        	  .attr("class", "node-backing")

        	// add node-main
        	d3.select(this).append("circle")
        	  .attr("r", d=>nodeBaseRadius+d.connections.length*nodeRadiusMultiplier)
        	  .attr("class", "node-main")

        })

        nodeLabels = nodeGroups.append("text")
	       .attr("dx", -nodeBaseRadius*.8)
	       .attr("dy", -nodeBaseRadius-20) // make spacing programatic TODO
	       .attr("font-size",15)
	       .attr("class", "node-label")
	       .text(d=>"ID: " + d.id)

        
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
