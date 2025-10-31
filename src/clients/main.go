package main

import (
	"os"
	"slices"
	"net/http"
	"encoding/json"
	"io/ioutil"

	"github.com/BurntSushi/toml"
  "github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/list"
)

var (
	confDir = "/.config/Supraboy981322/radio/tui"
	confPath = "/config.toml"
	server = "https://radio.my-lan.dev/"
	args = os.Args[1:]
	url string
	defaultConfig = []byte(`[server]
radio = "https://example.com/"

[style]
libraryTitleColor = "#05b4ff"
libraryItemNameColor = "#4287f5"
libraryItemDescColor = "#2d579c"
`)
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

func main() {
	var taken []string
	for i := 0; i < len(args); i++ {
		if !slices.Contains(taken, args[i]) {
			switch (args[i]) {
			case "-u":
				wrl("url passed")
				url = args[i+1]
				taken = append(taken, args[i+1])
			case "-h":
				help()
			default:
				help()
				fserr("\n\033[31minvalidArg: \033[0m" + args[i])
			}
		}
	}

	wrl("checking install")
	checkInstall()

	wrl("reading config")
	var conf Config
	_, err := toml.DecodeFile(confPath, &conf)
	hanFrr(err) //fn which handles `err`s as fatal

	server = conf.Server.Addr
	if server == "https://example.com/" {
		wserr("\033[31m..you don't appear to have configured the address for your server\033[0m")
		fserr("....see \033[32m-h\033[0m")
	}

	wrl("setting up ui")
	ListTitleColor = lipgloss.Color(conf.Style.LibraryTitleColor)
	ListItemNameColor = lipgloss.Color(conf.Style.LibraryItemNameColor)
	ListItemDescColor = lipgloss.Color(conf.Style.LibraryItemDescColor)

	
	if url == "" {
		wrl("loading ui")
		items := refreshStations()
		startUI(items)
		wrl("done")
	} else {
		fserr("TODO: skip menu and play from url")
	}
}

func help() {
	homeDir, _ := os.UserHomeDir()
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
	os.Exit(0)
}

func refreshStations() []list.Item {
	wrl("fetching stations")
	wrl("..server:  " + server)

	resp, err := http.Get(server + "library.json")
	if err != nil {
		fserr("err fetching library from server:  " + err.Error())
	}
	defer resp.Body.Close()
	
	wrl("reading stations")
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fserr("err reading library:  " + err.Error())
	}
	
	wrl("unmarshalling stations")
	var library [][]string
	err = json.Unmarshal(body, &library)
	if err != nil { 
		fserr("err unmarshalling library json:  " + err.Error())
	}

	wrl("generating list")
	var items []list.Item
	for _, itm := range library {
		items = append(items, item{title: string(itm[0]), desc: string(itm[1])})
	}

	return items
}

func checkInstall() {
	homeDir, err := os.UserHomeDir()
	hanFrr(err)
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
