
let nodes = null,
	edges = null,
	svg = null
const colors = d3.scale.category10();

function reveal() {
	console.log('force layout settled');
	svg.selectAll("circle, line").transition()
	   .duration(400)
	   .style("opacity", 1)
}

function svgInit() {
	svg = d3.select('#graph')
		    .style("border", "1px solid black")
		    .append('svg')
		    .attr('width', 640)
		    .attr('height', 480) // TODO dynamically pull values, particularly on resize
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
	newNodes[0].weight = 500
	return {nodes: newNodes, edges: newEdges}
}

function initGraphFromNodeMap (nodeMap) {
	nodeMap = arrayifyNodeMap(nodeMap)
	console.log(nodeMap)
	// Calculate force layout, only reveal when done
	const layout = d3.layout.force()
                     .nodes(nodeMap.nodes)
                     .links(nodeMap.edges)
                     .size([640, 480])
                     .linkDistance([60])
                     .charge([-2000])
                     .start()
                     .on('end', reveal)

	// DOM elements for edges
	edges = svg.selectAll("line")
			   .data(nodeMap.edges)
			   .enter()
			   .append("line")
			   .style("stroke","black")
			   .style("stroke-width", 2)

	// DOM elements for nodes
	nodes = svg.selectAll("circle")
			   .data(nodeMap.nodes)
			   .enter()
			   .append("circle")
			   .attr("r", 15)
			   .style("fill", "white")//(d,i) => colors(i))
			   .style("stroke-width", 4)
			   .style("stroke", "black")
			   .call(layout.drag) // Consider non physics trigging dragging TODO

	// everything starts hidden
	svg.selectAll("circle, line")
	   .style("opacity", 0)
	
	// Put these layout movements in the on('end') if you dont want the graph to come swinging into view
    layout.on("tick", function() {
		edges.attr("x1", function(d) { return d.source.x; })
		     .attr("y1", function(d) { return d.source.y; })
		     .attr("x2", function(d) { return d.target.x; })
		     .attr("y2", function(d) { return d.target.y; });

		nodes.attr("cx", function(d) { return d.x; })
		     .attr("cy", function(d) { return d.y; });
	});
	console.log("Calculation Layout...")

}