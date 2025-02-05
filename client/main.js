const { app, BrowserWindow } = require("electron/main");
const path = require("node:path");

function createWindow() {
  const win = new BrowserWindow({
    width: 1000,
    height: 800,
    webPreferences: {
      nodeIntegration: true,
    },
  });

  // Determine the operating mode: development or production
  const startURL =
    process.env.NODE_ENV === "development"
      ? "http://localhost:8080"
      : `file://${path.join(__dirname, "dist", "index.html")}`;

  win.loadURL("http://localhost:8080");
}

app.whenReady().then(() => {
  createWindow();

  app.on("activate", () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      createWindow();
    }
  });
});

app.on("window-all-closed", () => {
  if (process.platform !== "darwin") {
    app.quit();
  }
});
