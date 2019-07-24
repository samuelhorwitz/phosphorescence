const path = require('path');
const webpack = require('webpack');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const CopyWebpackPlugin = require('copy-webpack-plugin');
require('dotenv').config({path: path.resolve(process.cwd(), '../.env')});

module.exports = () => {
    return {
        mode: 'development',
        entry: './index.js',
        plugins: [
            new webpack.DefinePlugin({
              'process.env.PHOSPHOR_ORIGIN': JSON.stringify(process.env.PHOSPHOR_ORIGIN)
            }),
            new HtmlWebpackPlugin({
                template: './index.ejs',
                inject: 'head',
                filename: 'index.html'
            }),
            new CopyWebpackPlugin([{
                from: 'robots.txt',
                to: './'
            }])
        ],
        output: {
            path: path.resolve(__dirname, 'dist'),
            filename: 'bundle.js',
            publicPath: '/'
        }
    };
}