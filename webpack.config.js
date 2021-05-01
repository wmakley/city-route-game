const path = require('path')

const elmSource = path.resolve(__dirname, 'javascript/board-editor')

module.exports = function (env, argv) {
	process.env.NODE_ENV = argv.mode;

	const isProduction = argv.mode === 'production';

	return {
		entry: {
		},

		output: {
			path: path.resolve(__dirname, 'static'),
			filename: '[name].bundle.js',
		},

		module: {
		},

		devServer: {
			contentBase: path.join(__dirname, 'static'),
			compress: true,
    		port: 9000,
		}
	};
}
