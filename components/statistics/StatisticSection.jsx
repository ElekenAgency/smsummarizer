import React, {Component} from 'react';
import Statistic from './Statistic.jsx';
import StatisticTabs from './StatisticTabs.jsx';

class StatisticSection extends Component{
	render(){
		let {activeTag} = this.props;
		let section, name = "";
		if (activeTag !== undefined) {
			section = (
				<div className='panel-body statistics'>
					<StatisticTabs {...this.props}
						/>
					<Statistic {...this.props}
						stype={this.props.activeStatistic}/>
				</div>
			)
		}
		return (
			<div className='statistics-container panel panel-default'>
				<div className='panel-heading'>
					<strong>{name}</strong>
				</div>
				{section}
			</div>
		)
	}
}

StatisticSection.propTypes = {
	statisticsNames: React.PropTypes.array.isRequired,
	tabSelect: React.PropTypes.func.isRequired,
	data: React.PropTypes.object.isRequired,
	activeTag: React.PropTypes.object.isRequired,
	activeStatistic: React.PropTypes.object.isRequired,
}

export default StatisticSection
