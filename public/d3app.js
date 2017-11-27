const width = 640, height = 480
let node = null,
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
function arrayifyNodeMap (nodeMap) {
	newNodes = []
	newEdges = []
	for (key of Object.keys(nodeMap.nodes)){
		// console.log(nodeMap.nodes[key])
		newNodes.push(nodeMap.nodes[key])
	}
	for (key of Object.keys(nodeMap.edges)){
		// console.log(nodeMap.nodes[key])
		newEdges.push(nodeMap.edges[key])
	}

	return {nodes: newNodes, edges: newEdges}
}

function attachCoords(nodeMap) {
	const data = node.data()

	// console.log('nodeMap.nodes before:', nodeMap.nodes)
	for (let i=0; i<data.length; i++) {
		nodeMap.nodes[i].x = data[i].x
		nodeMap.nodes[i].y = data[i].y
	}
	// console.log('nodeMap.nodes after:', nodeMap.nodes)
}

function updateGraph (nodeMap) {
	attachCoords(nodeMap)
	console.log('updateGraph')

	


	// Since we are updating the data with new objects,
	// we need to point our simulation at the new objects 
	// to ensure continued tracking:
	simulation.nodes(nodeMap.nodes)
	simulation.force("link")
            .links(nodeMap.edges);

	console.log("d3 node", node)
	node.data(nodeMap.nodes)
		.style("fill", d => {
		 	if (d.poe.length > 0) {
		 		// console.log("d.poe",d.poe)
		 		if (d.poe[0].team != null)
		 			return "red"
		 			return d.poe[0].team.name
		 	}
		 	return "green"
		 })

	link.data(nodeMap.edges)

// 	// update node and edge traffic
		 
// 	// update player POEs and ongoing connections

// 	// update module contents

}

// initGraph creates the actual dom elements and provides necessary class tags etc
// as well as setting up rules for interactivity (i.e. zoom, drag, etc)
// updateGraph is responsible for actually mapping game data to individual nodes/edges
function initGraph (nodeMap) {
		console.log('Initializing Graph...');

		simulation = d3.forceSimulation()
            .force("link", d3.forceLink().id(function(d) { return d.index }))
            .force("collide",d3.forceCollide( function(d){return d.r + 8 }).iterations(16) )
            .force("charge", d3.forceManyBody().strength(-200))
            .force("center", d3.forceCenter(width / 2, height / 2))
            .force("y", d3.forceY(0))
            .force("x", d3.forceX(0))
    
        link = svg.append("g")
            .attr("class", "links")
            .selectAll("line")
            .data(nodeMap.edges)
            .enter()
            .append("line")
            .attr("stroke", "black")
        
        node = svg.append("g")
            .attr("class", "nodes")
            .selectAll("circle")
            .data(nodeMap.nodes)
            .enter().append("circle")
            .attr("r", function(d){  return 10 })
            // .call(d3.drag()
            //     .on("start", dragstarted)
            //     .on("drag", dragged)
            //     .on("end", dragended));    
        
        
        var ticked = function() {
            link
                .attr("x1", function(d) { return d.source.x; })
                .attr("y1", function(d) { return d.source.y; })
                .attr("x2", function(d) { return d.target.x; })
                .attr("y2", function(d) { return d.target.y; });
    
            node
                .attr("cx", function(d) { return d.x; })
                .attr("cy", function(d) { return d.y; });
        }  
        
        simulation
            .nodes(nodeMap.nodes)
            .on("tick", ticked);
    
        simulation.force("link")
            .links(nodeMap.edges);  

	console.log("Calculation Layout...")

}