import React, {Component} from 'react';
import TagSection from './tags/TagSection.jsx';
import StatisticSection from './statistics/StatisticSection.jsx';
import Socket from '../socket.js';

class App extends Component{
	constructor(props){
		super(props);
		this.state = {
			activeStatistic: "Tweets",
			data: {},
			tags: [],
			statisticsNames: ['Tweets', 'Links']
		};
	}

	componentDidMount() {
		let socket = this.socket = new Socket();
		socket.on('connect', this.onConnect.bind(this));
		socket.on('disconnect', this.onDisconnect.bind(this));
		socket.on('tag update', this.onTagUpdate.bind(this));
		socket.on('tag list', this.onTagList.bind(this));
	}

	onConnect() {
		this.setState({connected: true});
		this.socket.emit('tag list', '');
	}

	onDisconnect() {
		this.setState({connected: false});
	}

	onTagList(data) {
		let tags = data;
		tags.map(tag=> this.addTag(tag));
	}

	onTagUpdate(data) {
		let tags = data;
		console.log(data);
	}

	addTag(name){
		let {tags} = this.state;
		let tag = {id: tags.length, name};
		tags.push(tag);
		this.setState({tags});
		this.setState({activeTag: tag});
	}

	setTag(activeTag){
		let currectActiveTab = this.state.activeTag;
		/* socket.emit('tag unsubscribe', {tagName: currentActiveTag});*/
		this.setState({activeTag});
		this.socket.emit('tag update', {tagName: activeTag.name});
	}

	tabSelect(tabName) {
		this.setState({activeStatistic: tabName});
	}

	render(){
		return (
			<div className='app'>
				<div className='nav'>
					<TagSection
						{...this.state}
						addTag={this.addTag.bind(this)}
						setTag={this.setTag.bind(this)}
					/>
				</div>
					<StatisticSection
						{...this.state}
						tabSelect={this.tabSelect.bind(this)}
					/>
			</div>
		)
	}
}

export default App
