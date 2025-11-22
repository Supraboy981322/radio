const githubURL = "https://github.com/Supraboy981322/radio";
const musicElm = document.getElementById("music");
const wall = document.getElementById("wall");
const settingsWall = document.getElementById("settingsWall");
const settingsElm = document.getElementById("settings");
var currentStationButton = null;

window.onload = init();

async function fetchJSONasArray(url) {
  try {
    const response = await fetch(url);

    if (!response.ok) {
      throw new Error(`http err: ${response.status}`);
    }

    const data = await response.json();
    return data;
  } catch (error) {
    console.error("err fetching data:  ", error);
    return [];
  }
}

function init() {
  (async () => {
    fetch("logoLong.svg")
      .then(response => response.text())
      .then(xmlText => {
        const pageTitleCont = document.getElementById("pageTitleCont");
        pageTitleCont.innerHTML = xmlText
      })
      .catch(error => console.error("error fetching svg as xml: ", error));
  })();
  (async () => {
    const libraryArray = await fetchJSONasArray("library.json");
        
    if (Array.isArray(libraryArray)) {
      console.log("success fetching library");
      loadLibrary(libraryArray);
    } else {
      console.error("err expected json data as an array for fetching notes, but data is not an array!");
    }
  })();
}

function loadLibrary(library) {
  currentStationButton = document.getElementById("playButton0");
  for (i = 0; i < library.length; i++) {
    let stationData = library[i];
    let station = document.createElement("div");
    station.setAttribute("class", "station");
    station.setAttribute("id", `#playButton${i}`);
    station.setAttribute("onclick", `play("${stationData[1]}", "#playButton${i}");`);
    
    let title = document.createElement("span");
    title.setAttribute("class", "title");
    title.innerText = stationData[0];

    let playButton = document.createElement("button");
    playButton.setAttribute("id", `#playButton${i}`);
    playButton.innerText = "play";

    station.appendChild(title);
    station.appendChild(playButton);
    wall.appendChild(station);
    if (i == 0) { currentStationButton = playButton }
  }
  navigator.mediaSession.setActionHandler('pause', function () {
    let musicSRC = musicElm.getAttribute("src");
    let buttonId = currentStationButton.getAttribute("id");
    console.log(buttonId);
    console.log(musicSRC);
    stop(
      musicSRC,
      buttonId
    );
  });

  navigator.mediaSession.setActionHandler('play', function () {
    console.log(currentStationButton);
    let buttonOnclickAttr = currentStationButton.parentElement.getAttribute("onclick");
    let buttonId = currentStationButton.getAttribute("id");
    let musicSRC = buttonOnclickAttr.replace("play(\"", "").replace(`", "${buttonId}");`, "");
    console.log(musicSRC);
    play(
      musicSRC,
      buttonId
    );
    console.log(currentStationButton);
  });
}

function resetButtons() {
  const buttons = document.querySelectorAll(".station button");
  for (i = 0; i < buttons.length; i++) {
    buttons[i].innerText = "play";
    oldPlayFunc = buttons[i].parentElement.getAttribute("onclick");
    newPlayFunc = oldPlayFunc.replace("stop(", "play(");
    buttons[i].parentElement.setAttribute("onclick", newPlayFunc);
    buttons[i].removeAttribute("class");
  }
}

function play(url, playButtonId) {
  resetButtons();
  console.log(playButtonId);
  let playButton = document.querySelector(`button[id="${playButtonId}"]`);
  let stationElm = document.querySelector(`div[id="${playButtonId}"]`);
  currentStationButton = playButton;
  endStream();
  oldPlayFunc = stationElm.getAttribute("onclick");
  newPlayFunc = oldPlayFunc.replace("play(", "stop(");
  musicElm.setAttribute("src", `${url}?t=${Date.now()}`);
  musicElm.load();
  musicElm.play().catch(console.error);
  playButton.innerText = "stop";
  playButton.setAttribute("class", "active");
  stationElm.setAttribute("onclick", newPlayFunc);
}

function stop(url, playButtonId) { 
  endStream();
  console.log(playButtonId);
  let playButton = document.querySelector(`button[id="${playButtonId}"]`);
  playButton.innerText = "play";
  playButton.removeAttribute("class");
  oldPlayFunc = playButton.parentElement.getAttribute("onclick");
  newPlayFunc = oldPlayFunc.replace("stop(", "play(");
  playButton.parentElement.setAttribute("onclick", newPlayFunc);
}

function endStream() {
  musicElm.pause();
  musicElm.currentTime = 0;
  musicElm.removeAttribute("src");
  musicElm.load();
}

musicElm.onended = () => {
  musicElm.load();
  musicElm.play();
};

function menu(show) {
  //check if menu is currently active
  let menu = document.getElementById("menu");
  let menuItems = Array.from(menu.querySelectorAll("#item"));
  menuItems.push(menu.querySelector(".close"));
  let menuIcon = [
    menu.querySelector(".top"),
    menu.querySelector(".middle"),
    menu.querySelector(".bottom")
  ]
  if (show) {
    menu.setAttribute("class", "menu");
    menu.removeAttribute("onclick");
    for (let i = 0; i < menuItems.length; i++) {
      let menuItemContent = menuItems[i].innerText;
      menuItems[i].removeAttribute("style");
      if (menuItems[i].getAttribute("class") != "close") {
        menuItems[i].setAttribute("onclick", `link("${menuItemContent}")`);
      } else {
        menuItems[i].setAttribute("onclick", "menu(false)");
      }
    }
      for (let i = 0; i < menuIcon.length; i++) {
        menuIcon[i].setAttribute("style", "display: none;");
      }
    } else {
      menu.removeAttribute("class");
      for (let i = 0; i < menuItems.length; i++) {
        menuItems[i].removeAttribute("onclick");
        menuItems[i].setAttribute("style", "display: none;");
      }
      for (let i = 0; i < menuIcon.length; i++) {
        menuIcon[i].removeAttribute("style");
      }
      menu.addEventListener('mouseup', function() {
        menu.setAttribute("onclick", "menu(true);");
    }, { once: true });
  }
}

function link(what) {
  switch (what) {
    case "github": 
      console.log("    redirecting to GitHub page");
      window.location.assign(githubURL);
      break;
    case "settings":
      settings(true);
      break;
    case "clients":
      console.log(`    todo:  ${what}"`);
      break;
    default:
      console.Error(`    attempted to execute fn link() with undefined value:  ${what}`);
  }
  menu(false);
}

function settings(open) {
  if (open) {
    wall.setAttribute("style", "display: none");
    (async () => {
      const settings = await fetchJSONasArray("settings.json");
      const config = settings.config;

      settingsWall.removeAttribute("style");
      console.log("success fetching settings json as array");
      for (const [key, value] of Object.entries(config)) {
        let settingCont = document.createElement("div");
        settingCont.setAttribute("class", "settingsItem");
        
        let settingName = document.createElement("p");
        settingName.innerText = key;
        settingName.setAttribute("class", "settingName");
        settingCont.appendChild(settingName);
        
        let settingInput = document.createElement("input");
        switch (key) {
          case "log level":
          case "theme":
            settingInput.setAttribute("id", `${key}value`);
            settingInput.setAttribute("name", key);
            settingInput.setAttribute("type", "range");
            settingInput.setAttribute("min", "0");
            settingInput.setAttribute("step", "1");
            settingInput.setAttribute("value", "0");
            if (key == "log level") {
              settingInput.setAttribute("max", "4");
              switch (value) {
              case "fatal":
                settingInput.setAttribute("value", "0");
                break;
              case "error", "err":
                settingInput.setAttribute("value", "1");
                break;
              case "warn":
                settingInput.setAttribute("value", "2");
                break;
              case "info":
                settingInput.setAttribute("value", "3"); 
                break;
              case "debug":
                settingInput.setAttribute("value", "4"); 
                break;
              }
            } else{
              settingInput.setAttribute("max", "2");
            }

            let settingIndicator = document.createElement("p");
            settingIndicator.setAttribute("class", "settingIndicator");
            settingIndicator.setAttribute("id", `${key}Indicator`);
            settingIndicator.innerText = value;
                      
            settingInput.addEventListener("change", function(event) {
              updateIndicator(`${key}`, settingIndicator, event.target.value);
            });

            settingCont.appendChild(settingInput);
            settingCont.appendChild(settingIndicator);
            break;
          case "enable external radios":
            let toggleButton = document.createElement("div");
            toggleButton.setAttribute("class", "toggleButton");
            toggleButton.setAttribute("onclick", `toggle(this)`);
            toggleButton.setAttribute("enabled", "false");
            toggleButton.innerText = value; 
            settingCont.appendChild(toggleButton);
            break;
          default:
            settingInput.setAttribute("id", `${key}value`);
            settingInput.setAttribute("name", key);
            settingInput.setAttribute("type", "input");
            settingInput.setAttribute("value", value);
            settingCont.appendChild(settingInput);
        }

        settingsElm.appendChild(settingCont);
      }/*

      
   

      if (settings[i][3] != null) { 
        settingCont.appendChild(settingDesc);
      }
      if (settings[i][4] != null) {
        settingCont.appendChild(settingIndicator);
      }*/
    })();
  } else {
    settingsElm.innerHTML = "";
    settingsWall.setAttribute("style", "display: none;");
    wall.removeAttribute("style");
  } 
}

function updateIndicator(which, indicatorDummy, valueSTR) {
  let indicator = document.getElementById(`${which}Indicator`);
  let value = parseInt(valueSTR);
  switch (which) {
    case "log level":
      switch (value) {
        case 0 | "fatal":
          indicator.innerText = "fatal";
          break;
        case 1 | "error" |"err":
          indicator.innerText = "error";
          break;
        case 2 | "warn":
          indicator.innerText = "warn";
          break;
        case 3 | "info":
          indicator.innerText = "info";
          break;
        case 4 | "debug":
          indicator.innerText = "debug";
          break;
      }
      break;
    case "theme":
      switch (value) {
        case 0:
          indicator.innerText = "dark";
          break;
        case 1:
          indicator.innerText = "custom";
          break;
        case 2:
          indicator.innerText = "light";
          break;
        default:
          indicator.innerText = "ERROR";
      }
      break;
    case "custom css url":
      indicator.innerText = "todo: custom by url";
      break;
    default:
        console.error(`err! undefined setting change:  ${which}`);
  }
}

function toggle(elm) {
  if (elm.getAttribute("enabled") == "true") {
    elm.setAttribute("enabled", "false");
    elm.innerText = "false"
  } else {
    elm.setAttribute("enabled", "true");
    elm.innerText = "true"
  }
}
