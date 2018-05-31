import * as d3 from 'd3'
import UI_Toggles from './NWUI_Toggles'
import './NWGraph.css'
import CEQ from 'css-element-queries'
// TODO review the class structure and make those methods that don't rely on instance data STATIC
let log = false

// function newDefaultDict(defValue){
// 	var handler = {
// 	  get: function(target, name) {
// 	    return target.hasOwnProperty(name) ? target[name] : defValue;
// 	  }
// 	};
// 	return new Proxy({}, handler);
// }

const TRANSITIONS = {
	enter_slow: () => d3.transition().duration(1000).ease(d3.easeLinear),
	exit_medium: () => d3.transition().duration(750).ease(d3.easeLinear),
	exit_fast: () => d3.transition().duration(400).ease(d3.easeLinear),
}

const COIN_TIERS = {
	high: {
		cycle_rate: 300,
		cutoff: 3
	},
	med: {
		cycle_rate: 600,
		cutoff: 1
	},
	low: {
		cycle_rate: 1000,
		cutoff: 0
	},
}

function coinTier(value) {
	// console.log('coinTier value', value)
	switch (true) {
		case (value > COIN_TIERS.high.cutoff):
			return "high"
		case (value > COIN_TIERS.med.cutoff):
			return "med"
		case (value > COIN_TIERS.low.cutoff):
			return "low"
		default:
		return "ERROR"
	}
}


// function coinPulseRate(value) {
// 	switch (true) {
// 		case (value > COIN_TIERS.high.cutoff):
// 			return COIN_TIERS.high.cycle_rate
// 		case (value > COIN_TIERS.med.cutoff):
// 			return COIN_TIERS.med.cycle_rate
// 		case (value > COIN_TIERS.low.cutoff):
// 			return COIN_TIERS.low.cycle_rate
// 		default:
// 		return "ERROR"
// 	}
// }

// TODO grab width and height from parent container. Esp on resize
const width = 525,
	  height = 360
	  // NODE_RAD_MULTIPLIER = SIZES.node_outer_radius/4

const SIZES = {
	ui_icon: 15,
	node_icon: 50,
	stroke_width: 2,
	node_outer_radius: 40,
	node_inner_ratio: .7,
	poe_size:16, // node_outer_radius * .4
	power_token: 40-28-1, // node_outer - node_inner - stroke_width/2
	packet_radius: 4,
}

const unfocus_opacity = .1

const TEAMCOLORS = {
	red: d3.hsl(360, 1, .75),
	red_light: d3.hsl(360, 1, .9),
	red_unpowered: d3.hsl(360, .6, .8),
	blue: d3.hsl(230, 1, .75),
	blue_light: d3.hsl(230, 1, .9),
	blue_unpowered: d3.hsl(230, .25, .8),
	// none: "white",
	// none_light: "white", // fix this with || selection TODO
}

const BGCOLOR = "white"

const ICONS = {
		cloak: './icons/feature_cloak.png',
		firewall: './icons/feature_firewall.png',
		overclock: './icons/feature_overclock.png',
		poe: './icons/feature_poe.png',
}

// animation functions
function alertBlip(parent, loc, startR, endR, duration, startColor, endColor) {
	const opac = parent.classed('focused') ? 1 : .1*3

	endColor = endColor || startColor
	loc = loc || [0,0] // if you provide a parent group location probably not necessary

	const signal = parent.append('circle').attr('class','blip')
		.attr('r', startR)
		.attr('cx', loc[0])
		.attr('cy', loc[1])
   		.style('fill', 'none')
	   	.style('stroke', startColor)
	   	.style('stroke-width', SIZES.stroke_width*1.5)
	   	.style('opacity', opac)
	   	

	return signal.transition()
		.ease(d3.easeBounce)
		.duration(duration)
		.attr('r', endR)
		         .transition()
		.ease(d3.easeLinear)
		.duration(duration)
		.attr('r', startR)
		         .transition()
		.ease(d3.easeBounce)
		.duration(duration)
		.attr('r', endR)
		         .transition()
		.ease(d3.easeLinear)
		.duration(duration)
		.attr('r', startR)
		         .transition()
		.ease(d3.easeBounce)
		.duration(duration)
		.attr('r', endR)
				 .transition().ease(d3.easeLinear).duration(500)
		.style('opacity', 0)
		.remove()
}

function poeBlip(targ) { // if target's feature is powered, then we loop.
	const parent = d3.select(targ.node().parentNode)
	// console.log(targ.datum())
	if (targ && targ.datum().owner != "") {
		// console.log('making blip')
		
		const signal = parent.append('circle')
		   	.attr('r', SIZES.poe_size)
	   		.style('fill', 'none')
		   	.style('stroke', 'black')
		   	.style('stroke-width', 2)
		   	.attr('class','poe-signal')
		   	.style('opacity', d3.select(parent.node().parentNode).classed('focused') ? .5 : unfocus_opacity)

		signal.transition()
			.ease(d3.easeLinear)
			.duration(2000)
			.style('opacity', 0)
			.attr('r', SIZES.poe_size*2)
			.remove()
	}
}

function powerAlertPulse(target, transform) {
  target.transition().duration(700).ease(d3.easeLinear)
        .attr('transform', transform + 'scale(1.8)')
        .transition().duration(700).ease(d3.easeLinear)
        .attr('transform', transform + 'scale(1)')
        .on('end', () => powerAlertPulse(target, transform))
}

class NWGraph {
	constructor(targetID) {

		// Make sure container has settled before we draw things
		this.graphDiv = document.querySelector(targetID)
		this.startResizeTimer()

        new CEQ.ResizeSensor(this.graphDiv, ()=> {
        	if (this.resizeTimer)
        		clearTimeout(this.resizeTimer)
        		this.startResizeTimer()
        	
		    // console.log('Changed to ' + this.graphDiv.clientWidth);
		});
	}

	startResizeTimer() {
		this.resizeTimer = setTimeout(() => {
			console.log('div settled')
			this.width = this.graphDiv.clientWidth// + graphOffset[0]
	        this.height = this.graphDiv.clientHeight// + graphOffset[1]
			
			if (this.initialized) {
				this.resize()
			}
			else {
				this.init()
			}
		}, 200)
	}

	init() {
		console.log("width", this.width, "height", this.height)

		console.log("NWGraph creating svg...")
		this.svg = d3.select(this.graphDiv)
		    .append('svg')
		    .attr('width', this.width)
		    .attr('height', this.height)
		    // .attr('width', "100%")
		    // .attr('height', "100%")
		    // .style('border', '3px solid black')
		    // .style('fill', BGCOLOR)
		
		this.background = this.svg.append('g')
		this.background = this.svg.append('rect').classed('background', true)
								  .attr('x', 0)
								  .attr('y', 0)
								  .attr('width', this.width)
								  .attr('height', this.height)
								  .style('fill', BGCOLOR)

		this.root = this.svg.append('g').classed('root', true)
		this.uiLayer = this.svg.append('g').classed('ui-layer', true)
		this.linkLayer = this.root.append('g').classed('link-layer', true)
		this.trafficLayer = this.root.append('g').classed('traffic-layer', true)
		this.nodeLayer = this.root.append('g').classed('node-layer', true)

		this.simulation = d3.forceSimulation()
            .force("link", d3.forceLink().distance(SIZES.node_outer_radius*3))
            .force("collide",d3.forceCollide( function(d){ return d.r + 8 }) )
            .force("charge", d3.forceManyBody().strength(-300))
            .force("center", d3.forceCenter(this.width / 2, this.height / 2))
            .on('end', () => this.layoutStopped(this))
            // .on('start', this.layoutStarted)

		
		// Zoom and Fit
		this.zoom = d3.zoom()
			// .scaleExtent([1/4, 4])
			.on('zoom', () => this.root.attr('transform', d3.event.transform))
		
		const self = this
		this.background.on('click', () => this.resetZoom() )

		// Flags and Animations
        this.FLAGS = {
        	Power: true,
        	Traffic: true,
        	Production: true,
        	Alerts: true,
        	zooming: false,
        }

        

		// console.log('FLAGS', this.FLAGS )

		this.simulation.stop()
        this.stopped = true
		
		// this.FLAGS.traffic_running = false // little hackery needed to start daemon, fix one day TODO
		this.deadTrafficCount = {}
		this.lastTrafficTick = 0
		this.trafficClock = 0
		this.runTrafficDaemon(true)

		const poeInterval = setInterval(() => this.poeDaemon(), 2000)
		this.runCoinPulseDaemon(true)

		

		this.createFills()

		this.drawUI()

		this.initialized = true
		if (this.initState)
			this.update(this.initState)
	}

	zoomFit(zoomFocus, duration) {

		duration = (duration == undefined) ? 750 : duration

		if (this.zoomPending) // supercede any pending zoomings
			clearTimeout(this.zoomPending)
		// console.log('ZOOMFIT before', zoomFocus)
		
		// const zoomFocus = this.nodeGroups.filter('.focused')
		let paddingPercent = .8
		// const nodeCount = zoomFocus.selectAll('.node-group').size()

		const bounds = groupBounds(zoomFocus, this.root)
		// console.log('zoomFit bounds', bounds)

		// // diagnostic
		// if (this.tempBoundBox)
		// 	this.tempBoundBox.remove()
		// this.tempBoundBox = this.root.append('rect')
		// this.tempBoundBox
		//     .attr('x', bounds.x)
		//     .attr('y', bounds.y)
		//     .attr('width', bounds.width)
		//     .attr('height', bounds.height)
		// 	.style('fill', 'none')
		// 	.style('stroke', 'red')
		// 	.style('stroke-width', SIZES.stroke_width)

		// const parent = zoomFocus.node().parentElement;
		// console.log('zoomFit parent', parent)
		const parent = this.svg.node()

		const fullWidth = parent.clientWidth,
		    fullHeight = parent.clientHeight;

		const width = bounds.width,
		    height = bounds.height;

		const midX = bounds.x + width / 2,
		    midY = bounds.y + height / 2;

		if (width == 0 || height == 0) return; // nothing to fit
		const scale = paddingPercent / Math.max(width / fullWidth, height / fullHeight);
		const translate = [fullWidth / 2 - scale * midX, fullHeight / 2 - scale * midY];

		// console.trace("zoomFit", translate, scale);
		this.FLAGS.zooming = true
		this.root
			.transition()
			.duration(duration) // milliseconds
			.call(this.zoom.transform, d3.zoomIdentity.translate(translate[0],translate[1]).scale(scale))
			.on('end', () => this.FLAGS.zooming = false)
			// .on('end', () => console.log('ZOOMFIT after', this.root.node()))
	}

	createFills() {
		const piePattern = this.svg.append("defs")
	        .append('pattern')
	          .attr('id', 'unpowered-pie')
	          .attr('x', 0)
	          .attr('y', 0)
	          .attr('width', 3)
	          .attr('height', 3)
	          .attr('patternUnits', 'userSpaceOnUse')

		// add dots
		piePattern.append('circle')
	           .attr('r', 1)
	           .attr('cx', 1)
	           .attr('cy', 1)
	           .style('stroke', 'none')

	    const featurePattern = this.svg.append("defs")
	        .append('pattern')
	          .attr('id', 'unpowered-feature')
	          .attr('x', 0)
	          .attr('y', 0)
	          .attr('width', 20)
	          .attr('height', 20)
	          .attr('patternUnits', 'userSpaceOnUse')

		featurePattern.append('text')
           .text('!')
           .attr('x', 0)
           .attr('y', 10)
           .style('font-size', 10)

	    featurePattern.append('text')
           .text('!')
           .attr('x', 15)
           .attr('y', 20)
           .style('font-size', 10)
	}

	drawUI() {
		for (let i = 0; i < UI_Toggles.length; i++) {
	      const thisUI = UI_Toggles[i]

	      const sizeMod = thisUI.sizeMod || 1

	      const x = 25
	      const y = SIZES.ui_icon*2*i+25

	      const iconHolder = this.uiLayer.append('g')

	      const backing = iconHolder.append('circle').attr('class','ui-backing')
	        .attr('r', SIZES.ui_icon*.8)
	        .style('fill', this.FLAGS[thisUI.label] ? 'lightpink' : 'white')
	        .style('stroke-width', SIZES.stroke_width)  
	        .style('stroke', 'black')

	      iconHolder.append('image').attr('class','ui-icon')
	        .style('x', -SIZES.ui_icon*sizeMod/2)
	        .style('y', -SIZES.ui_icon*sizeMod/2)
	        .style('width', SIZES.ui_icon*sizeMod)
	        .style('height', SIZES.ui_icon*sizeMod)
	        .attr("xlink:href",  thisUI.icon)

	      const self = this
	      iconHolder.attr('transform','translate(' + x + ',' + y + ')')
	      			// .on('mousedown', function() { console.log('event.target', event.target,'this',this) })//; d3.event.stopPropagation(); })
	                .on('click', () => { 
						this.FLAGS[thisUI.label] = !this.FLAGS[thisUI.label]
						backing.style('fill', this.FLAGS[thisUI.label] ? 'lightpink' : 'white')
						if (thisUI.handler)
							thisUI.handler(this, this.FLAGS[thisUI.label])
	                })
	    }
	}

	

	restartSimulation(alpha = 1) {
		this.simulation.alpha(alpha)
		this.simulation.restart()
		
		this.layoutStarted()
	}

	layoutStarted() {
		// console.log('layout starting')
		this.stopped = false
	}

	layoutStopped(self) {
		// console.log('layout resting')
		// if (!this.FLAGS.zooming)
		// this.trackZoom()
		this.stopped = true
	}

	runTrafficDaemon(bool) {
		if (bool && !this.FLAGS.traffic_running) {
			console.log('starting traffic daemon')
			this.FLAGS.traffic_running = true
			if (this.nodeGroups)
				this.drawTraffic()

			this.trafficSemaphore = this.newSemaphore()

			window.requestAnimationFrame(this.trafficDaemon)

			if (this.traffic)
				this.traffic.transition().duration(500).style('opacity', 1)
		} else {
			console.log('stopping traffic daemon')
			clearInterval(this.trafficSemaphore)
			this.traffic.transition().duration(500)
				.style('opacity', 0)
				.on('end', () => this.FLAGS.traffic_running = false)
				.remove()
		}
	}

	trafficDaemon = (time, other) => {
		const self = this
		// console.log('trafficDaemon e', time, 'other', other)

		// non spammy clock for logging
		const delta = time - this.lastTrafficTick
		this.lastTrafficTick = time
		log = false
		this.trafficClock += delta
		if (this.trafficClock > 2000){
			// console.log(this)
			log = true
			this.trafficClock = 0 
		}


		if (this.traffic){
            this.traffic
				.each(function (d) {
					
					const hideDist = SIZES.node_outer_radius - SIZES.packet_radius
					const totalDist = distance(d.source, d.target)

					const startDistMod = hideDist/totalDist

					// Increment packet's progress
					if (!d.dead) d.progress = progInc(d.progress, .01*d.traffic.length/1.5)

					let source = d.source,
						target = d.target

					// Make sure packet's headed the right way
					const startDir = source.id > target.id ? "down" : "up"
					if (startDir != d.traffic[d.packetNum].dir) {

						source = d.target
						target = d.source
					}

					// Position packet
					const progress = d.progress
					// const ease = d3.easeLinear
					// const progress = ease(d.progress)

					const transX = source.x+(target.x-source.x) * (progress + startDistMod),
				          transY = source.y+(target.y-source.y) * (progress + startDistMod)

					d3.select(this)
					  .attr('transform', 'translate(' + transX + ',' + transY + ')')

					// Check if packet is out of sight
					const p1 = {x: transX, y:transY}

					// console.log('distance', distance(p1,target), 'hideDist', hideDist)
					if (distance(p1,target) < hideDist) {
						d.dead = true

						self.deadTrafficCount[d.traffic.length] = (d.traffic.length in self.deadTrafficCount) ? self.deadTrafficCount[d.traffic.length] + 1 : 1
						// console.log('Dead traffic', self.deadTrafficCount)

						// restart trip
						d.progress = 0
					}
				})
			
		}
		if (this.FLAGS.traffic_running)
			window.requestAnimationFrame(this.trafficDaemon)
	}

	newSemaphore() {
		const self = this
		return setInterval(()=>{
			if (this.traffic){ // if we have traffic packets

				for (let tier of Object.keys(this.deadTrafficCount)){
					const thisTier = this.trafficLayer.selectAll('.tier-'+tier)

					// console.log('Checking tier', tier)

					if (this.deadTrafficCount[tier] >= thisTier.size()) {

						// console.log('All traffic in tier', tier, 'is dead')

						thisTier.each(function(d) {
							if (d.traffic.length > 1){
									if (d.packetNum == d.traffic.length-1)
										d.packetNum = 0
									else 
										d.packetNum++
									// console.log('Semaphore, cycling traffic. PacketNum:', d.packetNum, 'packet:', d.traffic[d.packetNum])
									d3.select(this).select('.traffic-front').style('fill', d => TEAMCOLORS[d.traffic[d.packetNum].owner])
								}
							d.dead = false
							self.deadTrafficCount[tier] = 0
						})
					}
				}
			}
		}, 100)
	}

	poeDaemon = () => {
		if (this.nodeGroups) {
			this.nodeGroups.selectAll('.poe').each(function(d) {
				poeBlip(d3.select(this))
			})
		}
	}

	runCoinPulseDaemon(bool) {
		if (!this.coinPulses)
			this.coinPulses = {}
		if (bool) {
			for (let tierName of Object.keys(COIN_TIERS)) {
				const thisTier = COIN_TIERS[tierName]
				this.coinPulses[tierName] = setInterval(() => this.coinPulseDaemon(tierName), thisTier.cycle_rate+10)
			}
			// this.coinPulse = setInterval(() => this.coinPulseDaemon(), 100)
		} else {
			for (let pulse of Object.values(this.coinPulses)) {
				clearInterval(pulse)
			}
		}
	}

	coinPulseDaemon(tierName) {
		// console.log('coinPulseDaemon')
		if (this.nodeGroups) {
			const thisTier = COIN_TIERS[tierName]
			this.nodeGroups.selectAll('.pie-piece.' + tierName + '-coin').each(function(d) {

				const sel = d3.select(this)
				// if (d.data.powered && d.data.owner in TEAMCOLORS)
					// console.log('pulsing?', this._pulsing)
				if (!this._pulsing && d.data.powered && d.data.owner in TEAMCOLORS) {
					this._pulsing = true
					
					sel.transition('coinPulse').duration(thisTier.cycle_rate).ease(d3.easeLinear)
					   .style('fill', d => TEAMCOLORS[d.data.owner+"_light"])
					   .transition('coinPulse').duration(thisTier.cycle_rate).ease(d3.easeLinear)
					   .style('fill', d => TEAMCOLORS[d.data.owner])
				  	   .on('end', () => {
					 	 	this._pulsing = false
			   		   })
				}
			})
		}
	}

	reset() {
		this.stopped = true

		this.simulation.stop()
		if (this.links)
			this.links.remove()
		if (this.nodeGroups)
			this.nodeGroups.remove()

		if (this.traffic){
			this.traffic.remove()
		}

		this.resize()

		this.gameState = null
	}	

	// TODO resize only works after drag... need to alpha?
	resize() {
		// this.width = this.graphDiv.clientWidth
  //       this.height = this.graphDiv.clientHeight
  		// this.width = newSize.x
  		// this.height = newSize.y


        this.svg.attr('width', this.width)
        	.attr('height', this.height)
        	

       	this.background.attr('width', this.width)
        	.attr('height', this.height)

        console.log("resize width", this.width, "height", this.height)
		
        this.simulation.force("center", d3.forceCenter(this.width / 2, this.height / 2))
        this.restartSimulation(1)
        setTimeout(() => this.trackZoom(), 500)
	}

	draw() {
	// Order of these is important as D3 handles z-index by draw order only
	this.drawLinks()
	this.drawTraffic()
	this.drawNodeGroups()
	}

	update(newState) {

		if (this.initialized){
			console.log("NWGraph updating...")
			console.log('newState:', newState)
			let newMap = this.gameState ? false : true

			if (this.gameState && this.gameState.nodes.length != newState.nodes.length){
				console.log('newState node count mismatch, resetting map before binding new data...')
				newMap = true
				this.reset()
			}

			const self = this

			NWGraph.makeEdges(newState)
			NWGraph.attachTraffic(newState) // puts traffic data onto the newly made edges... 

			// NWGraph.attachPOEs(newState) replacing this with feature-based poe system...
			// NWGraph.attachRoutes(newState)

			// NWGraph.arrayifyModules(newState)
			
			// if we're updating pre-existing state
			if (this.gameState) {
				NWGraph.attachCoords(this.nodeGroups.data(), newState)
			}

			// Adopted newState now that we've done all the prep-work
			this.gameState = newState

			if (this.stopped){
				this.restartSimulation(1)
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
	            .nodes(this.gameState.nodes)
	            .on("tick", ticked);
	    
	        this.simulation.force("link")
	            .links(this.gameState.edges); 
			// console.log("NWGraph state post update:", this.gameState)

			this.draw()

			// show notifications
			if (this.FLAGS.Alerts) {
				if (this.gameState.alerts){
					for (alert of this.gameState.alerts) {
						alertBlip(this.root.select('#node-'+alert.location), undefined, SIZES.node_outer_radius+SIZES.stroke_width, SIZES.node_outer_radius*1.3, 700, TEAMCOLORS[alert.team])
					}
				}
			}

			// Resize
			// console.log('New Map?', newMap)
			if (newMap) {
				for (var i = 0, n = Math.ceil(Math.log(this.simulation.alphaMin()) / Math.log(1 - this.simulation.alphaDecay())); i < n; ++i) {
				    this.simulation.tick();
				  }
				this.resetZoom(0)
			}
		} else {
			this.initState = newState
		}
	}
	
	drawLinks() {
		const update = this.linkLayer.selectAll('.edge').data(this.gameState.edges)

		this.links = update.enter()
            .append("line")
            .attr("class", d => "edge edgeID-" + d.id)
            .attr("stroke", "black")
            .attr("stroke-width", SIZES.stroke_width)
          .merge(update)
            // .classed("traffic", d => d.traffic.length>0)

        update.exit().remove()

        // update traffic info on all links, new and old alike
        // this.links.classed("traffic", d => d.traffic.length>0)
	}

	drawTraffic() {
		// console.log('DRAW TRAFFIC')
		const self = this
		this.links.each(function (d) {
			const thisLink = d3.select(this)

			let traffic
			const trafficSelector = '.traffic-'+d.id

			if (d.traffic.length > 0){

				traffic = self.trafficLayer.selectAll(trafficSelector)

				// transfer any old traffic info if we have it... 
				if (!traffic.empty()){
					d.progress = traffic.datum().progress
					d.packetNum = traffic.datum().packetNum
				} else {
					d.packetNum = 0
					d.progress = 0
				}

				traffic = traffic.data([d])
			} else {
				traffic = self.trafficLayer.selectAll(trafficSelector).data([])
			}
			
			// enter
			traffic.enter().append('g').attr('class', d => 'traffic')
				   .attr('transform', 'translate(' + d.source.x + ',' + d.source.y + ')')
				   .each(function(d) {
				// backing to allow transparency without revealing underlying edge
				d3.select(this).append('circle').attr('class', d => 'traffic-backing')
					.attr('r', SIZES.packet_radius)
					.style('fill', BGCOLOR)
					.style('stroke', BGCOLOR)
					.style('stroke-width', SIZES.stroke_width)
					// .attr('test', d => console.log('TRAFFIC ENTER'))

				d3.select(this).append('circle').attr('class', d => 'traffic-front')
					.attr('r', SIZES.packet_radius)
					.style('stroke', 'black')
					.style('stroke-width', SIZES.stroke_width)
					.style('opacity', () => thisLink.classed('focused') ? 1 : unfocus_opacity)

				})
		
					
			traffic.exit().transition().duration(250).style('opacity',0).remove()
		})
		
		// update 
		this.traffic = this.trafficLayer.selectAll('.traffic')

		this.traffic
			.attr('class', d => 'traffic traffic-'+ d.id + ' tier-' + d.traffic.length)
			.select('.traffic-front')
			.style('fill', d => {
				// ensures when we're losing packets we don't index out of range
				if (d.packetNum > d.traffic.length-1)
					d.packetNum = 0
				return TEAMCOLORS[d.traffic[d.packetNum].owner]
			})
	}

	drawNodeGroups() {
		const self = this
		const update = this.nodeLayer.selectAll(".node-group").data(this.gameState.nodes)

		update.exit().remove()
			


		this.nodeGroups = update.enter()
        	.append("g")
            .attr("class", "node-group focused")
            .attr('id', d => 'node-' + d.id)
            // .style('opacity',.3) //used to see beneath nodes
            .each(function(d) {
            	// draw node Backing
            	const radius = SIZES.node_outer_radius
            	d3.select(this).append('circle').attr('class', 'node-backing backing')
            	  .attr('r', radius*1.07)
            	  .style('stroke', 'none')
            	  .style('fill', BGCOLOR)

            	// draw labels:
            	self.drawNodeLabels(d3.select(this))
            })
            .merge(update)
            .each(function(d) {
            	// const radius = d3.select(this).select('.node-backing').attr('r')
            	const radius = SIZES.node_outer_radius

				// draw node core coloring showing feature control and feature icons				
				self.drawFeatures(this, radius)

            	// draw pie chart showing machine control
				self.drawMachinePie(this, radius)

	        })
	        .on('mousedown', function(e) {
	        	var event = new Event('mousedown', { bubbles: true });
				self.graphDiv.dispatchEvent(event);
	        })
	        .on('click', (clicked) => this.setFocus(clicked.id))
            .call(d3.drag()
			    .on("start", this.dragstarted.bind(this))
			    .on("drag", this.dragged.bind(this))
			    .on("end", this.dragended.bind(this)))

		// indicate of player is present
    	this.drawPlayerLocation(this)            
		self.drawPowerTokens()
		// this.nodeGroups = enter.merge(update)
	}

	drawNodeLabels(node) { // TODO make follow update pattern so that node labels can be changed...
		node.append("text")
	       .attr("dy", function(d) {
	       		const thisPos =  -SIZES.node_outer_radius-18 // doesn't scale or deal with variable sizes
	       		return thisPos
	       	}) 
	       .style('opacity', 0)
	       .attr("text-anchor", "middle")
	       .attr("alignment-baseline", "middle")
	       .attr("font-size", "20")
	       .attr("class", "node-label")
	       .text(d=>d.id)
	       .transition().duration(1000)
	       .style("opacity", 1)
	}
	
	// Viewport tracking functions

	resetZoom(duration) {
		this.resetFocus()
		// should also reset viewport 
		// console.log('root', this.root)
		// console.log('nodeLayer', this.nodeLayer)
		this.zoomFit(this.nodeGroups.filter('.focused'), duration)
	}

	trackZoom() {
		this.zoomFit(this.nodeGroups.filter('.focused'))
	}

	resetFocus() {
		// at the end of resetFocus all nodeGroups, links, and traffic, should be focused and opacity 1
		// console.log('CLEAR nodeGroups size', this.nodeGroups.size())

		this.nodeGroups.classed('focused', true)
		this.links.classed('focused', true)
		
		this.root.selectAll('*').style('opacity', 1)
		
		// clear primary focus
		this.root.select('.focused-primary')
			.classed('focused-primary', false)
			// .classed('focused', false)
	}

	setFocus(id) {
		console.log('Map is focusing on node', id)

		if (this.nodeLayer.select('#node-'+id).classed('focused-primary')) {
			this.resetZoom()
			return
		}

		// resets focus
		this.resetFocus()

		this.nodeLayer.select('#node-'+id).classed('focused-primary', true)

		const focusedIDs = []
		

		// focus nodes
		// console.log('FOCUS nodeGroups size', this.nodeGroups.size())
		this.nodeGroups.each(function(d) {
			if (d.connections.indexOf(id) == -1 && d.id != id) {
				d3.select(this)
				  .classed('focused', false)
				
				// console.log('UNFOCUS', d3.select(this).selectAll('*').size())
				d3.select(this).selectAll(':not(.backing):not(g)')
				   .style('opacity', unfocus_opacity)
				   // .filter('.node-backing')
				   // .style('opacity', 1)

			} else {
				focusedIDs.push(d.id)
				d3.select(this).classed('focused', true)
			}
		})

		// unfocus edges
		this.links.each(function(d) {
			const nodes = d.id.split('e')
			if (focusedIDs.indexOf(parseInt(nodes[0])) == -1 || focusedIDs.indexOf(parseInt(nodes[1])) == -1) {
				d3.select(this)
				  .style('opacity', unfocus_opacity)
				  .classed('focused', false)
			} else {
				d3.select(this)
				  .classed('focused', true)
			}
		})
	
		// unfocus traffic
		this.trafficLayer.selectAll('.traffic').each(function(d) {
			// console.log('unfocus traffic', d3.select(this).data())
			const nodes = d.id.split('e')
			if (focusedIDs.indexOf(parseInt(nodes[0])) == -1 || focusedIDs.indexOf(parseInt(nodes[1])) == -1) {

				d3.select(this).classed('focused', false)
				d3.select(this).select('.traffic-front').style('opacity', unfocus_opacity)

			} else {
				d3.select(this)
				  .classed('focused', true)
			}
		})

		// this.zoomFit(this.root.selectAll('.focused'))
		this.zoomFit(this.nodeGroups.filter('.focused'))
		// groupBounds(this.root.selectAll('.node-group'), this.root)
		// console.log('IDS', focusedIDs)
	}

	drawPlayerLocation() {
		const location = this.gameState.player_location
		// console.log('Player location!', location)
		let node

		if (this.focusBox) {
				this.focusBox.transition().duration(250).style('opacity',0).remove()
			}

		if (location>-1){
			 node = this.root.select('#node-'+location)

			const factor = 2.1
			const radius = node.select('.node-backing').attr('r')
			const size = radius*factor

			this.focusBox = node.append('rect').attr('class','focus-box')
					.attr('stroke-dasharray', [size*.2, size*.6, size*.4, size*.6, size*.4, size*.6, size*.4, size*.6, size*.2])
					.attr('x', -.5*size)
					.attr('y', -.5*size)
					.attr('width', size)
					.attr('height', size)
					.style('fill', 'none')
					.style('stroke-width', SIZES.stroke_width)
					.style('stroke', 'black')
					.style('opacity',0)

			this.focusBox
					.transition().duration(250)
					.style('opacity',1)
		}
	}

	drawFeatures(context, nodeRadius) {
		const self = this
		const data = d3.select(context).datum()
		// console.log('drawFeatures data', data.feature.owner)

		// pre-treat feature data:
		// need type to be none instead of "" because using classed on "" results in true always
		data.feature.type = data.feature.type != "" ? data.feature.type : "none"

		// Enter feature data as group // Entering in this way makes exit never happen but keeps core layered below pie
		// const update = d3.select(context).selectAll('.feature-group').data(data.feature.owner == "" ? [] : [data.feature])
		const update = d3.select(context).selectAll('.feature-group').data([data.feature])
		const enter = update.enter().append('g').attr('class', 'feature-group')

		// shared transition 
		const t = d3.transition().duration(750)

		// append cores
		enter.append('circle').attr('class', 'feature-core')
			// .attr('r', nodeRadius*SIZES.node_inner_ratio-SIZES.stroke_width/2)
			.attr('r', nodeRadius*SIZES.node_inner_ratio)
			.style('stroke', 'none')
			.style('fill', BGCOLOR)
		
		// append icon to enter
		const newIcon = enter.append('image').attr('class','feature-icon ' + data.feature.type)

		newIcon.attr('type', data.feature.type)
		    .style('opacity', 0)
			.style('x', -SIZES.node_icon/2)
			.style('y', -SIZES.node_icon/2)
			.style('width', SIZES.node_icon)
			.style('height', SIZES.node_icon)
			.attr("xlink:href",  d => ICONS[data.feature.type])
		    .transition(t)
			.style('opacity',1)

		// merge enter with old features
		const merged = update.merge(enter)

		
		merged.each(function(mergeD) {
			// update cores
			d3.select(this).select('.feature-core')
			  .attr('class', d => 'feature-core' + (d.powered ? "" : " unpowered"))
			  .transition(t)
		      .style('fill', d => { 
			      	if (d.powered) {
						return TEAMCOLORS[d.owner] || BGCOLOR
					} else {
						// return self.FLAGS.Power ? TEAMCOLORS[d.owner] || BGCOLOR :  TEAMCOLORS[d.owner+"_unpowered"] || BGCOLOR
						return TEAMCOLORS[d.owner+"_unpowered"] || BGCOLOR
					}
				})

			// update icon
			const thisIcon = d3.select(this).select('.feature-icon')
			if (!thisIcon.classed(mergeD.type)) {
				thisIcon.attr('class', 'feature-icon ' + mergeD.type)
					.transition().duration(750)
					.style('opacity', 0)
					.style('x', -SIZES.node_icon/2)
					.style('y', -SIZES.node_icon/2)
					.style('width', SIZES.node_icon)
					.style('height', SIZES.node_icon)
					.on('end', function(d) {
						d3.select(this).attr("xlink:href",  d => ICONS[mergeD.type])
									    .transition().duration(750)
										.style('opacity',1)	
					})
					

			}
		  })

		// Exit
		update.exit().transition(t)
			   .style('opacity', 0)
			   .remove()

	}

	composeMachineData(data) {
		const newData = []
			for (let mac of data.machines){
				// console.log(mac)
				const owner = mac.owner//mac.module != null ? mac.module.owner : "none"
				const size = 1
				// const value = data.remoteness/1//parseInt(data.id)%3 + 1 // Math.ceil(Math.random()*3) // groupData.coinVal TODO implement node value system on backend
				const value = mac.coinval
				// console.log('composeMachineData', data)
				const powered = mac.powered//mac.module != null ? data.powered.indexOf(mac.module.owner) != -1 : true// powered if node is powered by owner
				// console.log('slot:' + slot)
				newData.push({ size, owner, value, powered })
			}
		return newData
	}

	drawMachinePie(context, nodeRadius) {
		// console.log('drawMachinePie context', d3.select(context).data())

		const thisNode = d3.select(context)
		const groupData = d3.select(context).data()[0]
		const data = this.composeMachineData(groupData)
		
		// Tell d3.pie where to find slice info in data structure
		const pie = d3.pie()
		    .sort(null)
		    .value(function(d) { return 1; });

		// Set up arc utility:
		const arc = d3.arc()
		    .outerRadius(nodeRadius)
		    .innerRadius(nodeRadius*SIZES.node_inner_ratio);
		// console.log('macpie noderadius', nodeRadius*SIZES.node_inner_ratio)

		// Enter
		// const piePaths = d3.select(context).selectAll('.pie').data(pie(data))
		const pieGroups = d3.select(context).selectAll('.pie-group').data(pie(data))

		const enter = pieGroups.enter().append('g').attr('class', 'pie-group')

		const newPaths = enter.append('path').attr('class','pie-piece')
							.attr('d', arc)
							.attr('opacity',1)
							.style('fill', BGCOLOR)
							.style('stroke', 'black')
							.style('stroke-width', 2)

		// store last owner to detect change of ownership
		// newPaths.each(function(d) {
		// 	this._owner = d.data.owner
		// })
			
		const t = d3.transition().duration(700)
		
		enter.transition(t)
			.attr("transform", "rotate(180)") // enter in style
		
		const self = this
		
		// Update
		// enter.merge(pieGroups).each(function(d) {
		// 	d3.select(this).select('.pie-piece').attr('class', d => { // this select is necessary to rebind data i believe
		// 		return 'pie-piece' + (d.data.powered ? "" : " unpowered")
		// 	}).each(function(d) {
		// 		let chain = d3.active(this)
		// 		while (chain)
		// 			chain = d3.active(chain)
		// 		chain = chain || d3.select(this)
		// 		console.log('chain', chain)
		// 		chain.transition(TRANSITIONS.enter_slow).style('fill', d => d.data.powered ? TEAMCOLORS[d.data.owner] || BGCOLOR : TEAMCOLORS[d.data.owner+"_unpowered"] || BGCOLOR)
		// 	})
		// })

		enter.merge(pieGroups).each(function(d) {
			d3.select(this).select('.pie-piece').attr('class', d => { // this select is necessary to rebind data i believe
				// console.log('pie update d', d)
				const classes = []
				classes.push('pie-piece') // base class
				classes.push(coinTier(d.data.value) + '-coin') // coin class
				classes.push((d.data.powered ? "" : "unpowered")) // power class

				return classes.join(' ')
			}).each(function(d) {
				// console.log('this owner', this._owner, 'data owner', d.data.owner)
				if (this._owner != d.data.owner || this._powered != d.data.powered) {
					// console.log('status changed!!')
					d3.select(this).interrupt('coinPulse')
					
					d3.select(this).transition(t)
						.style('fill', d => d.data.powered ? TEAMCOLORS[d.data.owner] || BGCOLOR : TEAMCOLORS[d.data.owner+"_unpowered"] || BGCOLOR)
						.on('end', () => this._pulsing = false) // prevents coinPulse from interrupting color transition.

					this._owner = d.data.owner
					this._powered = d.data.powered
				}
			})
		})

		// Exit
		pieGroups.exit()
				.transition().duration(1000)
				.style('opacity', 0)
				.remove()
	}

	drawPowerTokens() {
		const self = this
		
		const t = d3.transition().duration(500)

		if (this.FLAGS.Power) {
			// handle out of power tokens for pie piece
			const selection = this.nodeGroups.selectAll('.pie-piece,.feature-core')//.merge(this.nodeGroups.selectAll('.feature-core'))
			selection.each(function(d) { // for every pie piece / core
				const sel = d3.select(this)
				const parent = d3.select(sel.node().parentNode)
				if (!sel.classed('unpowered')) { 			   // if it has power, remove any out of power tokens
					parent.selectAll('.power-group')
								   .transition(t)
								   .style('opacity', 0)
								   .remove()
				} else if (parent.select('.power-group').empty()) { // if it doesn't have power AND there's no out of power token, make one
					self.makePowerAlert(sel, d, t)
				}
			})
		} else {
			this.nodeGroups.selectAll('.power-group')
				.transition(t)
				.style('opacity', 0)
				.remove()
		}
	}

	makePowerAlert(target, d, t) {
		// needs to differentiate type of selection to know if we have a core or a pie
		// console.log('makePowerAlert target', target.node())

		const parent = d3.select(target.node().parentNode)
		const grandParent = d3.select(target.node().parentNode.parentNode)
		const nodeRadius = grandParent.select('.node-backing').attr('r')
		
		const icon = './icons/ui_power3.png'
		
		const powerGroup = parent.append('g').attr('class','power-group')
		// const prevTransform = parent.attr('transform') || ""
		const backing = powerGroup.append('circle').attr('class', 'power-backing backing')
			.attr('r', SIZES.power_token/2)
			.style('fill', "white")
			.style('stroke-width', 2)
			.style('stroke', BGCOLOR)

		powerGroup.append('circle')
			.attr('r', SIZES.power_token/2)
			.style('fill', 'none')
			.style('stroke-width', 2)
			.style('stroke', 'black')

		powerGroup.append('image')
			.style('x', -SIZES.power_token/2)
			.style('y', -SIZES.power_token/2)
			.style('width', SIZES.power_token)
			.style('height', SIZES.power_token)
			.attr("xlink:href",  icon)
		

		let transform = ""
		if (target.classed('pie-piece')) {
			// must be same arc as in drawMachinePie
			const centroidAdjust = nodeRadius/10 // why do we need this??? TODO can't figure out why this is off center
			const arc = d3.arc()
				.outerRadius(nodeRadius)
				.innerRadius(nodeRadius*SIZES.node_inner_ratio-centroidAdjust); 

			const x = arc.centroid(d)[0],
				  y = arc.centroid(d)[1]

			transform = 'translate(' + x + ',' + y + ') ' //+ prevTransform
		}

		// console.log('makePowerAlert tranform', transform)
		powerGroup.attr('transform', transform)
		
		backing.style('opacity', 0)
				  .transition(t)
				  .style('opacity', 1)

		powerGroup.selectAll(':not(.backing)').style('opacity', 0)
				  .transition(t)
				  .style('opacity', grandParent.classed('focused') ? 1 : unfocus_opacity)
				  .on('end', () => powerAlertPulse(powerGroup, transform))
  }


	getEdge(edgeID) {
		for (let edge of this.gameState.edges) {
			// console.log("getEdge is comparing",edge.id, "and", edgeID)
			if (edge.id == edgeID)
				return edge
		}
		return null 
	}

	

	// Simulation drag helpers -------------------------------------------------------------------

	dragstarted(d) {
        if (!d3.event.active){
        	this.simulation.alphaTarget(0.3)
        	this.restartSimulation()
        }
        d.fx = d.x;
        d.fy = d.y;

        this.dragStart = {x:d.x, y:d.y}
    }
    
    dragged(d) {
        d.fx = d3.event.x;
        d.fy = d3.event.y;
    }
    
    dragended(d) {
    	// console.log('dragended this',  this, 'event target', event.target, 'graphDiv', this.graphDiv)
    	// this.graphDiv.onclick()
        if (!d3.event.active)
        	this.simulation.alphaTarget(0);
        d.fx = null;
        d.fy = null;

		// console.log('dragStart', this.dragStart)

        const distX = d.x-this.dragStart.x,
              distY = d.y-this.dragStart.y

		// console.log('dragEnd d', d)

		const zoomTriggerDist = 60
		// console.log('drag distance', Math.sqrt(distX*distX+distY*distY))
		if (Math.sqrt(distX*distX+distY*distY) > zoomTriggerDist){
			// if (!this.FLAGS.zooming)
			// 	setTimeout(() => this.resetZoom(), 500)
			this.trackZoom()
		}
        
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

	// static arrayifyModules(nodeMap) {
		// for (let node of nodeMap.nodes) {
		// 	node.modList = []
		// 	for (let modID of Object.keys(node.modules)){
		// 		node.modList.push(node.modules[modID])
		// 	}
		// }
	// }

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
	// static attachPOEs(gameState) {
	// 	for (let node of gameState.nodes) {
	// 		node.poes = []
	// 	}

	// 	// console.log("pre attachPoes", gameState.nodes)

	// 	for (let playerID of Object.keys(gameState.poes)) {
	// 		const poeID = gameState.poes[playerID].id
	// 		gameState.nodes[poeID].poes.push(gameState.players[playerID])
	// 	}
	// 	// console.log("post attachPoes", gameState.nodes)
	// }

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

	static attachTraffic(gameState) {
		for (let edge of gameState.edges) {
			if (edge.id in gameState.traffic)
				edge.traffic = gameState.traffic[edge.id]
			else
				edge.traffic = []
		}
	}

	static attachRoutes(gameState) {
		for (let node of gameState.nodes) {
			node.connectedPlayers = []
			node.traffic = []
		}

		for (let edge of gameState.edges) {
			edge.traffic = []
		}

		// console.log("pre attachRoutes", gameState)

		for (let playerID of Object.keys(gameState.players)) {
			const player = gameState.players[playerID]
			const route = player.route

			if (route) {
				// iterate in reverse since routes are reverse ordered
				for (let i = route.nodes.length-1; i > -1; i--) {
					// // attach traffic to nodes UPDATE TODO we don't actually care about node traffic ATM
					const thisNode = route.nodes[i]
					gameState.nodes[thisNode.id].traffic.push(player)

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

						NWGraph.getEdgeIn(thisEdgeID, gameState.edges).traffic.push(player)
					}
				}
				// attach endpoints
				gameState.nodes[route.endpoint.id].connectedPlayers.push(player)
			}
		}
		// console.log("post attachRoutes", gameState)
	}
}

function distance(p1, p2) {
	const a = p1.x-p2.x
	const b = p1.y-p2.y
	return Math.sqrt(a*a + b*b)
}

function progInc(counter, amount) {
	let c = counter
	c += amount
	if (c > 1)
		c = 0
	return c
}

function groupBounds(selection, baseGroup) {
  // compose data
  const boxes = []

  selection.each(function(d) {
  	// console.log('groupBounds selection', d3.select(this).node())
  	const bbox = d3.select(this).node().getBBox()
    boxes.push({x:d.x-bbox.width/2, y:d.y-bbox.height/2-12, width:bbox.width, height:bbox.height}) // odd -12 is to accomodate for text offset... TODO code into variable
  })
  
	// // diagnostic :
	// baseGroup.selectAll('.temp-bounds').remove()
	// boxes.map((bounds,i) => {
 //  	// console.log('drawwing bounds for i: ', i, 'of boxes: ', boxes)
 //  	// console.log('TEST', bounds)

	// baseGroup.append('rect').attr('class', 'temp-bounds')
	//     .attr('x', bounds.x)
	//     .attr('y', bounds.y)
	//     .attr('width', bounds.width)
	//     .attr('height', bounds.height)
	// 	.style('fill', 'none')
	// 	.style('stroke', 'black')
	// 	.style('stroke-width', SIZES.stroke_width)
	//   })

  // find bounds
  let low = {x: undefined, y:undefined}
  let high = {x: undefined, y:undefined}

  boxes.forEach(box => {
    box.x2 = box.x + box.width
    box.y2 = box.y + box.height

    // origin x/y's
    if (low.x == undefined || box.x < low.x)
      low.x = box.x
    if (low.y == undefined || box.y < low.y)
      low.y = box.y

    if (high.x == undefined || box.x > high.x)
      high.x = box.x
    if (high.y == undefined || box.y > high.y)
      high.y = box.y

    // computed x/y's
    if (low.x == undefined || box.x2 < low.x)
      low.x = box.x2
    if (low.y == undefined || box.y2 < low.y)
      low.y = box.y2

    if (high.x == undefined || box.x2 > high.x)
      high.x = box.x2
    if (high.y == undefined || box.y2 > high.y)
      high.y = box.y2

  })
  return {x:low.x, y:low.y, width: high.x-low.x, height: high.y-low.y}
}


export {NWGraph}