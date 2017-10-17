// This is the entry point of the electron app
// It manages the overall state of the native app (menus, windows, ...)
// TODO: launch cameradar process

const { app, BrowserWindow } = require('electron');
const path = require('path');
const url = require('url');

let mainWindow;

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.on('ready', () => {
  // create the browser window.
  mainWindow = new BrowserWindow({ width: 800, height: 600 });

  // and load the index.html of the app.
  mainWindow.loadURL(
    url.format({
      pathname: `${__dirname}/../build/index.html`,
      protocol: 'file:',
    }),
  );

  // Open the DevTools.
  mainWindow.webContents.openDevTools({
    mode: 'bottom',
  });

  mainWindow.on('closed', () => {
    mainWindow = null;
  });
});
