import React, { Component } from 'react';
import { PropTypes } from 'prop-types';
import * as d3 from 'd3'

//import { MAX_HEIGHT } from './strings';
const MAX_HEIGHT = 400
let tries = 0

// React Wrapper

function uiButtonHandler(context, icon, d) {
  context.FLAGS.UI[d.label] = !context.FLAGS.UI[d.label]
  d3.select(icon).select('circle').style('fill', context.FLAGS.UI[d.label] ? 'lightpink' : 'white')
  // console.log(d.label + context.FLAGS.UI[d.label])

  d.handler(context, context.FLAGS.UI[d.label])
}

function toggle_traffic(context, flag) {

}

function toggle_power(context, flag) {
  
}

function toggle_production(context, flag) {
  
}

function toggle_alert(context, flag) {
  
}

const UI_ITEMS = [
  {
    label: 'Traffic',
    icon: './icons/ui_traffic.png',
    handler: toggle_traffic,
  },
  {
    label: 'Power',
    icon: './icons/ui_power3.png',
    sizeMod: 1.2,
    handler: toggle_power,
  },
  {
    label: 'Production',
    icon: './icons/ui_production4.png',
    handler: toggle_production,
  },
  {
    label: 'Alert',
    icon: './icons/ui_alert.png',
    handler: toggle_alert,
  },
]


export class UIMenu extends Component {
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

   const SIZES = {
    ui_icon: 15
   }

  this.FLAGS = {
    UI: {}  
  }
  
  for (let ui of UI_ITEMS) {
    this.FLAGS.UI[ui.label] = true
  }

    // Drawing code goes here:
    for (let i = 0; i < UI_ITEMS.length; i++) {
      const thisUI = UI_ITEMS[i]

      const sizeMod = thisUI.sizeMod || 1

      const x = SIZES.ui_icon
      const y = SIZES.ui_icon*2*(i+1)

      
      
      const iconHolder = svg.append('g')
      const self = this

      

      iconHolder.append('circle').attr('class','ui-backing')
        .attr('r', SIZES.ui_icon*.8)
        .style('fill', this.FLAGS.UI[thisUI.label] ? 'lightpink' : 'white')
        .style('stroke-width', 4)  
        .style('stroke', 'black')

      iconHolder.append('image').attr('class','ui-icon')
        .style('x', -SIZES.ui_icon*sizeMod/2)
        .style('y', -SIZES.ui_icon*sizeMod/2)
        .style('width', SIZES.ui_icon*sizeMod)
        .style('height', SIZES.ui_icon*sizeMod)
        .attr("xlink:href",  thisUI.icon)

      const label = iconHolder.append('text').attr('class', 'ui-label')
                .attr('dx', SIZES.ui_icon*1.5)
                .attr("alignment-baseline", "middle")
                .text(thisUI.label)
                .style('opacity', 0)


      iconHolder.attr('transform','translate(' + x + ',' + y + ')')
                .on('click', function(d) { uiButtonHandler(self, this, thisUI) })
                // .on('mouseenter', () => label.transition().duration(500).style('opacity',1))
                // .on('mouseleave', () => label.transition().duration(200).style('opacity',0))

    }
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
