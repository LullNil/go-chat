const { app, ipcMain, BrowserWindow, globalShortcut } = require("electron");
const path = require("node:path");
const fs = require("fs");

let mainWindow;

// Path to the file where the window state (position and size) is saved.
// This file will be stored in the directory returned by app.getPath('userData').
const stateFile = path.join(app.getPath("userData"), "window-state.json");

// Default window state
let windowState = {
  width: 950,
  height: 800,
};

// If the state file exists, read the saved window state
if (fs.existsSync(stateFile)) {
  try {
    const savedState = JSON.parse(fs.readFileSync(stateFile, "utf8"));
    windowState = { ...windowState, ...savedState };
  } catch (err) {
    console.error("Error reading window state file:", err);
  }
}

// Variable to track whether DevTools are docked (true) or detached (false)
// Default is detached mode
let devtoolsDocked = false;

function createWindow() {
  mainWindow = new BrowserWindow({
    width: windowState.width,
    height: windowState.height,
    x: windowState.x, // Use saved x position if available
    y: windowState.y, // Use saved y position if available
    frame: false,
    titleBarStyle: "hidden",
    webPreferences: {
      devTools: process.env.NODE_ENV === "development" ? true : false,
      nodeIntegration: false,
      contextIsolation: true,
      preload: path.join(__dirname, "preload.js"),
    },
  });

  // Determine the start URL based on the environment
  const startURL =
    process.env.NODE_ENV === "development"
      ? "http://localhost:8080"
      : `file://${path.join(__dirname, "dist", "index.html")}`;

  mainWindow.loadURL(startURL);

  // In development mode, open DevTools in detached mode by default
  if (process.env.NODE_ENV === "development") {
    mainWindow.webContents.openDevTools({ mode: "detach" });
    devtoolsDocked = false;
  }

  // Save window state (size and position) when the window is resized or moved
  mainWindow.on("resize", saveWindowState);
  mainWindow.on("move", saveWindowState);
}

function saveWindowState() {
  if (!mainWindow) return;
  const bounds = mainWindow.getBounds();
  try {
    fs.writeFileSync(stateFile, JSON.stringify(bounds));
  } catch (err) {
    console.error("Error saving window state:", err);
  }
}

// IPC events for custom title bar actions
ipcMain.on("window-minimize", (event) => {
  const win = BrowserWindow.fromWebContents(event.sender);
  if (win) win.minimize();
});

ipcMain.on("window-maximize", (event) => {
  const win = BrowserWindow.fromWebContents(event.sender);
  if (win) {
    if (win.isMaximized()) {
      win.unmaximize();
    } else {
      win.maximize();
    }
  }
});

ipcMain.on("window-close", (event) => {
  const win = BrowserWindow.fromWebContents(event.sender);
  if (win) win.close();
});

app.whenReady().then(() => {
  createWindow();

  // Register a global shortcut to toggle DevTools docking mode.
  // Pressing Ctrl+Shift+D will switch between detached and docked modes.
  if (process.env.NODE_ENV === "development") {
    globalShortcut.register("CommandOrControl+Shift+D", () => {
      if (mainWindow && mainWindow.webContents.isDevToolsOpened()) {
        // Close current DevTools
        mainWindow.webContents.closeDevTools();
        // Toggle the DevTools mode
        if (devtoolsDocked) {
          // Currently docked, so open in detached mode
          mainWindow.webContents.openDevTools({ mode: "detach" });
          devtoolsDocked = false;
        } else {
          // Currently detached, so open in docked mode (e.g., docked to the right)
          mainWindow.webContents.openDevTools({ mode: "right" });
          devtoolsDocked = true;
        }
      }
    });
  }

  app.on("activate", () => {
    // On macOS, it's common to re-create a window in the app when the dock icon is clicked
    if (BrowserWindow.getAllWindows().length === 0) {
      createWindow();
    }
  });
});

// Unregister all shortcuts when the app is about to quit
app.on("will-quit", () => {
  globalShortcut.unregisterAll();
});

// Quit the app when all windows are closed (except on macOS)
app.on("window-all-closed", () => {
  if (process.platform !== "darwin") {
    app.quit();
  }
});
