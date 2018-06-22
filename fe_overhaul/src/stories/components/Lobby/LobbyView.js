import React, { Component } from 'react'
import GameList from './GameList'
import GameDetails from './GameDetails'
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
]

class LobbyView extends Component {
    constructor(props) {
        super(props)
        this.state = {
            show: "list",
            selected: {},
            games: games,
        }
    }

    handleSelect = (selected) => {
        // console.log('SELECT')
        this.setState({
            selected,
            show: "details",
        })
    }

    handleCancel = () => {
        // console.log('CANCEL')
        this.setState({
            focus: {},
            show: "list"
        })
    }

    render() {
        let view
        switch (this.state.show) {
            case "list":
            view = <GameList games={this.state.games} onSelect={this.handleSelect} onJoin={this.props.onJoin}/>
            break
            case "details":
            view = <GameDetails game={this.state.selected} onCancel={this.handleCancel} onJoin={this.props.onJoin}/>
            break
            // case "new":
            // return <NewGame onNewGame={this.props.onNewGame} onCancel={this.handleCancel}/>
        }

        return (
            <div style = {{height: "200px", width: "350px", border: "1px solid black"}}>
                {view}
            </div>
        )
    }
}

export default LobbyView