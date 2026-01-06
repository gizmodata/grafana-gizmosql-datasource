const path = require('path');
const CopyWebpackPlugin = require('copy-webpack-plugin');
const ForkTsCheckerWebpackPlugin = require('fork-ts-checker-webpack-plugin');

module.exports = (env, argv) => {
  const isProduction = argv.mode === 'production';

  return {
    mode: isProduction ? 'production' : 'development',
    devtool: isProduction ? 'source-map' : 'eval-source-map',
    entry: './src/module.ts',
    output: {
      path: path.resolve(__dirname, 'dist'),
      filename: 'module.js',
      library: {
        type: 'amd',
      },
      clean: false, // Don't clean - we want to keep the Go binary
      devtoolModuleFilenameTemplate: '[resource-path]',
    },
    externals: [
      'lodash',
      'react',
      'react-dom',
      '@grafana/data',
      '@grafana/runtime',
      '@grafana/ui',
      '@emotion/css',
      function ({ request }, callback) {
        // Externalize all @grafana/* packages
        if (/^@grafana\//.test(request)) {
          return callback(null, request);
        }
        callback();
      },
    ],
    resolve: {
      extensions: ['.ts', '.tsx', '.js', '.jsx', '.json'],
    },
    module: {
      rules: [
        {
          test: /\.[tj]sx?$/,
          exclude: /node_modules/,
          use: {
            loader: 'swc-loader',
            options: {
              jsc: {
                parser: {
                  syntax: 'typescript',
                  tsx: true,
                  decorators: false,
                  dynamicImport: true,
                },
                target: 'es2015',
                loose: false,
              },
              sourceMaps: true,
            },
          },
        },
        {
          test: /\.css$/,
          use: ['style-loader', 'css-loader'],
        },
        {
          test: /\.s[ac]ss$/,
          use: ['style-loader', 'css-loader', 'sass-loader'],
        },
        {
          test: /\.(png|jpe?g|gif|svg)$/,
          type: 'asset/resource',
          generator: {
            filename: 'img/[name][ext]',
          },
        },
      ],
    },
    plugins: [
      new CopyWebpackPlugin({
        patterns: [
          { from: 'src/plugin.json', to: 'plugin.json' },
          { from: 'src/img', to: 'img', noErrorOnMissing: true },
          { from: 'README.md', to: 'README.md', noErrorOnMissing: true },
          { from: 'LICENSE', to: 'LICENSE', toType: 'file', noErrorOnMissing: true },
          { from: 'CHANGELOG.md', to: 'CHANGELOG.md', noErrorOnMissing: true },
        ],
      }),
      new ForkTsCheckerWebpackPlugin({
        typescript: {
          diagnosticOptions: {
            semantic: true,
            syntactic: true,
          },
        },
      }),
    ],
  };
};
