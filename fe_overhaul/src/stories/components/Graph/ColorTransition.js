import React, { Component } from 'react';
import { PropTypes } from 'prop-types';
import * as d3 from 'd3'

//import { MAX_HEIGHT } from './strings';
const MAX_HEIGHT = 400
let tries = 0

// React Wrapper

export class ColorTransition extends Component {
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

    const myscale = d3.scaleLinear().domain([-10,10]).range(['red', 'blue'])
    const myRevScale = d3.scaleLinear().domain([10,-10]).range(['white', 'blue', 'black'])
    // const colorScale = d3.scale.linear().domain([-10, 0, 10]).range(['red', '#ddd', 'green']);

    const data = d3.range(-10,10,.2).map((d,i) => {
      svg.append('rect')
         .attr('width', 40)
         .attr('height', 40)
         .attr('x', (i%10)*40)
         .attr('y', Math.floor(i/10)*40)
         .style('fill', myscale(d))
         .style('stroke', myRevScale(d))
    })
    
    // console.log(data) 

    // Drawing code goes here:
    // const circ = this.state.svg.append('circle')
    //     .attr('r', 20)
    //     .attr('cx', 200)
    //     .attr('cy', 200)
    
    // transLoop(circ)
    
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

function transLoop(item) {
  item.transition().duration(1000)
        .style('fill', d3.hsl(220, .5, .5))
        .transition().duration(1000)
        .style('fill', d3.hsl(360, .5, .5))
        .on('end', () => transLoop(item))
}
