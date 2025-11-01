package main

import (
	"os"
	"slices"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/BurntSushi/toml"
  "github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/list"
)

var (
	confDir = "/.config/Supraboy981322/radio/tui" //updated to be in user home dir
	confPath = "/config.toml"	//updated later to be in conf dir
	server string	 //radio server address
	args = os.Args[1:]	//cli args
	url string	//for if the user passes the url arg
	curTask = 0  //this is global because it simplifies the startup logging
	defaultConfig = []byte(`[server]
radio = "https://example.com/"

[style]
libraryTitleColor = "#359dde"
libraryItemNameColor = "#fccd12"
libraryItemDescColor = "#ba810f"`)
	//this simplifies the startup code
	//  (used in a fn later)
	startupTasks = []string{
		"starting...",
		"checking install",
		"reading config",
		"settings up ui",
		"fetching stations",
		"reading stations",
		"unmarshalling stations",
		"constructing list",
		"loading ui",
		"done.",
	}
)

type (
	ConfigServer struct {
		Addr  string `toml:"radio"`
	}
	ConfigStyle struct {
		LibraryTitleColor string `toml:"libraryTitleColor"`
		LibraryItemNameColor string `toml:"libraryItemNameColor"`
		LibraryItemDescColor string `toml:"libraryItemDescColor"`
	}
	Config struct {
		Server  ConfigServer
		Style   ConfigStyle
	}
)

//I know this looks like trash, but
//  it's better than writing 
//  `fmt.Println(...)` (or the `log`/`slog` equivalent)
//  for printing each startup tasks
func wrTsk() {               //basically...
	wrl("\033[36m[\033[33m" +  //just `[` with an ansii color code
		strconv.Itoa(curTask) +  //conv task int to str then append to line 
		"\033[36m/\033[32m" +    //just `/` with an ansii color code
		strconv.Itoa(len(startupTasks) - 1) +  //len of list from earlier
		"\033[36m]:\033[0m " +   //just `]` with an ansii color code
		startupTasks[curTask])   //add the task string from list
	curTask++                  //increment task counter
}

func main() {
	//first of many calls for that that startup
	//  task printing fn 
	wrTsk()
	
	//iterate through the args passed
	//  and keep track of which args 
	//  are values and which are keys
	var taken []string
	for i := 0; i < len(args); i++ {
		if !slices.Contains(taken, args[i]) {
			switch (args[i]) {
			case "-u":
				//fn that's basically just `fmt.Println`
				//  with less typing
				wrl("url passed")
				url = args[i+1]
				taken = append(taken, args[i+1])
			case "-h":
				help()
			default:
				help() //if unrecognized arg, print help
				//fn to fatal print err with less typing
				fserr("\n\033[31minvalidArg: \033[0m" + args[i])
			}
		}
	}

	wrTsk() //another call of that fn
	checkInstall() //call fn to verify install & conf

	wrTsk() //i know this might seem repetitive

	//read config file
	var conf Config
	_, err := toml.DecodeFile(confPath, &conf)
	hanFrr(err) //fn that's just `if err != nil { ... }`

	//set the server address from the config 
	server = conf.Server.Addr
	if server == "https://example.com/" {
		//a fn that prints string as err but with
		//  less typing (diff from fatal variant)
		wserr("\033[31m..you don't appear to " + 
						"have configured the address" + 
						"for your server\033[0m")
		//the same fn as before but fatal
		fserr("....see \033[32m-h\033[0m")
	}

	wrTsk() //i shall stop commenting on this fn
	//set the ui colors using conf values 
	ListTitleColor = lipgloss.Color(
						conf.Style.LibraryTitleColor)
	ListItemNameColor = lipgloss.Color(
						conf.Style.LibraryItemNameColor)
	ListItemDescColor = lipgloss.Color(
						conf.Style.LibraryItemDescColor)

	
	if url == "" {
		wrTsk()
		//refresh the stations library
		items := refreshStations()
		//strt bubble tea
		startUI(items)
		wrTsk()
	} else {
		//skipping to playing a song isn't
		//  yet implemented 
		fserr("TODO: skip menu and play from url")
	}
}

func help() {
	//get the uer's home dir
	homeDir, _ := os.UserHomeDir()

	//another "needlessly complex" block
	//  but it's to save typing (in my case)
	//  `wrl(...)` so many times, basically
	//  create a str arr (1 val == 1 line)
	//  then iterate through the arr and
	//  print each val as a ln
	var li = []string{
		"\033[34mradio\033[0m --> \033[32mhelp\033[0m",
		"  \033[32m-h:\033[0m",
		"    this screen",
		"  \033[33m-u:\033[0m",
		"    skip menu and play url",
		"     eg: \033[34mradio\033[0m \033[33m-u\033[0m https://example.com/foo.mp3",
		"\n  how to edit your config:",
		"    your config path: \033[36m" + homeDir + confDir + confPath + "\033[0m",
		"    \033[32m1)\033[0m with the editor of your choice, open config file using the path above",
		"    \033[32m2)\033[0m enter the url for your \033[34mradio\033[0m server (in quotes)",
		"        eg: \033[30;46m radio = \"https://example.com/\" \033[0m",
		"    for customization, see refer to the online docs", 
	}
	for i := 0; i < len(li); i++ {
		wrl(li[i])
	}

	//then exit
	os.Exit(0)
}

func refreshStations() []list.Item {
	wrTsk()

	//get the library from server 
	resp, err := http.Get(server + "library.json")
	hanFrr(err)
	defer resp.Body.Close()
	
	wrTsk()
	//read the library
	body, err := ioutil.ReadAll(resp.Body)
	hanFrr(err)
	
	wrTsk()
	//unmarshal library
	var library [][]string
	err = json.Unmarshal(body, &library)
	hanFrr(err)

	wrTsk()
	//generate the library list for Bubble Tea
	var items []list.Item
	for _, itm := range library {
		items = append(items, item{
			title: string(itm[0]),
			desc: string(itm[1])})
	}

	return items
}

func checkInstall() {
	//get usr home
	homeDir, err := os.UserHomeDir()
	hanFrr(err)
	
	//remember at the start where I said I'd
	//  update confDir and confPath
	confDir = homeDir + confDir
	confPath = confDir + confPath

	//chk conf dir
	_, err = os.Stat(confPath)
	if os.IsNotExist(err) {
		//assume dir doesn't exist
		hanErr(os.MkdirAll(confDir, 0755))
		
		//create conf with defaults
		hanFrr(os.WriteFile(confPath, defaultConfig, 0644))
	}
}
