import React, { Component } from 'react';
import { PropTypes } from 'prop-types';
import { select } from "d3-selection";
import 'd3-transition';
import {NWGraph} from "./NWGraph"
import './Graph.css'

import * as Maps from '../../maps'

// React Wrapper

class Graph extends Component {
  constructor(props) {
    super(props);
    this.state = {
      graph:null,
      graphOffset: props.graphOffset || [0,0]
    }
    // this.update = this.update.bind(this);
  }
  componentDidMount() {
    // console.log('did mount size', this.props.size)
    if (this.props.dataset)
      this.update(this.props.dataset);
    // this.resize()
  }
  componentWillReceiveProps(newProps) {
    // console.log('will receive size', this.props.size)
    if (newProps.dataset)
      this.update(newProps.dataset);
  }

  shouldComponentUpdate() {
    return false;
  }

  initGraph = (callback) => {
    const graph = new NWGraph('#graph', this.state.graphOffset)
    this.setState({
      graph
    }, callback)
  }
  
  update = (dataset) => {
    if (!this.state.graph){
      this.initGraph(() => this.update(dataset))
      return
    }
    // console.log('graph has been initialized?', this.state.graph)

    console.log('dataset', dataset)
    this.state.graph.update( dataset)
    // this.state.graph.draw()
  }

  reset = () => {
    if (this.state.graph)
      this.state.graph.reset()
  }

  resize(newSize) {
    if (this.state.graph)
      this.state.graph.resize(newSize)
  }

  // call(f, argument) {
  //   console.log('testing ref function calls')
  //   if (this.state.graph) {
  //     console.log('calling', f, 'on graph')
  //     this.state.graph[f](argument)
  //   } else {
  //     console.log('graph not initialized yet...')
  //   }
  // }

  render() {
    if (this.props.debug) {
      return (
          <div style={{height:400, width:400, border:'1px solid black'}}>
            <div id='graph'/>
            {Object.keys(Maps).map((mapName) => {
              return <button key={mapName} onClick={() => this.update(Maps[mapName])} > { mapName } </button>  
            })}
            
          </div>
      );
    } else {
      return (
          <div style={{boxSizing: 'border-box'}} id='graph'/>
      )
    }
  }
}

// Graph.propTypes = {
//   dataset: PropTypes.object.isRequired,
// }

export default Graph