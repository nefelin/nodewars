import React, { Component } from 'react';
import { PropTypes } from 'prop-types';
import * as d3 from 'd3'

//import { MAX_HEIGHT } from './strings';
const MAX_HEIGHT = 400
let tries = 0
const nodeRadius = 60

// React Wrapper

export class PiePower extends Component {
  constructor(props) {
    super(props);
    this.state = {
      svg:null,
    }

    this.update = this.update.bind(this);
    this.checkSVG = this.checkSVG.bind(this);
  }
  componentDidMount() {
    this.update() 
  }

  componentWillReceiveProps() {
    this.update();
  }

  shouldComponentUpdate() {
    return false;
  }

  checkSVG(callback) {
    if (!this.state.svg){

      const svg = d3.select('#graph')
          .attr('width', 400)
          .attr('height', 400)
          .style('border', '1px solid black')

      this.setState({ svg }, callback)
      return 0
    }
    return 1
  }

  update(dataset) {
    if (!this.checkSVG(() => this.update(dataset)))
      return
	 
	 const svg = this.state.svg

    // Drawing code goes here:
    // svg.append('circle')
    //     .attr('r', 20)
    //     .attr('cx', 200)
    //     .attr('cy', 200)

    const data = [1,1,1,1,1,1,1,1]
    
    const context = svg.append('g').attr('transform', 'translate(200,200)')
    // Tell d3.pie where to find slice info in data structure
    const pie = d3.pie()
        .sort(null)
        .value(function(d) { return 1; });

    // Set up arc utility:
    const arc = d3.arc()
        .outerRadius(nodeRadius)
        .innerRadius(nodeRadius*.7);

    // Enter
    const pieGroups = context.selectAll('.pie-group').data(pie(data)).enter().append('g').attr('class', 'pie-group')

    pieGroups.append('path')//.attr('class', 'pie')
      .attr('d', arc)
      .attr('opacity',1)
      .style('fill', d => {
        const coin = Math.floor(Math.random()*3)
        switch (coin) {
          case 0: 
            return "white"
          case 1:
            return "lightpink"
          case 2:
            return "lightblue"
          }
      })
      .style('stroke', 'black')
      .style('stroke-width', 2)
      .attr('class', 'pie-piece')

    const self = this
    svg.selectAll('.pie-piece').each(function(d) {
      if (d3.select(this).style('fill') != 'white' && Math.random() > .3)
        self.makePower(d3.select(this), d)
    })

    // var label = d3.arc()
    // .outerRadius(radius - 40)
    // .innerRadius(radius - 40);

    // d3.csv("data.csv", function(d) {
    //   d.population = +d.population;
    //   return d;
    // }, function(error, data) {
    //   if (error) throw error;

      // var arc = g.selectAll(".arc")
      //   .data(pie(data))
      //   .enter().append("g")
      //     .attr("class", "arc");

      // arc.append("path")
      //     .attr("d", path)
      //     .attr("fill", function(d) { return color(d.data.age); });

      // arc.append("text")
      //     .attr("transform", function(d) { return "translate(" + label.centroid(d) + ")"; })
      //     .attr("dy", "0.35em")
      //     .text(function(d) { return d.data.age; });
    // });
    
  }

  makePower(where, d) {
    const arc = d3.arc()
        .outerRadius(nodeRadius)
        .innerRadius(nodeRadius*.7);

    const icon = './icons/ui_power3.png'
    const SIZES = {
      node_icon: nodeRadius - nodeRadius*.7
    }
    
    const x = arc.centroid(d)[0],
          y = arc.centroid(d)[1]

    // console.log('x',x,'y',y)
    const powerGroup = d3.select(where.node().parentNode).append('g').attr('class','power-group')
    const translate = 'translate(' + x + ',' + y + ')'
    powerGroup.attr('transform', translate)

    powerGroup.append('circle')
             .attr('r', SIZES.node_icon/2)
             .style('fill', "white")
             .style('stroke-width', 2)
             .style('stroke', 'black')

    powerGroup.append('image')
            .style('x', -SIZES.node_icon/2)
            .style('y', -SIZES.node_icon/2)
            .style('width', SIZES.node_icon)
            .style('height', SIZES.node_icon)
            .attr("xlink:href",  icon)

    pulse(powerGroup, translate)
  }

  handleButton = () => {
    this.state.svg.select('circle').remove()
  }

  render() {
    return (
      <div>
        <svg id='graph'/>
        <div>
          <button onClick = {this.handleButton}>Click</button>
        </div>
      </div>
    );
  }
}

function pulse(target, translate) {
  target.transition().duration(700).ease(d3.easeLinear)
        .attr('transform', translate + 'scale(1.3)')
        .transition().duration(700).ease(d3.easeLinear)
        .attr('transform', translate + 'scale(1)')
        .on('end', () => pulse(target, translate))
}