import React, { Component } from 'react';

import GridLayout from 'react-grid-layout';
import { Responsive, WidthProvider } from 'react-grid-layout';
import './STRMLGrid.css'

import AceEditor from 'react-ace'
import { modeMap, themeMap } from './brace-modes-themes.js'

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

      language: 'python',
      aceTheme: 'monokai',

      compilerOutput: null,
      testResults: {stdouts:[], grades:[], hints:[]},

      layout: this.readLocalLayout() || defaultLayout,

      // graph: null
    }

    window.onfocus = () => this.gatherFocus()

    this.terminal = React.createRef()
    this.editor = React.createRef()
    this.stdin = React.createRef()
    this.graph = React.createRef()

    this.focusedField = "TERM"
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
        this.focusedField = "TERM"
        break
      case 'editor':
        this.focusedField = "EDIT"
        break
      case 'stdin':
        this.focusedField = "STDIN"
        break
      default: // toggles between term and edit
        if (this.focusedField == "TERM")
          this.focusedField = "EDIT"
        else
          this.focusedField = "TERM"
    }
  }

  gatherFocus = (e) => {
    // console.log(e)
    // if (e){
    //   // e.preventDefault()
    //   e.stopPropagation()
    // }
    setTimeout(()=>{
      switch (this.focusedField){
        case "TERM":
        this.terminal.current.focus()
        break

        case "EDIT":
        this.editor.current.editor.focus()
        break

        case "STDIN":
        this.stdin.current.focus()
        break
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
          case 'Theme':
            this.setState({ aceTheme: themeMap[item] })
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
      outline: 'none',
      width: '100%',
      height: '100%',
      boxSizing: 'border-box',
      resize: 'none'
    }

    return (
      // <GridLayout style={strmlWindow} onLayoutChange={this.handleChange} className="layout" draggableCancel="input,textarea" layout={layout} cols={12} rowHeight={30} width={1260}>
      <div  style={{height: 1500}}>

      <STRMLWindow onMouseDown={this.gatherFocus}  bgOverride={this.state.team ? this.state.team : 'white'} className="full-screen" menuBar={[{ name: 'NodeWars' }, {name: 'Layout', items: ['Load', 'Save', 'Reset']}]} onSelect={this.handleSelect}>
        <GridLayout {...layoutProps}>

          <div key="terminal" >
            <STRMLWindow onMouseDown={this.gatherFocus} menuBar={[{ name: 'Terminal' }]} onSelect={this.handleSelect}>
              <TinyTerm ref={this.terminal}  onFocus={() => this.toggleFocus('terminal')}  grabFocus={true} onSend={this.handleTermSend}/>
            </STRMLWindow>
          </div>

          <div key="map">
            <STRMLWindow onMouseDown={this.gatherFocus} menuBar={[{ name: 'Map' }]} onSelect={this.handleSelect}>
              <Graph ref={this.graph} dataset={ this.state.graph }/>
            </STRMLWindow>
          </div>
          
          <div key="score">
            <STRMLWindow onMouseDown={this.gatherFocus} menuBar={[{ name: 'Score' }]} onSelect={this.handleSelect}>
            </STRMLWindow>
          </div>

          <div key="codepad">
            <STRMLWindow onMouseDown={this.gatherFocus} menuBar={[{ name: 'Ace Editor' }, { name: 'Build', items: ['Make, (ctrl-enter)', 'Test, (ctrl-\\)', 'Reset, (ctrl-r)']}, { name: 'Theme', items: Object.keys(themeMap) }]} onSelect={this.handleSelect}>
              <AceEditor
                  ref={this.editor}
                  style={fill}
                  mode={this.state.aceMode}
                  theme={this.state.aceTheme}
                  onFocus={() => this.toggleFocus('editor')}
                  onChange={this.handleAceChange}
                  name="UNIQUE_ID_OF_DIV"
                  editorProps={{$blockScrolling: true}}
                  value={this.state.aceContent}
                />
            </STRMLWindow>
          </div>

          <div key="challenge_details">
            <STRMLWindow onMouseDown={this.gatherFocus}  menuBar={[{ name: 'Challenge Details' }]} onSelect={this.handleSelect}>
              <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit. Mollitia officia, illo magni. Consequatur sapiente adipisci eos, fugit maxime velit necessitatibus corporis illo ut molestiae et temporibus ipsum quas voluptatum deleniti.</p>
            </STRMLWindow>
          </div>

          <div key="test_results">
            <STRMLWindow onMouseDown={this.gatherFocus} menuBar={[{ name: 'Test Results' }]} onSelect={this.handleSelect}>
              <TestResults results={this.state.testResults}/>
            </STRMLWindow>
          </div>

          <div key="compiler_output">
            <STRMLWindow onMouseDown={this.gatherFocus} menuBar={[{ name: 'Compiler Output' }]} onSelect={this.handleSelect}>
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
              <textarea style={fill} ref={this.stdin} onMouseDown={() => {this.toggleFocus("stdin");this.gatherFocus()}} onChange={this.handleStdinChange} value={this.state.stdin}></textarea>
            </STRMLWindow>
          </div>


        </GridLayout>
      </STRMLWindow>
      </div>
    )
  }
}
export { STRMLGrid }
