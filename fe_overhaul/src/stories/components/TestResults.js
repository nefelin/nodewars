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
				<th>Output</th>
				<th>Pass/Fail</th>
				<th>Hint</th>
			</tr>
			{results.stdouts.map((k,i) => (
					<tr key={i}>
						<td>{results.stdouts ? results.stdouts[i] : 'None'}</td>
						<td><strong>{results.grades ? results.grades[i] : 'None'}</strong></td>
						<td>{results.hints? results.hints[i] : 'None'}</td>	
					</tr>
				)
			)}
		</tbody>	
	</table>
	</div>
)

export default TestResults