const path = require('path')
const webpack = require('webpack')
const { VueLoaderPlugin } = require('vue-loader')
const WebpackDevServer = require('webpack-dev-server')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')

module.exports = function (env, argv) {
	process.env.NODE_ENV = argv.mode;

	const isProduction = argv.mode === 'production';

	const plugins = [
		new VueLoaderPlugin(),
		// Vue feature flags (shuts up warnings), see: https://github.com/vuejs/vue-next/tree/master/packages/vue#bundler-build-feature-flags
		// Using defaults for now.
		new webpack.DefinePlugin({
			__VUE_OPTIONS_API__: JSON.stringify(true),
			__VUE_PROD_DEVTOOLS__: JSON.stringify(false)
		})
	];

	if (isProduction) {
		plugins.push(new MiniCssExtractPlugin());
	}

	return {
		entry: {
			'admin': path.join(__dirname, 'javascript/admin/admin.js')
		},

		output: {
			path: path.join(__dirname, 'static/admin'),
			filename: '[name].bundle.js',
		},

		module: {
			rules: [
				{
					test: /\.js$/,
					enforce: 'pre',
					use: ['source-map-loader'],
				},
				{
					test: /\.vue$/,
					loader: 'vue-loader'
				},
				{
					test: /\.css$/,
					use: [
						!isProduction
							? 'vue-style-loader'
							: MiniCssExtractPlugin.loader,
						'css-loader'
					]
				}
			],
		},

		plugins: plugins,

		devServer: {
			contentBase: path.join(__dirname, 'static'),
			compress: true,
			port: 9000,
		}
	};
}
