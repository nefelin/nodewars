import React, { Component } from 'react';
import GridLayout from 'react-grid-layout';

import Graph from './Graph'
import * as Maps from '../../maps'
// import '../../../../node_modules/react-grid-layout/css/styles.css'
// import '../../../../node_modules/react-resizable/css/styles.css'


class GridTest extends React.Component {
  render() {
    // layout is an array of objects, see the demo for more complete usage
    var layout = [
      // {i: 'a', x: 0, y: 0, w: 1, h: 2, static: true},
      {i: 'terminal', x: 0, y: 0, w: 4, h: 4, minW: 4, maxW: 8},
      {i: 'score', x: 4, y: 0, w: 2, h: 4}, //isResizable: false,
      {i: 'map', x: 0, y: 4, w: 6, h: 15, minW: 6, maxW: 8, minH:12, maxH:15},
    ];

    const style = {
      border: '1px solid black',
      overflow: 'hidden'

    }
    return (
      <GridLayout style={style}  layout={layout} cols={12} margin={[5,5]} rowHeight={30} width={1200} height={1000}>
        <div style={style} key="map"><Graph dataset={ Maps.SimpleMap } onClick={() => alert('test')}/></div>
        <div style={style} key="terminal">b</div>
        <div style={style} key="score">b</div>
      </GridLayout>
    )
  }
}

// import SplitterLayout from 'react-splitter-layout';
 
// class GridTest extends React.Component {
//   render() {
//     const style = {
//       border: '1px solid black'
//     }

//     return (
//       <div style={style}>
//       <SplitterLayout style={style}>
//         <div>
//           <SplitterLayout vertical={true} primaryMinSize={200} secondaryMinSize={500}>
//             <div>
//               <SplitterLayout primaryMinSize={300} secondaryMinSize={50}>
//                 <div>
//                   - terminal
//                 </div>
//                 <div>
//                   - score
//                 </div>
//               </SplitterLayout>
//             </div>
//             <div>
//               <Graph dataset={ Maps.SimpleMap }/>
//             </div>
//           </SplitterLayout>
//         </div>
//         <div>
//           Contains:
//           - Challenge Details
//           - Ace Editor
//           - results
//           - stdin
//         </div>
//       </SplitterLayout>
//       </div>
//     );
//   }
// }
 
export { GridTest }



      

