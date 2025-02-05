const { defineConfig } = require("@vue/cli-service");

module.exports = defineConfig({
  transpileDependencies: true,
  publicPath: process.env.NODE_ENV === "production" ? "./" : "/",

  pluginOptions: {
    electronBuilder: {
      nodeIntegration: true,
    },
  },

  devServer: {
    port: 8080,
    hot: true,
  },
});
