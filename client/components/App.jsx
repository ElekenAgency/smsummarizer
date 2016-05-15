import React, {Component} from 'react';
import TagSection from './tags/TagSection.jsx';
import StatisticSection from './statistics/StatisticSection.jsx';

class App extends Component{
	constructor(props){
		super(props);
		this.state = {
			activeStatistic: "Tweets",
			data: {"Tweets":[
				{data:"Some tweet1", count:25},
				{data:"Some tweet2", count:22},
			],
				"Links":[
					{data:"Link 1", count:12},
					{data:"Link 2", count:12},
				]},
			tags: [],
			statisticsNames: ['Tweets', 'Links']
		};
	}

	addTag(name){
		let {tags} = this.state;
		let tag = {id: tags.length, name};
		tags.push(tag);
		this.setState({tags});
		this.setState({activeTag: tag});
	}

	setTag(activeTag){
		this.setState({activeTag});
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
