import React, { Component } from 'react';
import { PropTypes } from 'prop-types';
import * as d3 from 'd3'

//import { MAX_HEIGHT } from './strings';
const MAX_HEIGHT = 400
let tries = 0

// React Wrapper

export class SelectionSize extends Component {
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
    svg.selectAll('circle').remove()
    svg.selectAll('rect').remove()

    d3.range(4).map((d,i) => {
      svg.append('circle')
        .attr('r', Math.random()*15+20)
        .attr('cx', Math.random()*300+50)
        .attr('cy', Math.random()*300+50)
    })

    

    const bigBounds = groupBounds(svg.selectAll('circle'))    

    const r = svg.append('rect')
       .attr('x', bigBounds.x)
       .attr('y', bigBounds.y)
       .attr('width', bigBounds.width)
       .attr('height', bigBounds.height)
       .style('fill', 'none')
       .style('stroke', 'black')
    
    console.log(r.size())
  }

  

  handleButton = () => {
    this.state.svg.select('circle').remove()
  }

  render() {
    return (
      <div>
        <svg id='graph'/>
        <div>
          <button onClick = {this.update}>Click</button>
        </div>
      </div>
    );
  }
}

function groupBounds(selection) {
  // compose data
  const boxes = []

  selection.each(function(d) {
   boxes.push(d3.select(this).node().getBBox())
  })

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
