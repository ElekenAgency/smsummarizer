import React, {Component} from 'react';

class Statistic extends Component{
	render(){
		let {stype} = this.props;
		let {data} = this.props;
		let actualData = data[stype];
		return (
			<table className="table table-striped">
			<thead>
			<tr>
			<th>{stype}</th>
			<th>Likes</th>
			</tr>
			</thead>
			<tbody>{
				actualData.map(stat => {
					return (
						<tr>
							<td>{stat.data}</td>
							<td>{stat.count}</td>
						</tr>
					)
				})
			}
			</tbody>
			</table>
		);
	}
}

Statistic.propTypes = {
	stype: React.PropTypes.object.isRequired,
	data: React.PropTypes.array.isRequired,
}

export default Statistic
