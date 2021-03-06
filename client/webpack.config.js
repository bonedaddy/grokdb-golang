var
path = require('path'),
webpack = require('webpack');

// postcss
// var autoprefixer = require('autoprefixer-core');
// var cssimport = require('postcss-import');


// chokidar = require('chokidar'),
// anybar = require('anybar');

// chokidar.watch('./src/**/*', {ignored: /[\/\\]\./}).on('all', function(event, path) {
//     anybar('red');
// });

// chokidar.watch('./static/**/*', {ignored: /[\/\\]\./}).on('all', function(event, path) {
//     anybar('green');
// });

var appRoot = __dirname;

module.exports = {
    colors: true,
    watch: true,
    entry: {
        app: "./entry.js",
        vendor: [
            'babel/polyfill',
            'babel-runtime/regenerator',
            'babel-runtime/core-js',
            'react',
            'react-dom',
            'dataloader',
            'immutable',
            'page',
            'bluebird',
            'minitrue',
            'orwell',
            'co',
            'lodash',
            'react-either',
            'superagent/lib/client',
            'classnames',
            'shallowequal',
            'slug',
            'markdown-it',
            'react-textarea-autosize',
            'qs',
            'moment',
            'react-prop-once',
            'scriptjs',
            'markdown-it-abbr',
            'markdown-it-footnote',
            'markdown-it-sub',
            'markdown-it-sup'
            // 'intl',
            // 'between2',
            // 'pouchdb',
            // 'pouchdb/extras/memory'
        ]
    },
    output: {
        path: path.join(appRoot, "../assets"),
        filename: "app.js"
    },
    devtool: "#source-map",
    resolve: {
        root: path.join(appRoot, "/"),
        modulesDirectories: ["node_modules"]
    },
    module: {
        loaders: [
            // needed for markdown-it
            {
                test: /\.json$/,
                loader: "json-loader"
            },
            // {
            //     test:   /\.css$/,
            //     loader: "style-loader!css-loader!postcss-loader"
            // },
            {
                test: /\.jsx?$/,
                exclude: /node_modules/,
                loader: 'babel-loader',
                query: {
                    optional: ['runtime', 'validation.react']
                }
            },
            {
                test: /\.jsx?$/,
                loader: "eslint-loader",
                exclude: /node_modules/
            }

        ]
    },
    plugins: [
        new webpack.DefinePlugin({
            __DEV__: JSON.stringify(JSON.parse(process.env.BUILD_DEV || 'true')),
            __PRERELEASE__: JSON.stringify(JSON.parse(process.env.BUILD_PRERELEASE || 'false'))
        }),
        new webpack.optimize.CommonsChunkPlugin(/* chunkName= */"vendor", /* filename= */"vendor.js")
    ],
    // postcss: function () {
    //     return {
    //         defaults: [
    //             autoprefixer,
    //             cssimport({
    //                 // see postcss-import docs to learn about onImport callback
    //                 // https://github.com/postcss/postcss-import

    //                 onImport: function (files) {
    //                     files.forEach(this.addDependency);
    //                 }.bind(this)
    //             })
    //         ],
    //         cleaner:  [autoprefixer({ browsers: [] })]
    //     };
    // }
};
