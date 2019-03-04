const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin')

module.exports = {
  entry: {
    index: path.resolve(__dirname, 'src', 'index.jsx'),
    sample_template: path.resolve(__dirname, 'src', 'sample_template.jsx'),
    notification_console: path.resolve(__dirname, 'src', 'notification_console.jsx'),
    explore_notification: path.resolve(__dirname, 'src', 'explore_notification.jsx'),
    notification_with_quote: path.resolve(__dirname, 'src', 'notification_with_quote.jsx'),
    remind_meetup_notification: path.resolve(__dirname, 'src', 'remind_meetup_notification.jsx'),
    signup_notification: path.resolve(__dirname, 'src', 'signup_notification.jsx'),
    generic_notification: path.resolve(__dirname, 'src', 'generic_notification.jsx'),
  },
  output: {
    path: path.resolve(__dirname, 'dist', 'assets'),
    filename: '[name].js'
  },
  resolve: {
    extensions: ['.js', '.jsx']
  },
  module: {
    rules: [
      {
        test: /\.jsx?$/,
        loader: 'babel-loader',
        exclude: /node_modules/,
        options: {
          presets: ['@babel/preset-react', '@babel/preset-env']
        }
      },
      {
        test: /\.scss/,
        use: ['style-loader', 'css-loader', 'sass-loader']
      }
    ]
  },
  devServer: {
   contentBase: './dist',
   publicPath: '/',
   port: 9000
  },
  externals: {
    'config': {
      'serverUrl': 'http://localhost',
    }
  },
  plugins: [
    new HtmlWebpackPlugin({
      template: 'src/index.html',
      filename: "../index.html",
      title: "{{.WebpageTitle}}",
      chunks: ['index']
    }),
    new HtmlWebpackPlugin({
      template: 'src/index.html',
      filename: "../sample_template.html",
      title: "{{.WebpageTitle}}",
      chunks: ['sample_template']
    }),
    new HtmlWebpackPlugin({
      template: 'src/index.html',
      filename: "../notification_console.html",
      title: "Notification Console",
      chunks: ['notification_console']
    }),
    new HtmlWebpackPlugin({
      template: 'src/index.html',
      filename: "../explore_notification.html",
      title: "Explore",
      chunks: ['explore_notification']
    }),
    new HtmlWebpackPlugin({
      template: 'src/index.html',
      filename: "../notification_with_quote.html",
      title: "Notification",
      chunks: ['notification_with_quote']
    }),
    new HtmlWebpackPlugin({
      template: 'src/index.html',
      filename: "../remind_meetup_notification.html",
      title: "Reminder to Meet Up",
      chunks: ['remind_meetup_notification']
    }),
    new HtmlWebpackPlugin({
      template: 'src/index.html',
      filename: "../signup_notification.html",
      title: "Welcome to Hive!",
      chunks: ['signup_notification']
    }),
    new HtmlWebpackPlugin({
      template: 'src/index.html',
      filename: "../generic_notification.html",
      title: "Notification",
      chunks: ['generic_notification']
    })
  ]
};
