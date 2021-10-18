const path = require("path");
const fileList = ["./src/validate_signup.ts", "./src/auto_complete.ts"];

module.exports = fileList.map((fileName) => ({
  entry: fileName,
  module: {
    rules: [
      {
        test: /\.tsx?$/,
        use: "ts-loader",
        exclude: /node_modules/,
      },
    ],
  },
  resolve: {
    extensions: [".ts", ".js"],
  },
  output: {
    path: path.join(__dirname, path.join("../", "static/js")),
    filename: fileName
      .split("/")
      [fileName.split("/").length - 1].replace(".ts", ".js"),
  },
  mode: "development",
}));
