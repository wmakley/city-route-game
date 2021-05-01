const path = require('path')
const webpack = require('webpack')
const { VueLoaderPlugin } = require('vue-loader');
const WebpackDevServer = require('webpack-dev-server');

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

	return {
		entry: {
			'admin': path.join(__dirname, 'javascript/admin/admin.js')
		},

		output: {
			path: path.join(__dirname, 'static'),
			filename: '[name].bundle.js',
		},

		module: {
			rules: [
				{
					test: /\.vue$/,
					loader: 'vue-loader'
				},
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
