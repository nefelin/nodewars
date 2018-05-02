import React, { Component } from 'react'
import './TestResults.css'
const style = {
	height: '100%',
	overflowY: 'scroll',
	overflowX: 'hidden',
	boxSizing: "border-box"
}
const TestResults = ({results}) => (
	<div style={style}>
	<table>
		<tbody>
			<tr>
				<th>#</th>
				<th>Pass/Fail</th>
				<th>Hint</th>
			</tr>
			{results.map((res, i) =>{
				return (
					<tr key={i}>
						<td>{i}</td>
						<td><strong>{res.grade}</strong></td>
						<td>{res.hint ? res.hint : 'None'}</td>	
					</tr>
				)
			})}
		</tbody>	
	</table>
	</div>
)

export default TestResults