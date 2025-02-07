const { app, BrowserWindow } = require("electron/main");
const path = require("node:path");

let mainWindow;

function createWindow() {
  mainWindow = new BrowserWindow({
    width: 1000,
    height: 800,
    webPreferences: {
      nodeIntegration: true,
    },
  });

  // mainWindow.setMenu(null);

  // Determine the operating mode - development or production
  const startURL =
    process.env.NODE_ENV === "development"
      ? "http://localhost:8080"
      : `file://${path.join(__dirname, "dist", "index.html")}`;

  mainWindow.loadURL(startURL);

  // Opening Developer Tools in Development Mode
  if (process.env.NODE_ENV === "development") {
    mainWindow.webContents.openDevTools();
  }
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
