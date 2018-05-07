import React, { Component } from 'react'
import './TinyTerm.css'

/*
props:
response - gets written to the page and triggers new prompt
*/

class TinyTerm extends React.Component {
	constructor(props){
		super(props)

		this.state = {
			commandHistory: [],
			page: "",
			command: "",
			cursor: this.props.cursor || "\u258c",
			focused: "false"
		}

		this.handleKeyPress = this.handleKeyPress.bind(this)
		this.handleChange = this.handleChange.bind(this)
		this.focus = this.focus.bind(this)
		this.handleLoseFocus = this.handleLoseFocus.bind(this)

		this.sendCommand = this.sendCommand.bind(this)
		this.recv = this.recv.bind(this)

		this.input = React.createRef()
		this.container = React.createRef()
	}

	componentDidMount() {
		if (this.props.grabFocus)
			this.focus()
		else {
			this.input.current.blur()
			this.handleLoseFocus()
		}

	}
	
	// componentWillReceiveProps(props) {
	// 	// console.log('tinyTerm receiving props',props)
	// 	if (props.incoming!=""){
	// 		console.log('TinyTerm receiving', props)
	// 		this.updatePage(props.incoming)
	// 	}
	// 	// this.setState({
	// 	// 	page: this.state.page + props.incoming
	// 	// })
	// }

	recv(newContent) {
		this.setState({
			page: this.state.page + newContent,
		}, () => this.container.current.scrollTop = this.container.current.scrollHeight)
		// TODO is there a way to hide the scrollbar here?
	}

	focus() {
		this.input.current.focus()
		this.setState({focused: true})
	}

	handleLoseFocus() {
		// console.log('lost focus')
		this.setState({focused: false})
	}

	handleKeyPress(e) {
		// console.log(e.keyCode)
		switch (e.keyCode){
			case 13:
				this.sendCommand(e.target.value)
				break
		}
	}
	
	sendCommand(cmd) {
		// console.log((this.state.page + cmd + '\n').split('\n'))
		this.recv(cmd + '\n')
		this.setState({ command: "" })
		this.props.onSend(cmd)
	}

	handleChange(e) {
		// console.log('tinyTerm input', e)
		// const command = e.target.value.split(this.state.prompt)[1] ? command : ""
		this.setState({ command: e.target.value })
	}

	render() {
		const bgColor = 'black',
			  textColor = 'white'
		const style = {
			width: '100%',
			height: '100%',
			// color: textColor,
			// backgroundColor: bgColor,
			position: 'relative',
			paddingBottom: '10px',
			overflowY: 'scroll',
			boxSizing: 'border-box',
			padding: 5,
		}

		const inputStyle = {
			width: 0,
			opacity: 0
		}

		const contentStyle = {
			width: '100%',
			boxSizing: 'border-box',
			padding: 5,
			overflowY: 'scroll',

		}

		let pageLines = this.state.page.split('\n')
		// pageLines = pageLines.slice(0, pageLines.length-1)

		const pageContent = pageLines.map((line, key) => {
			// console.log('key', key, 'length', pageLines.length)
			return <span key={key}>{line}{key != pageLines.length-1 ? <br/> : null}</span>			
		})

		// console.log(pageContent)

		return (
			<div className="TinyTerm" onClick={this.focus} ref={this.container} style={style}>
					{pageContent} {this.state.command}
					
					<span className={"TinyTerm-cursor" + (this.state.focused ? " focused" : "")}>{this.state.cursor}</span>
								
				<input style={inputStyle} ref={this.input} value={this.state.command} onKeyDown={this.handleKeyPress} onChange={this.handleChange} onBlur={this.handleLoseFocus}></input>
			</div>
			
		)
	}
}

export default TinyTerm