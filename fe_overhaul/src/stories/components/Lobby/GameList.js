import React, { Component } from 'react';
var games = [
    {
        name: "Main",
        players: '2',
        nodes: '12',
        map: "DemoMap",
        languages: ['Python', 'Golang', 'JavaScript'],
        languageLock: 'true',
        autoBalance: 'true',
        private: 'true'
    },
    {
        name: "Other",
        players: '2',
        nodes: '12',
        map: "DemoMap",
        languages: ['Python', 'Golang', 'JavaScript'],
        languageLock: 'true',
        autoBalance: 'true',
        private: 'true'
    },
    {
        name: "Funhouse",
        players: '2',
        nodes: '12',
        map: "DemoMap",
        languages: ['Python', 'Golang', 'JavaScript'],
        languageLock: 'true',
        autoBalance: 'true',
        private: 'true'
    }
]
function GameList() {
    const tableHead = <tr>
        <th>Name</th>
        <th>Players</th>
        <th>Size</th>
        <th>Map</th>
        {/* <th>Private</th> */}
    </tr>
    const tableBody = games.map((g, i) => (
        <tr key={i}>
            <td>{g.name}</td>
            <td>{g.players}</td>
            <td>{g.nodes}</td>
            <td>{g.map}</td>
            {/* <td>{g.private}</td> */}
            <td style={{backgroundColor: "white", border: "none"}}><a href="http://donothing">Join</a></td>
            <td style={{ backgroundColor: "white", border: "none" }}><a href="http://donothing">Details</a></td>
        </tr>
    ))

    return (
    <div>
        <table>
            <tbody>
                {tableHead}
                {tableBody}
            </tbody>
        </table>
        <button>Create New Game</button>
    </div>
    )    
}

export default GameList

