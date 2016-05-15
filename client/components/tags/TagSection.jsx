import React, {Component} from 'react';
import TagForm from './TagForm.jsx';
import TagList from './TagList.jsx';

class TagSection extends Component{
	render(){
		return (
			<div className='support panel panel-primary'>
				<div className='panel-heading'>
					<strong>Tags</strong>
				</div>
				<div className='panel-body tags'>
					<TagList {...this.props} />
					<TagForm {...this.props} />
				</div>
			</div>
		)
	}
}

TagSection.propTypes = {
	tags: React.PropTypes.array.isRequired,
	setTag: React.PropTypes.func.isRequired,
	addTag: React.PropTypes.func.isRequired,
	activeTag: React.PropTypes.object.isRequired
}

export default TagSection
