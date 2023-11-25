const CopyPlugin = require('copy-webpack-plugin');
const path = require('path');
const WebpackObfuscator = require('webpack-obfuscator');

module.exports = {
  entry: {
    main: './dist/main.js'
  },
  output: {
    path: path.resolve(__dirname, 'dist'),
    filename: '[name].js'
  },
  plugins: [
    new CopyPlugin({
      patterns: [{from: 'static'}]
    }),
    new WebpackObfuscator({
      optionsPreset: 'high-obfuscation'
    }, [])
  ]
}
