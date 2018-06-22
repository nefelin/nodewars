import React, { Component } from 'react';

const GameDetails = (props) => {
    console.log('Game Details', props.game)
    const dets = Object.keys(props.game).map( (f, i) => (
        <li>
        <strong> {f}: </strong> {props.game[f]}
        </li>
    ))

    return (
        <div>
            <ul>
                {dets}
            </ul>
            <button onClick={props.onCancel}>Back</button>
        </div>
    )
}

export default GameDetails