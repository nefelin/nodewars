import React, { Component } from 'react';
import { PropTypes } from 'prop-types';
import * as d3 from 'd3'
import CEQ from 'css-element-queries'

//import { MAX_HEIGHT } from './strings';
const MAX_HEIGHT = 400
let tries = 0

// React Wrapper

export default class ScoreBars extends Component {
  constructor(props) {
    super(props);
    this.state = {
      svg:null,
    }

    this.update = this.update.bind(this);
    this.makeSVG = this.makeSVG.bind(this);
    this.containerDiv = React.createRef()
  }

  componentDidMount() {
    // this.update()
    // new CEQ.ResizeSensor(this.containerDiv.current, () => console.log(this.containerDiv.current.clientWidth))
    new CEQ.ResizeSensor(this.containerDiv.current, ()=> {
          if (this.resizeTimer)
            clearTimeout(this.resizeTimer)
            this.startResizeTimer()
          
        // console.log('Changed to ' + this.graphDiv.clientWidth);
    });
  }

  startResizeTimer() {
    this.resizeTimer = setTimeout(() => {
      console.log('score div settled')
      this.width = this.containerDiv.current.clientWidth// + graphOffset[0]
          this.height = this.containerDiv.current.clientHeight// + graphOffset[1]
      
      if (this.svg) {
        this.resize()
        this.update()
      }
      else {
        this.makeSVG(this.update())
      }
    }, 400)
  }

  componentWillReceiveProps() {
    this.update();
  }

  shouldComponentUpdate() {
    return false;
  }

  resize() {
    // resize everything.

    this.margin = {
        top: 0,//this.height*.1,
        left: 0,//this.width*.5,
        bottom: 0,//this.height*.,
        right: 0//this.width*.001
      }
    // this.height = this.containerDiv.current.clientHeight
    // this.width = this.containerDiv.current.clientWidth

    // const totalWidth = this.width + this.margin.left + this.margin.right,
    //       totalHeight = this.height + this.margin.top + this.margin.bottom

    this.widthMargins = this.margin.left + this.margin.right,
    this.heightMargins = this.margin.top + this.margin.bottom,



    this.svg = d3.select("#score")
        .attr("width", this.width)
        .attr("height", this.height)
    
    this.root.attr("transform", 
              "translate(" + this.margin.left + "," + this.margin.top + ")");

    // redraw bars
    // this.update()
    // console.log('SHOULD RESIZE')
  }

  makeSVG(callback) {
    if (!this.svg){

    this.root = d3.select("#score").append("g")
    this.svg = d3.select("#score")
        .attr("width", this.width)
        .attr("height", this.height)
    this.resize()
        // .style('border', '1px solid black')
      // .append("g")
        // .attr("transform", 
        //       "translate(" + this.margin.left + "," + this.margin.top + ")");
      // // const svg = d3.select('#graph')
      //     .attr('width', 400)
      //     .attr('height', 400)
      //     .style('border', '1px solid black')

      // this.setState({ svg }, callback)
      if (callback)
        callback()
      return 0
    }
    return 1
  }

  update(override) {
    if (!this.makeSVG(() => this.update()))
      return
    const dummyData = [
            { name: 'R', color: 'black', cc: 700 },
            { name: 'B', color: 'grey', cc: 550 },
          ]
    const data = teamSort(override || this.props.data || dummyData)

    

    const maxScore = 1000
  
    const t = d3.transition().duration(1000).ease(d3.easeLinear)
    const [x, y] = genScales(this.width-this.widthMargins, this.height-this.heightMargins)
    y.domain([0, maxScore])

    // Drawing code goes here:
    if (data && data.length > 0) {
      
      
      x.domain(data.map(d => d.name))
      

      const update = this.svg.selectAll('.bar')
          .data(data)

      const enter = update.enter().append('rect')
        .attr('class', 'bar')

      enter.merge(update)
        .transition(t)
        .attr('x', d => x(d.name))
        .attr('width', x.bandwidth())
        .attr('y', d => y(d.cc))
        .attr('height', d => this.height - y(d.cc))
        .style('fill', d => d.name)
    }

    // // add the x Axis
    // this.svg.append("g")
    //     .attr("transform", "translate(0," + this.height + ")")
    //     .call(d3.axisBottom(x));

    // // add the y Axis
    // this.svg.append("g")
    //     .call(d3.axisLeft(y));


    // svg.append('circle')
    //     .attr('r', 20)
    //     .attr('cx', 200)
    //     .attr('cy', 200)
    
  }

  handleButton = () => {
    this.update([
            { name: 'R', color: 'red', cc: 110 },
            { name: 'B', color: 'blue', cc: 75 },
            { name: 'G', color: 'green', cc: 72 },
          ])
  }

  render() {
    return (
      <div style={{width:'100%', height:'100%'}} ref={this.containerDiv}>
        <svg id='score'/>
      </div>
    );

  }
}


function genScales(width, height) {
   const scaleX = d3.scaleBand()
          .range([0, width])
          .padding(0.1);
   const scaleY = d3.scaleLinear()
          .range([height, 0]);

  return [scaleX, scaleY]
}

function teamSort(data) {
  function comparitor(a, b) {
    if (a.name[0] > b.name[0])
      return 0
    return 1
  }

  return data.sort(comparitor)
}