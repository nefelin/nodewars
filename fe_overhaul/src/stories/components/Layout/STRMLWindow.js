import React, { Component } from 'react'
import STRMLMenu, { STRMLMenuBar } from './STRMLMenu'
import './STRMLWindow.css'

const STRMLWindow = (props) => {
	// console.log('STRMLWindow Rendering', props)
	const bgColor = props.bgOverride || 'white'
	return (
		<div className="strml-window">
			<div className="strml-window-header noselect">
				<STRMLMenuBar menus={props.menuBar} onSelect={(menu, selection) => props.onSelect(props.menuBar[0].name, menu, selection)}/>
				
			</div>
			<div style={{backgroundColor:bgColor}} className="strml-window-content">
				{props.children}
			</div>
		</div>
	)
}

export default STRMLWindow