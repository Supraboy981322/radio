const githubURL = "https://supraboy981322.github.io/subpages/text-break?t=todo%3A%20github%20page";
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

async function fetchSettings(url) {
    const res = await fetch(url);
    if (!res.ok) throw new Error(`failed to fetch ${url}:  ${res.status}`);
    const data = await res.json();
    return data
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
//        playButton.setAttribute("onclick", `play("${stationData[1]}", this)`);
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

async function settings(open) {
    if (open) {
        wall.setAttribute("style", "display: none");
        try {
            const settings = await fetchSettings("settings.json");
            
            if (Array.isArray(settings)) {
                settingsWall.removeAttribute("style");
                console.log("success fetching settings json as array");
                settingsElm.innerHTML = "";
                
                for (let i = 0; i < settings.length; i++) {
                    let setting = settings[i];
                    if (typeof setting !== "object" || setting === null) continue;
                    let settingCont = document.createElement("div");
                    settingCont.className = "settingItem";

                    let settingName = document.createElement("p");
                    settingName.innerText = setting.name ?? `setting_${i}`;
                    settingName.className = "settingName";
                    settingCont.appendChild(settingName);

                    let settingInput = document.createElement("input");
                    settingInput.id = `${setting.name}Value`;
                    settingInput.name = setting.name ?? "";
                    switch (setting.type) {
                        case "int":
                            settingInput.type = "range";
                            settingInput.min = String(setting.min);
                            settingInput.max = String(setting.max);
                            settingInput.step = "1";
                            break;
                        case "string":
                            settingInput.type = "text";
                            break;
                        default:
                            console.error(`err! undefined setting type:  ${setting.type}`);
                    }
                    if (setting.value != null) settingInput.value = String(setting.Value);

                    if (setting.desc != null) {
                        let settingDesc = document.createElement("p");
                        settingDesc.className = "settingDesc";
                        settingDesc.innerText = setting.desc;
                        settingCont.appendChild(settingDesc);
                    }


                    if (setting.indicator != null) {
                        let settingIndicator = document.createElement("p");
                        settingIndicator.setAttribute("class", "settingIndicator");
                        settingIndicator.setAttribute("id", `${setting.name}Indicator`);
                        settingIndicator.innerText = String(setting.indicator);
                        settingInput.addEventListener("change", function (event) {
                            updateIndicator(setting.name, settingIndicator, event.target.value);
                        });
                        settingCont.appendChild(settingIndicator);
                    }
                    settingCont.appendChild(settingInput);
                    settingsElm.appendChild(settingCont);
                }
            } else {
                console.error("err fetching settings json as array")
            }
        } catch (err) {
            console.error("err loading settings:  ", err);
        }
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
            console.error(`err! undefined setting change! this is very weird behavior and should never occur:  ${which}`);
    }
}

function saveSettings() {
    //let themeData = settingsData.item[0];
    //let customCSSdata = settingsData.item[1];
    for (let i = 0; i < settingsData.length(); i++) {
        settingsData.item[i].value = document.getElementById("themeValue").value;
        settingsData.item[1].value = document.getElementById("customCSSValue").value;
    }

    console.log("attempting to save settings");
    
}
