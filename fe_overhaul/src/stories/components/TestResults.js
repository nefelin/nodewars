import React, { Component } from 'react'
import './TestResults.css'
const style = {
	height: '100%',
	overflowY: 'scroll',
	overflowX: 'hidden',
	boxSizing: "border-box"
}
const TestResults = ({results}) => {
	console.log(results)

	let tableHead, tableBody

	if (results.stdouts != undefined) {
		tableHead = <tr><th>Output</th></tr>
		tableBody = <tr><td>{results.stdouts != "" ? results.stdouts[0] : 'No Output'}</td></tr>
	} else {
		tableHead = <tr>
						<th>#</th>
						<th>Pass/Fail</th>
						<th>Hint</th>
					</tr>
		tableBody = results.grades.map((k,i) => (
							<tr key={i}>
								<td>{i}</td>
								<td><strong>{results.grades ? results.grades[i] : 'None'}</strong></td>
								<td>{results.hints? results.hints[i] : 'None'}</td>	
							</tr>
						))
	}


	return (

		<div style={style}>
		<table>
			<tbody>
				{tableHead}
				{tableBody}
			</tbody>	
		</table>
		</div>
	)
}

export default TestResults