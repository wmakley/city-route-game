const path = require('path')

const elmSource = path.resolve(__dirname, 'javascript/board-editor')

module.exports = function (env, argv) {
	process.env.NODE_ENV = argv.mode;

	const isProduction = argv.mode === 'production';

	return {
		entry: {
			'board-editor': path.resolve(__dirname, 'javascript/board-editor/index.js')
		},

		output: {
			path: path.resolve(__dirname, 'static'),
			filename: '[name].bundle.js',
		},

		module: {
			rules: [{
				test: /\.elm$/,
				exclude: [/elm-stuff/, /node_modules/],
				use: {
					loader: 'elm-webpack-loader',
					options: {
						cwd: elmSource,
					}
				}
			}]
		},

		devServer: {
			contentBase: path.join(__dirname, 'static'),
			compress: true,
    		port: 9000,
		}
	};
}
