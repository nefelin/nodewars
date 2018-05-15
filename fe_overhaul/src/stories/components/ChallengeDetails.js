import React, {Component} from 'react'

const ChallengeDetails = ( { challenge } ) => {
	const style = {
		height: '100%',
		padding: 10,
		overflowY: 'scroll',
		overflowX: 'scroll',
		boxSizing: "border-box",
	}

	if (!challenge)
		return <div style={style}>No Challenge</div>

	console.log('ChallengeDetails', challenge)
	return (<div style={style}>
				{challenge.name}<br/><br/>
				
				{challenge.shortDesc}<br/><br/>
				
				{challenge.longDesc != "" ? "Details: " + challenge.longDesc : null}

				Example(s):<br/>
				{challenge.sampleIO.map((sample, i) => (
					<div>
					<span key={i}>{sample.input + " -> " + sample.expect}</span><br/>
					</div>
				))}

			</div>
	)
}

export default ChallengeDetails