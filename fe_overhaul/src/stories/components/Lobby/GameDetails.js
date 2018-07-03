import React, { Component } from 'react';
import './GameDetails.css'
const GameDetails = (props) => {
    console.log('Game Details', props.game)
    const dets = Object.keys(props.game).map( (f, i) => (
        <li key={f}>
        <strong> {f}: </strong> {props.game[f]}
        </li>
    ))

    return (
        <div id="flex-container">
            <div >
                <img style={{width:"32px", height:"32px", display: "inline"}} src="./icons/back.svg" onClick={props.onCancel}/>
                
            </div>
            <div style={{position: "relative"}} >
                <ul style={{listStyleType:"none"}}>
                    {dets}
                </ul>
                <a style={{cursor:"pointer", position: "absolute", bottom: "3px", right: "5px"}} onClick={props.onJoin}>Join</a>
            </div>
            {/* <div>
                
            </div> */}
        </div>
    )
}

export default GameDetails