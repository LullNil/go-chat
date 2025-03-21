const { contextBridge, ipcRenderer } = require("electron");

contextBridge.exposeInMainWorld("electronAPI", {
  // Methods for custom title bar actions
  minimize: () => ipcRenderer.send("window-minimize"),
  maximize: () => ipcRenderer.send("window-maximize"),
  close: () => ipcRenderer.send("window-close"),

  // Methods for token management
  saveToken: (token) => ipcRenderer.invoke("save-token", token),
  getToken: () => ipcRenderer.invoke("get-token"),
  deleteToken: () => ipcRenderer.invoke("delete-token"),
});
