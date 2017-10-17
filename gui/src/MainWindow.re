open Utils;

requireCSS "./MainWindow.css";

let logo = requireAssetURI "./images/logo.png";

let component = ReasonReact.statelessComponent "MainPage";

let make _children => {
  ...component,
  render: fun _self =>
    <div className="mainWindow">
      <h1 className="title"> (textEl "Coming soon !") </h1>
      <img className="logo" src=logo />
    </div>
};
