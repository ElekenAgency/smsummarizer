import React, {Component} from 'react';

class TagForm extends Component{
	onSubmit(e){
		e.preventDefault();
		const node = this.refs.tag;
		const tagName = node.value;
		this.props.addTag(tagName);
		node.value = '';
	}
	render(){
		return (
			<form onSubmit={this.onSubmit.bind(this)}>
				<div className='form-group'>
					<input
						className='form-control'
						placeholder='Add Tag'
						type='text'
						ref='tag'
					/>
				</div>
			</form>
		)
	}
}

TagForm.propTypes = {
	addTag: React.PropTypes.func.isRequired
}


export default TagForm
