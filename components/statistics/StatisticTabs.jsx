import React, {Component} from 'react';

class StatisticTabs extends Component{
	constructor(props) {
		super(props);
		this.onClick = this.onClick.bind(this);
	}

	onClick(name){
		this.props.tabSelect(name);
	}

	render(){
		const {activeStatistic} = this.props;
		return (
			<div>
				<ul className="nav nav-tabs">{
					this.props.statisticsNames.map(stat => {
						const active = stat	=== activeStatistic ? 'active' : '';
						return (
							<li className={active}>
								<a onClick={() => this.onClick(stat)} id={stat}>{stat}</a>
							</li>
						)
					})
																 }
				</ul>
			</div>
		)
	}
}

StatisticTabs.propTypes = {
	statisticsNames: React.PropTypes.array.isRequired,
	tabSelect: React.PropTypes.func.isRequired,
}


export default StatisticTabs
