import React, {Component} from 'react';

import { storiesOf } from '@storybook/react';
import { action } from '@storybook/addon-actions';
import { linkTo } from '@storybook/addon-links';

import { Button, Welcome } from '@storybook/react/demo';

import { TabView } from './components/TabList'
import { NodePie } from './components/Graph/NodePie'
import Graph from './components/Graph/Graph'

import * as Maps from './maps'
import { App } from './components/Graph/ReactForceLayout'
import { POEPing } from './components/Graph/POEPing'
import { CycleTest } from './components/Graph/cycle_test'
import { ColorTransition } from './components/Graph/ColorTransition'
import { ClipMask } from './components/Graph/ClipMask'
import { SelectionSize } from './components/Graph/SelectionSize'
import { UIMenu } from './components/Graph/UIMenu'
import { PiePower } from './components/Graph/PiePower'
import { SplitGrid } from './components/Graph/SplitGrid'
import { STRMLGrid } from './components/Layout/STRMLGrid'
import  STRMLWindow from './components/Layout/STRMLWindow'
import { STRMLMenu } from './components/Layout/STRMLMenu'

import TinyTerm from './components/Terminal/TinyTerm'

storiesOf('Welcome', module).add('to Storybook', () => <Welcome showApp={linkTo('Button')} />);

storiesOf('Button', module)
  .add('with text', () => <Button onClick={action('clicked')}>Hello Button</Button>)
  .add('with some emoji', () => (
    <Button onClick={action('clicked')}>
      <span role="img" aria-label="so cool">
        😀 😎 👍 💯
      </span>
    </Button>
  ));

storiesOf('Tab Stuff', module)
	.add('TabView', () => (
		<TabView startLabel="main" onChange="(d) => alert(d)"/>
	))

storiesOf('ScoreBoard', module)
	.add('Pie Chart', () => (
		<p>This should be an a pie chart that animates on update, to be used for territory and/or production</p>
	))
	.add('Score Bars', () => (
		<p>These should show overall CoinCoin score clearly</p>
	))
	.add('Score Board', () => (
		<p>This should combine Pie and Bar charts with text readout to display the game sitch clearly</p>
	))

storiesOf('Graph Map', module)
	.add('D3 Wrapper', () => {
		return (
			<div style={{height:400, width:400, border:'1px solid black'}}>
			<p>Graph should tie all elements together into a force layout</p>
			<Graph debug={true} dataset={ Maps.SingleNode } onClick={() => alert('test')}/>

			</div>
		)
	})
	.add('React Centric Force Layout', () => {
		return (
			<div>
			<p>The trick with this one is gonna be how to use pie charts.</p>
			<p>If pie charts aren't too difficult its really a question of weather we need animation</p>
			<App />
			</div>
		)
	})
	.add('Poe Animation Idea', () => (
		<POEPing dataset={[]}/>
	))

storiesOf('Layout', module)
	.add('Splitter Layout', () => {
		return <SplitGrid />
	})
	.add('STRML Grid', () => {
		return <STRMLGrid />
	})
	.add('STRML Window', () => {
		return (
			<div style={{height:350, width:500}}>
				<STRMLWindow keykey="map" menuBar={[{ name: 'Map' }, { name: 'Theme', items: ['Light', 'Dark'] }]} onSelect={this.handleSelect}>
		          <Graph dataset={ Maps.SimpleMap }/>
		        </STRMLWindow>
			</div>
		)

	})
	.add('BorderBox Lab', () => {
		const bigbox = {
			position: 'relative',
			height: 200,
			width: 300,
			border: '1px solid black',
			overflowY: 'scroll',
		}

		const title = {
			zIndex: 10,
			position:'absolute',

			paddingLeft: 5,//'2px 0 2px 0 5px',
			paddingTop: 2,
			paddingBottom: 2,

			width: '100%',
			borderBottom: '1px solid black',
			boxSizing: 'border-box',
			backgroundColor: 'white',

			fontSize: '13px',
		}

		const content = {
			position: 'absolute',
			top: 20,
			bottom: 0,
			left: 0,
			right: 0,
			// marginTop: 19,
			// height: '100%',
			// width: '100%',
			backgroundColor: 'blue',
			boxSizing: 'border-box'
		}


		return (
			<div style={bigbox}>
				<div style={title}> Title </div>
					<div style={content}>
					</div>
				
			</div>
		)
	})

storiesOf('Terminal', module)
	.add('TinyTerm', () => {
		class Temp extends React.Component {
			constructor(props) {
				super(props)
				this.state = {
					message:'hey'
				}
			}

			render() {
				return (
					<div style={{height:150}}>
						<TinyTerm incoming={this.state.message} />
						<button onClick={() => this.setState({message:'new message\n'})}>update</button>
					</div>
				)
			}
		}

		return <Temp />
	})

storiesOf('Sandbox', module)
	.add('Text Formatting', () => {
		const s1 = "a      - whitespace heavy string"
		return (
			<p> {s1} </p>
		)
	})
	.add('Cycle Test', () => {
		return <CycleTest />
	})
	.add('Color Transition', () => {
		return <ColorTransition />
	})
	.add('Clipping Masks', () => {
		return <ClipMask />
	})
	.add('Selection Size', () => {
		return <SelectionSize />
	})
	.add('UIMenu Test', () => {
		return <UIMenu />
	})
	.add('PiePower', () => {
		return <PiePower />
	})

