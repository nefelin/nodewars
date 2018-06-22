import React, { Component } from 'react';
var game = {
        name: "Main",
        players: '2',
        nodes: '12',
        map: "DemoMap",
        languages: ['Python', 'Golang', 'JavaScript'],
        languageLock: 'true',
        autoBalance: 'true',
        private: 'true'
    }

const GameDetails = () => {
    const dets = Object.keys(game).map( (f, i) => (
        <li>
        <strong> {f}: </strong> {game[f]}
        </li>
    ))

    return (
        <ul>
            {dets}
        </ul>
    )
}

export default GameDetails