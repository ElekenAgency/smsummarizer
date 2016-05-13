import React, {Component} from 'react';
import Tag from './Tag.jsx';

class TagList extends Component{
	render(){
		return (
			<ul>{
				this.props.tags.map( chan =>{
					return <Tag
						tag={chan}
						key={chan.id}
						{...this.props}
					/>
				})
			}</ul>
		)
	}
}

TagList.propTypes = {
	tags: React.PropTypes.array.isRequired,
	setTag: React.PropTypes.func.isRequired,
	activeTag: React.PropTypes.object.isRequired
}

export default TagList
