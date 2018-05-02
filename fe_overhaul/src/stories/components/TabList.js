import React, { Component } from 'react';

export class TabView extends Component {
	constructor(props){
		super(props)
		this.state = {
			selected: 0,
			currentContent: "",
			buffers: [
					new Buffer('Main'),
					new Buffer('Test')
			],
			
		}
	}
	
	componentWillMount(){
		if (!this.selected)
			this.handleSelect(0)
	}

	handleSelect = (selected) => {
		console.log('handleSelect state', this.state)
		if (this.state.selected < this.state.buffers.length){ // if selected isn't pointing to death

			// deep clone buffers
			const buffers = cloneBufs(this.state.buffers)
			
			console.log('buffer cloned', buffers)

			// write this.state.currentContent to buffer clone
			buffers[this.state.selected].content = this.state.currentContent
			console.log('currentContent written to buffer')
			// setState with new buffers
			this.setState({ buffers })
		}

		this.setState({
			selected,
			currentContent:this.state.buffers[selected].content
		})
		console.log('handleSelect is done')
	}

	handleNewTab = () => {
		// deep clone buffers
		const buffers = cloneBufs(this.state.buffers)

		// add new buffer to buffers
		buffers.push(new Buffer('New'))

		// set state with new buffers
		this.setState({buffers}, () => this.handleSelect(buffers.length-1))
		
	}

	handleRemoveTab = (i) => {
		console.log('handleRemove:', i)
		// switch to main buffer
		const switchTo = i > 0 ? i-1 : i+1
		this.handleSelect(switchTo)

		// deep clone buffers
		const buffers = cloneBufs(this.state.buffers)

		// remove tab
		buffers.splice(i, 1)

		// set state with new buffers
		this.setState({ buffers })
	}

	handleContentChange = (e) => {
		// should also track cursor position
		this.setState({currentContent:e.target.value})
	}

	render(){
		const tabs = this.state.buffers.map((o,i) => {
			return <span><a onClick={() => this.handleSelect(i)}>{o.label}</a> {this.state.buffers.length > 1 ? <button onClick={() => this.handleRemoveTab(i)}>-</button> : null} </span>
		})

		return (
			<div>
				<div className="tab-navbar"> {tabs} <button onClick={this.handleNewTab}>+</button></div>
				<input value={this.state.currentContent} onChange={this.handleContentChange}/>
			</div>
			)

	}
}

function Buffer(label="", content="", cursor=0, language="python") {
	this.label = label
	this.content = content
	this.cursor = cursor
	this.language = language
}

function cloneBufs(orig) {
	console.log('cloning:', orig)
	const newBufs = []
	
	for (let i = 0; i < orig.length; i++) {
		newBufs.push(new Buffer(orig[i].label, orig[i].content,orig[i].cursor,orig[i].language))
	}

	return newBufs
}

