import React, { Component } from 'react'
import './STRMLMenu.css'

export const STRMLMenuBar = ({menus, onSelect}) => {
	const bar = []
	for (let menu of menus) {
		bar.push(<STRMLMenu key={menu.name} name={menu.name} items={menu.items ? menu.items : []} onSelect={onSelect}/>)
	}

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

	}

	// componentDidMount() {
		// console.log('head', window.getComputedStyle(this.menuTitle.current).left)
		// console.log('body', this.menuBody.current.style.left)
	// }

	showMenu = (event) => {	    
		this.setState({ show: true }, () => {
			document.addEventListener('click', this.closeMenu);
		});
	}

	closeMenu = () => {
		this.setState({ show: false }, () => {
			document.removeEventListener('click', this.closeMenu);
		});
	}

	handleSelect = (item) => {
		console.log(item)
	}

	render() {
		return (
			<div className="strml-menu">
				<div ref={this.menuTitle} className={"strml-menu-title" + (this.state.items.length > 0 ? " has-contents" : "") + (this.state.show ? " active" : "")}onClick={this.state.items.length > 0 ? this.showMenu : null}>
					{this.props.name}
				</div>
					<div ref={this.menuBody} className={"strml-menu-body" + (this.state.show ? "" : " hidden")}>
						{this.state.items.map((item) => (
							<div key={item} className="strml-menu-item" onClick={() => this.props.onSelect(this.props.name, item)}>
								{item}
							</div>
						))}

					</div>

			</div>
		)
	}
}

export default STRMLMenu