import React, { Component } from 'react';
import { PropTypes } from 'prop-types';
import * as d3 from 'd3'

//import { MAX_HEIGHT } from './strings';
const MAX_HEIGHT = 400
let tries = 0

// React Wrapper

export class POEPing extends Component {
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
      // tries++
      // if (tries > 10)
      //   return

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

    console.log(this.state)
    const center = this.state.svg.append('g').attr('translate', "translate(" + 200 + "," + 200 + ")")
    const circ = center.append('circle')
          .attr('r', 20)
          .attr('cx', Math.random()*380)
          .attr('cy', Math.random()*380)
          .style('fill', 'black')
          // .style('stroke', 'black')
          // .style('stroke-width',2)

    pulse(circ)
  }

  handleButton = () => {
    this.state.svg.select('circle').remove()
  }

  render() {
    return (
      <div>
        <svg id='graph'/>
        <button onClick = {this.update}>Create</button>
        <button onClick = {this.handleButton}>Click</button>
      </div>
    );
  }
}


function pulse (target){
  if (target) {
    target.transition().ease(d3.easeLinear).duration(1000)
          .attr('r', target.attr('r')*2)
          .transition().ease(d3.easeLinear).duration(200)
          .attr('r', target.attr('r'))
          .on('end', () => pulse(target))
  }
}