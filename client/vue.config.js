const { defineConfig } = require("@vue/cli-service");

module.exports = defineConfig({
  transpileDependencies: true,
  publicPath: "./",

  pluginOptions: {
    electronBuilder: {
      nodeIntegration: true, // Разрешает использование Node.js API в Electron
    },
  },

  devServer: {
    port: 8080,
    hot: true,
  },
});
