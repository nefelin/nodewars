'use strict'
// TODO review the class structure and make those methods that don't rely on instance data STATIC

// class constants (only used internally)
// TODO grab width and height from parent container. Esp on resize
const width = 640,
	  height = 480,
	  NODE_BASE_RADIUS = 16,
	  NODE_RAD_MULTIPLIER = NODE_BASE_RADIUS/4,
	  BASE_STROKE_WIDTH = 2

const t = d3.transition()
      .duration(1000)


class NWGraph {
	constructor(targetID) {

		this.graphDiv = document.querySelector(targetID)
		this.width = this.graphDiv.clientWidth
        this.height = this.graphDiv.clientHeight
        console.log("width", this.width, "height", this.height)

		console.log("NWGraph creating svg...")
		this.svg = d3.select(targetID)
		    // .style("border-right", "1px solid black")
		    .append('svg')
		    .attr('width', "100%")
		    .attr('height', "100%")

		this.simulation = d3.forceSimulation()
            .force("link", d3.forceLink().distance(NODE_BASE_RADIUS*4))
            .force("collide",d3.forceCollide( function(d){ return d.r + 8 }) )
            .force("charge", d3.forceManyBody().strength(-300))
            .force("center", d3.forceCenter(this.width / 2, this.height / 2))

		this.simulation.stop()
        this.stopped = true

	}

	reset() {
		this.stopped = true

		this.simulation.stop()
		if (this.links)
			this.links.remove()
		if (this.nodeGroups)
			this.nodeGroups.remove()

		this.resize()

		this.gameState = null
	}	

	resize() {
		this.width = this.graphDiv.clientWidth
        this.height = this.graphDiv.clientHeight
  		// this.width = 640
  		// this.height = 480
        // console.log("width", this.width, "height", this.height)
        this.simulation.force("center", d3.forceCenter(this.width / 2, this.height / 2))
	}

	draw() {
	// const self = this

	// Order of these is important as D3 handles z-index by draw order only
	this.drawLinks()
	this.drawNodeGroups()
    this.drawNodeLabels()

    // this.drawModules()
    this.drawSlots()

	}

	update(newState) {
		console.log("NWGraph updating...")

		const self = this

		NWGraph.makeEdges(newState.map)
		NWGraph.attachPOEs(newState)
		NWGraph.attachRoutes(newState)
		NWGraph.arrayifyModules(newState.map)
		
		// if we're updating pre-existing state
		if (this.gameState) {
			NWGraph.attachCoords(this.nodeGroups.data(), newState.map)
		}

		this.gameState = newState

		if (this.stopped){
			this.simulation.alpha(1)
			this.simulation.restart()
			console.log("simulation restarting")
			this.stopped = false
		}

		var ticked = function() {
            self.links
                .attr("x1", function(d) { return d.source.x; })
                .attr("y1", function(d) { return d.source.y; })
                .attr("x2", function(d) { return d.target.x; })
                .attr("y2", function(d) { return d.target.y; });
    
    		self.nodeGroups.attr("transform", d => { return "translate(" + d.x + "," + d.y + ")"; })
        }

        this.simulation
            .nodes(this.gameState.map.nodes)
            .on("tick", ticked);
    
        this.simulation.force("link")
            .links(this.gameState.map.edges); 
		// console.log("NWGraph state post update:", this.gameState)
	}

	drawLinks() {
		const update = this.svg.selectAll('.edge').data(this.gameState.map.edges)

		this.links = update.enter()
            .append("line")
            .attr("class", d => "edge edgeID-" + d.id)
            .attr("stroke", "black")
            .attr("stroke-width", BASE_STROKE_WIDTH)
          .merge(update)
            .classed("traffic", d => d.traffic.length>0)

        update.exit().remove()

        // update traffic info on all links, new and old alike
        // this.links.classed("traffic", d => d.traffic.length>0)
	}

	drawNodeGroups() {
		const self = this
		const update = this.svg.selectAll(".node-group").data(this.gameState.map.nodes)
		const enter = update.enter()
        	.append("g")
            .attr("class", "node-group")
            // .merge(update)
            .call(d3.drag()
                .on("start", this.dragstarted.bind(this))
                .on("drag", this.dragged.bind(this))
                .on("end", this.dragended.bind(this))); 
            
		enter.each(function(d) {
        	const radius = NODE_BASE_RADIUS+d.connections.length*NODE_RAD_MULTIPLIER

        	// add backing to hide main's transparency when animating POEs
        	d3.select(this).append("circle")
        	  .attr("r", d=>radius)
        	  .attr("class", "node-backing")
        	  .attr("id", d => "node-backing-" + d.id)
        	  

        	// add node-main
        	d3.select(this).append("circle")
        	  .attr("r", d=>radius)
        	  .attr("class", "node-main")
        	  .attr("id", d => "node-main-" + d.id)


        	// add potential poe indicators 
        	const potentialPoe = Object.keys(self.gameState.map.poes).indexOf(String(d.id))

        	if (potentialPoe > -1) {
        		d3.select(this).append("circle")
        		  .attr("r", d=>radius*.85)
        		  .attr("class", "potentialPOE")
        		  .style("stroke-width", 2)
          		  .style("stroke", "black")
        		  .style("fill-opacity", 0)
        	}
        })

		this.nodeGroups = enter.merge(update)

        update.exit().remove()

		// update node-main classes for new and old nodes alike

        this.nodeGroups.select(".node-main")
			.classed("node-poe", d => d.poes.length>0)
			.classed("player-connected", d => d.connectedPlayers.length>0)
			.classed("traffic", d => d.traffic.length>0)
			.transition(t)
			.style("fill", d => {
				// console.log(d)
				if (d.poes.length>0)
					return d.poes[0].team
				return d3.hsl(165, 1*d.remoteness*.6, 1*(1-d.remoteness*.9)).brighter(3)
				// return d3.interpolateBuGn(d.remoteness*.8+Math.random()/15)
				// return d3.interpolateSpectral(d.remoteness*.6+Math.random()/15)
				// return "white"
				// return d3.hsl(160+(30*d.remoteness), (d.remoteness+(d.connections.length*.2)), (1-d.remoteness)+(d.connections.length*.1)).brighter(3)
				// return d3.hsl(d.remoteness*300, d.remoteness, 1-d.remoteness*.4)
				// return d3.hsl(250, (d.remoteness+(d.connections.length*.2)), sigmoid((1-d.remoteness)+(d.connections.length*.1)))
			})
	}

	// sigmoid(t) {
	//     return 1/(1+Math.pow(Math.E, -t));
	// }	

	drawNodeLabels() {
		// const nodeRadius = 
		this.nodeLabels = this.nodeGroups.append("text")
	       // .attr("dx", function(d) {
	       		// console.log("nodegroup?", this.parentNode)
	       		// d3.select(this.parentNode).select(".node-main").attr("r")*.8
	        // })
	       .attr("dy", function(d) {
		       	// console.log("this", d3.select(this).getBoundingClientRect())
		       	// console.log("this getBBox", d3.select(this).node().getBBox())
		       	// console.log("this getBoundingCliRect", d3.select(this).node().getBoundingClientRect())
	       		// const textBounds = d3.select(this).node().getBBox()
	       		const parentRad = d3.select(this.parentNode).select(".node-main").attr("r")
	       		const thisPos =  -parentRad-12
	       		return thisPos
	       	}) 
	       .attr("text-anchor", "middle")
	       .attr("alignment-baseline", "middle")
	       .attr("font-size", "20")
	       .attr("class", "node-label")
	       .text(d=>d.id)
	       // .style("opacity", .3)
	}

	drawSlots() {
		this.nodeGroups.each(function(d) {
			const slots = d3.select(this).selectAll(".mod-slot")
				.data(d.slots)
			
			const parentD = d
			// console.log("slots d: ", d)
	
			const thisGroup = d3.select(this)
			const nodeRadius = thisGroup.select(".node-main").attr("r");

			// console.log("parent", d3.select(this))
			// console.log(nodeRadius)

			const modRad = NODE_BASE_RADIUS/2.5
			const slotRad = modRad/4
			const slotWidth = 10
			const slotHeight = 10
			const spacing = 5
			const colCount = 2
			// const spacing = modRad*3
			// const angleInc = 70*0.017453; //convert to radian

			const angleInc = (360/d.slots.length)*0.017453
			
			slots.exit()
					.transition(t)
					.attr("r", 0)
					.remove()

			slots.enter()
				   .append("rect")
				   .attr("class", "mod-slot")
	   			   .style("fill", d => {
						// console.log("slot d", d)
	   			   		return d.module == null ? "white" : d.module.team
	   				})
	   			   .style("fill-opacity", d => d.module == null ? .3 : d.module.health/d.module.maxHealth)
	   			   .style("stroke-opacity", d => d.module == null ? .8 : 1)
		           .style("stroke", "black")
		           .style("stroke-width", 2)
		           .attr("opacity", 0)
		           .attr("width", d => d.module == null ? slotWidth : slotWidth)
		           .attr("height", d => d.module == null ? slotHeight : slotHeight)
		           .transition(t)
		           .style("opacity", 1)
		           // .attr("x", (d,i) => ((nodeRadius/3) * Math.cos(-1.5708+angleInc*i))-5)
		           // .attr("y", (d,i) => ((nodeRadius/3) * Math.sin(-1.5708+angleInc*i))-5)
		           .attr("x", (d,i) => {
		           		const offset = -(colCount * slotWidth + (colCount-1)*spacing)/2
		           		return offset+(slotWidth+spacing)*(i%colCount)
		           	})
		           .attr("y", (d,i) => {
		           		const rows = Math.ceil(parentD.slots.length/colCount)
		           		const offset = -(rows * slotHeight + (rows-1)*spacing)/2/2 // Should need last /2, something is off TODO
		           		return offset+(-slotHeight/2)+(Math.floor(i/2)*(slotHeight+spacing))
		           })
		           

		    // update new slots
		    slots.transition(t)
				.style("fill", d => d.module == null ? "white" : d.module.team)
			    .style("fill-opacity", d => d.module == null ? .3 : d.module.health/d.module.maxHealth)
			    .style("stroke-opacity", d => d.module == null ? .8 : 1)
		})
	}

	// Simulation drag helpers -------------------------------------------------------------------

	dragstarted(d) {
        if (!d3.event.active) this.simulation.alphaTarget(0.3).restart();
        d.fx = d.x;
        d.fy = d.y;
    }
    
    dragged(d) {
        d.fx = d3.event.x;
        d.fy = d3.event.y;
    }
    
    dragended(d) {
        if (!d3.event.active) this.simulation.alphaTarget(0);
        d.fx = null;
        d.fy = null;
    } 

	// Helper methods and methods for pre-treating gameState before drawing ----------------------

	alertFlash(color, targetNode) {
	
		const startColor = this.svg.style("background-color")

		this.svg.transition()
			 .style("background-color", d3.hsl(color).brighter(1))
			 .transition()
			 .ease(d3.easeLinear)
			 .duration(800)
			 .style("background-color", startColor)
	}

	static arrayifyModules(nodeMap) {
		for (let node of nodeMap.nodes) {
			node.modList = []
			for (let modID of Object.keys(node.modules)){
				node.modList.push(node.modules[modID])
			}
		}
	}

	static attachCoords(oldData, nodeMap) {
		// const data = this.nodeGroups.data()

		// console.log('nodeMap.nodes before:', nodeMap.nodes)
		for (let i=0; i<oldData.length; i++) {
			nodeMap.nodes[i].x = oldData[i].x
			nodeMap.nodes[i].y = oldData[i].y
		}
		// console.log('nodeMap.nodes after:', nodeMap.nodes)
	}


	// Combine all functions that iterate over all nodes into single function to reduce runtime TODO
	static attachPOEs(gameState) {
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

	static getEdgeIn(edgeID, edgeList) {
		for (let edge of edgeList) {
			// console.log("getEdge is comparing",edge.id, "and", edgeID)
			if (edge.id == edgeID)
				return edge
		}
		return null 
	}

	static makeEdgeID(id1, id2) {
		let edgeID
		if (id1 > id2)
			edgeID = id1 + "e" + id2
		else
			edgeID = id2 + "e" + id1
		return edgeID
	}

	static makeEdges(nodeMap) {

		// console.log('edgifying', nodeMap)
		const seenEdges = {}
		nodeMap.edges = []

		for (let i in nodeMap.nodes){

			i = parseInt(i)
			for (let connectionID of nodeMap.nodes[i].connections) {
				let edgeID = NWGraph.makeEdgeID(i, connectionID)

				if (!seenEdges[edgeID]){
					//TODO wtf is this?
					// console.log('edge object', {id:edgeID, source:i, target:connectionID})
					seenEdges[edgeID] = {id:edgeID, source:i, target:connectionID}
					// console.log("seenEdges[edgeID]", seenEdges[edgeID])
					// console.log("seenEdges[edgeID]", seenEdges)
				}
			}
		}

		for (let edgeID of Object.keys(seenEdges)) {
			nodeMap.edges.push(seenEdges[edgeID])
		}
		// console.log("makeEdges produced:",nodeMap)
	}

	static attachRoutes(gameState) {
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
							thisEdgeID = NWGraph.makeEdgeID(route.nodes[i].id, route.nodes[i-1].id)
						} else {
							// otherwise add traffic between last node and target
							thisEdgeID = NWGraph.makeEdgeID(route.nodes[0].id, route.endpoint.id)
						}

						NWGraph.getEdgeIn(thisEdgeID, gameState.map.edges).traffic.push(player)
					}
				}
				// attach endpoints
				gameState.map.nodes[route.endpoint.id].connectedPlayers.push(player)
			}
		}
		// console.log("post attachRoutes", gameState.map)
	}

}

// drawModules() {
	// 	this.nodeGroups.each(function(d) {
	//        	const modules = d3.select(this).selectAll(".node-module") // select instead of selectAll auto binds to parent data
	// 			.data(d.modList)

	// 		// console.log("modules: ", modules)

	// 		const parentRadius = d3.select(this.parentNode).select(".node-main").attr("r");
	// 		const nodeRadius = parentRadius

	// 		// const nodeRadius = (nodeBaseRadius+d.connections.length*nodeRadiusMultiplier)
	// 		const modRad = nodeBaseRadius/2.5
	// 		// const spacing = modRad*3
	// 		const angleInc = 70*0.017453; //convert to radian

	// 		modules.exit()
	// 			.transition(t)
	// 			.attr("r", 0)
	// 			.remove()

	// 		modules.enter()
	// 			   .append("circle")
	// 			   .attr("class", "node-module")
	//    			   .style("fill", d => d.team)
	//    			   .style("fill-opacity", d => {
	//    			   		console.log("making module of fill at:" , d);
	//    			   		return d.health/d.maxHealth
	//    			   	})
	// 	           .style("stroke", "black")
	// 	           .style("stroke-width", 2)
	// 	           .attr("opacity", 0)
	// 	           .attr("r", d => modRad)
	// 	           // .attr("cy", -nodeRadius/2)
	// 	           .transition(t)
	// 	           .style("opacity", 1)
	// 	           .attr("cx", (d,i) => (nodeRadius/2) * Math.cos(-1.5708+angleInc*i))
	// 	           .attr("cy", (d,i) => (nodeRadius/2) * Math.sin(-1.5708+angleInc*i))

	// 	    // update new modules
	// 	    modules.transition(t)
	// 			.style("fill", d => d.team)
	// 		    .style("fill-opacity", d => {
	// 		   		// console.log("making module of fill at:" , d);
	// 		   		return d.health/d.maxHealth
	// 		    })

	// 	})
	// }
