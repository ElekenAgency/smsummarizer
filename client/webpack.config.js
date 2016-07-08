var path = require('path');

module.exports = {
	entry: path.join(__dirname, 'index.js'),
	output: {
		filename: 'bundle.js'
	},
	module: {
		loaders: [
			{
				test: /\.jsx?$/,
				loader: 'babel-loader',
				exclude: /node_modules/
			}
		]
	}
}
