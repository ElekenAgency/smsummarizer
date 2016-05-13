import React, {Component} from 'react';

class Tag extends Component{
	onClick(e){
		e.preventDefault();
		const {setTag, tag} = this.props;
		setTag(tag);
	}
	render(){
		const {tag, activeTag} = this.props;
		const active = tag === activeTag ? 'active' : '';
		return (
			<li className={active}>
				<a onClick={this.onClick.bind(this)}>
					{tag.name}
				</a>
			</li>
		)
	}
}

Tag.propTypes = {
	tag: React.PropTypes.object.isRequired,
	setTag: React.PropTypes.func.isRequired,
	activeTag: React.PropTypes.object.isRequired
}

export default Tag
