import * as d3 from 'd3'
import React from 'react'
import ReactDOM from 'react-dom'
///////////////////////////////////////////////////////////
/////// Functions and variables

const width = 300;
const height = 250;
const color = d3.scaleOrdinal(d3.schemeCategory10);


const enterNode = (selection) => {
    selection.select('.outer-circle')
        .attr("r", 30)
        .style("fill","white")
        .style("stroke","black")

    selection.select('.inner-circle')
        .attr("r", 23)
        .style("fill","white")
        .style("stroke","black")

    selection.select('text')
        .attr("dy", ".35em")
        .style("transform", "translateX(-50%,-50%")
};

const updateNode = (selection) => {
    selection.attr("transform", (d) => "translate(" + d.x + "," + d.y + ")")

};

const enterLink = (selection) => {
    selection.attr("stroke-width", 2)
    .style("stroke","black")
        // .style("opacity","1")
};

const updateLink = (selection) => {
    selection
        .attr("x1", (d) => d.source.x)
        .attr("y1", (d) => d.source.y)
        .attr("x2", (d) => d.target.x)
        .attr("y2", (d) => d.target.y);
};

const updateGraph = (selection) => {
    selection.selectAll('.node')
        .call(updateNode)
    selection.selectAll('.link')
        .call(updateLink);
};

////////////////////////////////////////////////////////////////////////////
/////// App component. Hold graph data in state and renders Graph component.
/////// Graph component in turn renders Link and Node components.


export class App extends React.Component {
  constructor(props){
    super(props)
    this.state = {
      data : {"nodes":
        [
          {"name": "fruit", "id": 1},
          {"name": "apple", "id": 2},
          {"name": "orange", "id": 3},
          {"name": "banana", "id": 4}
        ],
      "links": 
        [
          {"source": 1, "target": 2},
          {"source": 1, "target": 3},
          {"source": 0, "target": 1},
        ]
      }
    }
  }
  
  handleAdd = () => {

    this.setState({data : {"nodes":
        [
          {"name": "fruit", "id": 1},
          {"name": "apple", "id": 2},
          {"name": "orange", "id": 3},
          {"name": "banana", "id": 4},
          {"name": "banana", "id": 5}
        ],
      "links": 
        [
          {"source": 1, "target": 2},
          {"source": 1, "target": 3},
          {"source": 0, "target": 1},
          {"source": 0, "target": 4},
        ]
      }}, () => alert('done'))
  }
  render() {
    return (
            <div className="graphContainer">
                <Graph data={this.state.data} />
                <button onClick={this.handleAdd}>add</button>
            </div>
        )
    }
}

///////////////////////////////////////////////////////////
/////// Graph component. Holds Link and Node components

export class Graph extends React.Component {

    componentDidMount() {
      this.d3Graph = d3.select(ReactDOM.findDOMNode(this));
      
      var force = d3.forceSimulation(this.props.data.nodes)
        .force("charge", d3.forceManyBody().strength(-50))
        .force("link", d3.forceLink(this.props.data.links).distance(90))
        .force("center", d3.forceCenter().x(width / 2).y(height / 2))
        .force("collide", d3.forceCollide([5]).iterations([5]))

      function dragStarted(d) {
          if (!d3.event.active) force.alphaTarget(0.8).restart()
          d.fx = d.x
          d.fy = d.y

      }

      function dragging(d) {
          d.fx = d3.event.x
          d.fy = d3.event.y
      }

      function dragEnded(d) {
          if (!d3.event.active) force.alphaTarget(0)
          d.fx = null
          d.fy = null
      }

      const node = d3.selectAll('g.node')
        .call(d3.drag()
                  .on("start", dragStarted)
                  .on("drag", dragging)
                  .on("end", dragEnded)
             );
      
        force.on('tick', () => {
            this.d3Graph.call(updateGraph)
        });
    }

    render() {
        var nodes = this.props.data.nodes.map( (node) => {
            return (
            <Node
                data={node}
                name={node.name}
                key={node.id}
            />);
        });
        var links = this.props.data.links.map( (link,i) => {
            return (
                <Link
                    key={link.target+i}
                    data={link}
                />);
        });
        return (
            <svg className="graph" width={width} height={height}>
                <g>
                    {links}
                </g>
                <g>
                    {nodes}
                </g>
            </svg>
        );
    }
}

///////////////////////////////////////////////////////////
/////// Link component

class Link extends React.Component {

    componentDidMount() {
        this.d3Link = d3.select(ReactDOM.findDOMNode(this))
            .datum(this.props.data)
            .call(enterLink);
    }

    componentDidUpdate() {
        this.d3Link.datum(this.props.data)
            .call(updateLink);
    }

    render() {
        return (
             <line className='link' />
        );
    }
}

///////////////////////////////////////////////////////////
/////// Node component

class Node extends React.Component {

    componentDidMount() {
        this.d3Node = d3.select(ReactDOM.findDOMNode(this))
            .datum(this.props.data)
            .call(enterNode)
    }

    componentDidUpdate() {
        this.d3Node.datum(this.props.data)
            .call(updateNode)
    }

    handle(e){
        console.log(this.props.data.id + ' been clicked')
    }

    render() {
        return (
            <g className='node'>
                <circle className="outer-circle" ref="dragMe" onClick={this.handle.bind(this)}/>
                <circle className="inner-circle" onClick={this.handle.bind(this)}/>
                {false ? <text>{this.props.data.id}</text> : null}
            </g>
        );
    }
}

// ReactDOM.render(<App />, document.querySelector('#root'))