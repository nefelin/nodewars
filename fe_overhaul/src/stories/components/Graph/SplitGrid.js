import React, { Component } from 'react';

import Graph from './Graph'
import * as Maps from '../../maps'
import './SplitGrid.css'
// import '../../../../node_modules/react-grid-layout/css/styles.css'
// import '../../../../node_modules/react-resizable/css/styles.css'

import SplitterLayout from 'react-splitter-layout';
 
class SplitGrid extends React.Component {
  render() {
    const style = {
      border: '1px solid black',
      height: '100%',
      width: '100%',

    }

    return (
      <div>
      <SplitterLayout>
        <div>
          <SplitterLayout vertical={true} primaryMinSize={200} secondaryMinSize={500}>
            <div>
              <SplitterLayout primaryMinSize={300} secondaryMinSize={50}>
                <div style={style}>
                  - terminal
                </div>
                <div style={style}>
                  - score
                </div>
              </SplitterLayout>
            </div>
            <div style={style}>
              <Graph dataset={ Maps.SimpleMap }/>
            </div>
          </SplitterLayout>
        </div>
        <div style={style}>
          Contains:
          - Challenge Details
          - Ace Editor
          - results
          - stdin
        </div>
      </SplitterLayout>
      </div>
    );
  }
}
 
export { SplitGrid }



      

