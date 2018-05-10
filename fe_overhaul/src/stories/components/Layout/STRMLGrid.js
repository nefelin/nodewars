import React, { Component } from 'react';

import GridLayout from 'react-grid-layout';
import { Responsive, WidthProvider } from 'react-grid-layout';
import './STRMLGrid.css'

import AceEditor from 'react-ace'
import 'brace/mode/golang';
import 'brace/theme/monokai';

// const ResponsiveGridLayout = WidthProvider(Responsive);


import '../../../../node_modules/react-grid-layout/css/styles.css'
import '../../../../node_modules/react-resizable/css/styles.css'


import TinyTerm from '../Terminal/TinyTerm'

import TestResults from '../TestResults'

import { defaultLayout } from './Layouts'

import STRMLWindow from './STRMLWindow'

import NWSocket from '../NWSocket'
import * as Parsers from '../Parsers'
// map stuff
import Graph from '../Graph/Graph'
import * as Maps from '../../maps'

class STRMLGrid extends React.Component {
  constructor(props) {
    super(props)

    this.graph = React.createRef()
    this.state = {
      mapSize: {x:5*105, y:12*30} ,
      map: Maps.SimpleMap,
      
      stdin: 'Put your own stdin here',
      aceContent: "Code Here",

      // tinyTermIn: "",

      compilerOutput: null,
      testResults: {stdouts:[], grades:[], hints:[]},

      layout: this.readLocalLayout() || defaultLayout,

      // graph: null
    }
    
    this.terminal = React.createRef()
    this.editor = React.createRef()
    this.graph = React.createRef()

    this.termHasFocus = true
  }
  

  componentDidMount() {
    // Setup for websocket + handlers
    const inParser = new Parsers.Incoming(this, true)
    const ws = new NWSocket(inParser)
    this.outgoing = new Parsers.Outgoing(ws, this, true)

    // Set up keyboard shortcuts
    
    document.onkeydown = this.handleKeyPress;
  }

  handleKeyPress = (e) => {
      var evtobj = window.event ? event : e
      // console.log('keypress', e.keyCode)
      if (evtobj.ctrlKey) {
        switch(evtobj.keyCode) {
          case 192: // '`'
            // console.log('switch focus')
            this.toggleFocus()
            this.gatherFocus()
            break

          case 188: // ','
            this.toggleFocus('terminal')
            this.gatherFocus()
            break

          case 190: // '.'
            this.toggleFocus('editor')
            this.gatherFocus()
            break

          case 13: // 'enter'
            this.outgoing.parseCmd('make')
            break

          case 220: // '\'
            this.outgoing.parseCmd('test')
            break
          
           case 82: // 'r'
            this.outgoing.parseCmd('reset')
            break

          default:
            
        }
      }
  }

  toggleFocus = (target) => {
    // console.log('termhasfocus', this.termHasFocus)
    switch (target) {
      case 'terminal':
        this.termHasFocus = true
        break
      case 'editor':
        this.termHasFocus = false
        break
      default:
        this.termHasFocus = !this.termHasFocus
    }
  }

  gatherFocus = () => {
    setTimeout(()=>{
      if (this.termHasFocus){
        this.terminal.current.focus()
      }
      else {
        this.editor.current.editor.focus()
      }
    }, 0)
  }

  handleTermSend = (cmd) => {
    this.outgoing.parseCmd(cmd) 
  }

  handleChange = (e) => {
    // track current layout
    this.currentLayout = JSON.stringify(e)
  }

  handleAceChange = (val) => {
    this.setState({ aceContent: val },
      () => this.outgoing.componentState("editor", this.state.aceContent)
      )
  }

  handleStdinChange = (e) => {
    this.setState({ stdin: e.target.value },
      () => this.outgoing.componentState("stdin", this.state.stdin)
      )
  }

  handleSelect = (win, menu, item) => {
    switch (win) {
      case 'NodeWars':
        switch (menu){
          case 'Layout':
            switch(item){
              case 'Save':
                this.saveLayout()
                break

              case 'Load':
                this.loadLayout()
                break

              case 'Reset':
                this.resetLayout()
                break
            }
        }
        break
      
      case 'Ace Editor':
        switch (menu) {
          case 'Build':
            switch (item){
              case 'Make, (ctrl-enter)':
                this.outgoing.parseCmd('make')
                break
              case 'Test, (ctrl-\\)':
                this.outgoing.parseCmd('test')
                break
              case 'Reset, (ctrl-r)':
                this.outgoing.parseCmd('reset')
                break
            }
            break
        }
        break
      default:
        console.log('Routing for', win, menu, item, 'missing')
    }
  }

  readLocalLayout = () => {
    let layout = localStorage.getItem('gridLayout')
    if (layout){
      console.log('layout cookie loaded')
      return layout = JSON.parse(layout)
    } else {
      console.log('error: No layout cookie found!')
      return null 
    }
  }

  saveLayout = () => {
    localStorage.setItem('gridLayout', this.currentLayout)
  }

  resetLayout = () => {
    this.setState({ layout: defaultLayout })
  }

  loadLayout = () => {
    const layout = this.readLocalLayout()
    if (layout)
      this.setState({layout})
  }


  render() {

    const layoutProps = {
      onLayoutChange: this.handleChange,
      className: 'layout',
      draggableCancel: 'input,textarea,.strml-menu-title',
      draggableHandle: '.strml-window-header',
      layout: this.state.layout,
      cols: 12,
      rowHeight: 30,
      width: 1260,
    }
    
    const fill = {
      width: '100%',
      height: '100%',
      boxSizing: 'border-box',
      resize: 'none'
    }

    return (
      // <GridLayout style={strmlWindow} onLayoutChange={this.handleChange} className="layout" draggableCancel="input,textarea" layout={layout} cols={12} rowHeight={30} width={1260}>
      <div  style={{height: 1500}}>

      <STRMLWindow bgOverride={this.state.team ? this.state.team : 'white'} className="full-screen" menuBar={[{ name: 'NodeWars' }, {name: 'Layout', items: ['Load', 'Save', 'Reset']}]} onSelect={this.handleSelect}>
        <GridLayout {...layoutProps}>

          <div key="terminal" >
            <STRMLWindow menuBar={[{ name: 'Terminal' }]} onSelect={this.handleSelect}>
              <TinyTerm ref={this.terminal}  onFocus={() => this.toggleFocus('terminal')}  grabFocus={true} onSend={this.handleTermSend}/>
            </STRMLWindow>
          </div>

          <div key="map">
            <STRMLWindow menuBar={[{ name: 'Map' }, { name: 'Theme', items: ['Light', 'Dark'] }]} onSelect={this.handleSelect}>
              <Graph onMouseDown={() => this.gatherFocus()} ref={this.graph} dataset={ this.state.graph }/>
            </STRMLWindow>
          </div>
          
          <div key="score">
            <STRMLWindow menuBar={[{ name: 'Score' }]} onSelect={this.handleSelect}>
            </STRMLWindow>
          </div>

          <div key="codepad">
            <STRMLWindow menuBar={[{ name: 'Ace Editor' }, { name: 'Build', items: ['Make, (ctrl-enter)', 'Test, (ctrl-\\)', 'Reset, (ctrl-r)']}]} onSelect={this.handleSelect}>
              <AceEditor
                  ref={this.editor}
                  style={fill}
                  mode="golang"
                  theme="monokai"
                  onFocus={() => this.toggleFocus('editor')}
                  onChange={this.handleAceChange}
                  name="UNIQUE_ID_OF_DIV"
                  editorProps={{$blockScrolling: true}}
                  value={this.state.aceContent}
                />
            </STRMLWindow>
          </div>

          <div key="challenge_details">
            <STRMLWindow menuBar={[{ name: 'Challenge Details' }]} onSelect={this.handleSelect}>
              <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit. Mollitia officia, illo magni. Consequatur sapiente adipisci eos, fugit maxime velit necessitatibus corporis illo ut molestiae et temporibus ipsum quas voluptatum deleniti.</p>
            </STRMLWindow>
          </div>

          <div key="test_results">
            <STRMLWindow menuBar={[{ name: 'Test Results' }]} onSelect={this.handleSelect}>
              <TestResults results={this.state.testResults}/>
            </STRMLWindow>
          </div>

          <div key="compiler_output">
            <STRMLWindow menuBar={[{ name: 'Compiler Output' }]} onSelect={this.handleSelect}>
              <div style={{margin:10}}>
                <p>
                  {
                    this.state.compilerOutput ?
                      this.state.compilerOutput.message.split('\n').map((line, key) => {
                        // console.log('key', key, 'length', pageLines.length)
                        return <span key={key}>{line}<br/></span>      
                      })
                      : null
                  }
                </p>
              </div>
            </STRMLWindow>
          </div>

          <div key="stdin">
            <STRMLWindow menuBar={[{ name: 'Stdin' }]} onSelect={this.handleSelect}>
              <textarea style={fill} onChange={this.handleStdinChange} value={this.state.stdin}></textarea>
            </STRMLWindow>
          </div>


        </GridLayout>
      </STRMLWindow>
      </div>
    )
  }
}
export { STRMLGrid }



        /*<div style={strmlWindow} key="terminal">
          <div className="grid-item-header" style={header}>Termnal</div>
            <div style={content}>
              <TinyTerm/>
            </div>
        </div>
        <div style={strmlWindow} key="map">
          <div className="grid-item-header" style={header}>Map</div>
          <div style={content}>
          <Graph ref={this.graph} graphOffset={[0, -header.height]} dataset={ Maps.SimpleMap }/>
          </div>
        </div>
        <div style={strmlWindow} key="score">
          <div className="grid-item-header" style={header}>Score</div>
        </div>
        <div style={strmlWindow} key="codepad">
          <div className="grid-item-header" style={header}>Ace</div>
            <div style={content}>
              <AceEditor
                mode="golang"
                theme="monokai"
                onChange={this.handleAceChange}
                name="UNIQUE_ID_OF_DIV"
                editorProps={{$blockScrolling: true}}
                value={this.state.aceContent}
              />
            </div>
        </div>
        <div style={strmlWindow} key="challenge_details">
          <div className="grid-item-header" style={header}>Challenge Details</div>
            <div style={{margin:10}}>
              <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit. Mollitia officia, illo magni. Consequatur sapiente adipisci eos, fugit maxime velit necessitatibus corporis illo ut molestiae et temporibus ipsum quas voluptatum deleniti.</p>
            </div>
        </div>
        <div style={strmlWindow} key="test_results">
          <div className="grid-item-header" style={header}>Test Results</div>
            <div style={content}>
              <TestResults results={results}/>
            </div>
        </div>
        <div style={strmlWindow} key="compiler_output">
          <div className="grid-item-header" style={header}>Compiler Output</div>
          <div style={{margin:10}}>{this.state.compilerOutput}</div>
        </div>
        <div style={strmlWindow} key="stdin">
          <div className="grid-item-header" style={header}>Stdin</div>
          <textarea style={stdin} onChange={this.handleStdinChange} value={this.state.stdin}></textarea>
        </div>



      

*/