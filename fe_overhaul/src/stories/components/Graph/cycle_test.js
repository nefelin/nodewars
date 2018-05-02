import React, { Component } from 'react';
import { PropTypes } from 'prop-types';
import * as d3 from 'd3'

//import { MAX_HEIGHT } from './strings';
const MAX_HEIGHT = 400
let tries = 0

// React Wrapper

export class CycleTest extends Component {
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

    const group = svg.append('g').attr('transform', 'translate(200, 200)')

    group.append('rect')
         .attr('x', -50)
         .attr('width', 20)
         .attr('height', 20)

     group.append('rect')
         .attr('x', 50)
         .attr('width', 20)
         .attr('height', 20)

    group.transition().duration(1000)
         .attr('transform', 'translate(200, 200) rotate(90)')

  }

  handleButton = () => {
  }

  render() {
    return (
      <div>
        <svg id='graph'/>
        <div>
          <button onClick = {() => this.update(null)}>update</button>
          <button onClick = {this.handleButton}>Click</button>
        </div>
      </div>
    );
  }
}
