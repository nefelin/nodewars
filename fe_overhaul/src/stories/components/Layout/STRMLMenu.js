import React, { Component } from 'react'
import './STRMLMenu.css'

export const STRMLMenuBar = ({menus, onSelect}) => {
	const bar = []
	for (let menu of menus) {
		bar.push(<STRMLMenu key={menu.name} name={menu.name} items={menu.items ? menu.items : []} onSelect={onSelect}/>)
	}
	// if (menus[0].name == 'Ace Editor'){
	// 	console.log('EDITOR rendering', menus)
	// }

	return <div>{bar}</div>
}

class STRMLMenu extends Component {
	constructor(props) {
		super(props)
		this.state = {
			show: false,
			items: this.props.items || [],
		}

		this.menuTitle = React.createRef()
		this.menuBody = React.createRef()

		this.clickedOpen = false

	}

	showMenu = (e) => {	
		// console.log('show', e.type)
		if (this.state.show)
			this.closeMenu(e)

		if (e.type == 'click'){
			this.clickedOpen = true
			// return
		} else this.clickedOpen = false


		if (this.state.items.length > 0) {
			this.setState({ show: true }, () => {
				document.addEventListener('mousedown', this.closeMenu);
				document.addEventListener('mouseup', this.closeMenu);
				// document.addEventListener('click', this.closeMenu);
			});
		}
	}

	closeMenu = (e) => {
		// console.log('close', e.type)

		if (this.clickedOpen == false || (this.clickedOpen == true && e.type == 'mousedown')) {
			this.setState({ show: false }, () => {
				document.removeEventListener('mousedown', this.closeMenu);
				document.removeEventListener('mouseup', this.closeMenu);
				// document.removeEventListener('click', this.closeMenu);
			});
		}
	}


	render() {
		return (
			<div className="strml-menu">
				<div ref={this.menuTitle} className={"strml-menu-title" + (this.state.items.length > 0 ? " has-contents" : "") + (this.state.show ? " active" : "")} onClick={this.showMenu} onMouseDown={this.showMenu}>
					{this.props.name}
				</div>
					<div ref={this.menuBody} className={"strml-menu-body" + (this.state.show ? "" : " hidden")}>
						{this.state.items.map((item) => (
							<div key={item} className="strml-menu-item" onMouseDown={() => this.props.onSelect(this.props.name, item)} onMouseUp={() => this.props.onSelect(this.props.name, item)}>
								{item}
							</div>
						))}

					</div>

			</div>
		)
	}
}

export default STRMLMenu