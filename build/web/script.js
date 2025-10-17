const musicElm = document.getElementById("music");
const wall = document.getElementById("wall");

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
    for (i = 0; i < library.length; i++) {
        let stationData = library[i];
        let station = document.createElement("div");
        station.setAttribute("class", "station");
        
        let title = document.createElement("span");
        title.setAttribute("class", "title");
        title.innerText = stationData[0];

        let playButton = document.createElement("button");
        playButton.setAttribute("id", `playButton${i}`);
        playButton.setAttribute("onclick", "play('" + stationData[1] + "', this)");
        playButton.innerText = "play";

        station.appendChild(title);
        station.appendChild(playButton);
        wall.appendChild(station);
    }
}

function play(url, playButton) {
    const buttons = document.querySelectorAll("button");
    for (i = 0; i < buttons.length; i++) {
        buttons[i].innerText = "play";
        oldPlayFunc = buttons[i].getAttribute("onclick");
        newPlayFunc = oldPlayFunc.replace("stop(", "play(");
        buttons[i].setAttribute("onclick", newPlayFunc);
        buttons[i].removeAttribute("style");
    }
    endStream();
    oldPlayFunc = playButton.getAttribute("onclick");
    newPlayFunc = oldPlayFunc.replace("play(", "stop(");
    musicElm.setAttribute("src", `${url}?t=${Date.now()}`);
    musicElm.load();
    musicElm.play().catch(console.error);
    playButton.innerText = "stop";
    playButton.setAttribute("class", "active");
    playButton.setAttribute("onclick", newPlayFunc);
}

function stop(url, playButton) { 
    endStream();
    playButton.innerText = "play";
    playButton.removeAttribute("class");
    oldPlayFunc = playButton.getAttribute("onclick");
    newPlayFunc = oldPlayFunc.replace("stop(", "play(");
    playButton.setAttribute("onclick", newPlayFunc);
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
