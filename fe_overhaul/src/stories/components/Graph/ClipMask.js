import React, { Component } from 'react';
import { PropTypes } from 'prop-types';
import * as d3 from 'd3'

//import { MAX_HEIGHT } from './strings';
const MAX_HEIGHT = 400
let tries = 0

// React Wrapper

export class ClipMask extends Component {
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
    // this.state.svg.append('clipPath').attr('id', 'circ')
    //     .append('circle')
    //     .attr('r', 20)
    //     .attr('cx', 200)
    //     .attr('cy', 200)

    // this.state.svg.append('rect')
    //     .attr('x', 200)
    //     .attr('y', 100)
    //     .attr('width', 50)
    //     .attr('height', 50)
    //     .attr('clip-path', 'url(#circ)')
    //     .style('fill', 'grey')

    var pattern = this.state.svg.append("defs")
        .append('pattern')
          .attr('id', 'pattern-test')
          .attr('x', 0)
          .attr('y', 0)
          .attr('width', 20)
          .attr('height', 20)
          .attr('patternUnits', 'userSpaceOnUse')
          .style('stroke','none')

    pattern.append('text')
           .text('!')
           .attr('x', 5)
           .attr('y', 10)
           .style('font-size', 10)

    pattern.append('text')
           .text('!')
           .attr('x', 15)
           .attr('y', 20)
           .style('font-size', 10)

    // pattern.append('circle')
    //        .attr('r', 2)
    //        .attr('cx', 2)
    //        .attr('cy', 2)
    //        .style('stroke', 'none')
    // pattern.append('rect')
    //         .attr('x', 0)
    //         .attr('y', 0)
    //         .attr('width', 1000)
    //         .attr('height', 1000)
    //         .style('fill', 'lightpink')

    // pattern.append('path')
    //          .attr('d', 'M-1,1 l2,-2 M0,10 l10,-10 M9,11 l2,-2')
    //          .style('stroke', 'black')
    //          .style('stroke-width', 1)
    //          .style('fill', 'none');
             // .transition().duration(10000)
             // .attr('transform','rotate(180)')

    //Shape design
    this.state.svg.append("g").attr("id","shape")
        .attr('transform', 'translate(200,200)')
        .append("circle")
        .attr('r', 50)
        .style('fill', 'url(#pattern-test)')
        .transition().duration(5000)
        
        // .style('fill', 'black')
        // .attr({cx:"200",cy:"200",r:"50" })

  svg.append('text')
     .text('hello world')
     .attr('x', 5)
     .attr('y', 5)
  // svg.append('path')
  //            .attr('d', 'M0,0' + ' L100,100')
  //            .style('stroke', 'black')
  //            .style('stroke-width', 2)
  //            .style('fill', 'none');
    
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
