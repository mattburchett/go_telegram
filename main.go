package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/yanzay/tbot"
)

// Config - Specify what to look for in Config file
var Config struct {
	BotToken          string
	SonarrAPIURL      string
	SonarrAPIKey      string
	PlexPyAPIURL      string
	PlexPyAPIKey      string
	RadarrAPIURL      string
	RadarrAPIKey      string
	CouchPotatoAPIURL string
	CouchPotatoAPIKey string
	PlexAPIURL        string
	PlexAPIKey        string
}

func sonarrStatus(message *tbot.Message) {
	response, err := http.Get(Config.SonarrAPIURL + "system/status?apikey=" + Config.SonarrAPIKey)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	message.Replyf("%s", responseData)

}

func sonarrVersion(message *tbot.Message) {
	response, err := http.Get(Config.SonarrAPIURL + "system/status?apikey=" + Config.SonarrAPIKey)
	if err != nil {
		log.Fatal(err)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	type Version struct {
		Version string `json:"version"`
	}

	version := Version{}
	jsonErr := json.Unmarshal(responseData, &version)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	message.Replyf("%s", version.Version)

}

// func activeSteamers(message *tbot.Message) {
// 	response, err := http.Get(Config.PlexAPIURL + "api/v2?apikey=" + Config.PlexAPIKey + "&cmd=")
// }

func main() {
	c := flag.String("c", "./config.json", "Specify the configuration file.")
	flag.Parse()
	file, err := os.Open(*c)
	if err != nil {
		log.Fatal("can't open config file: ", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Config)
	if err != nil {
		log.Fatal("can't decode config JSON: ", err)
	}

	bot, err := tbot.NewServer(Config.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	whitelist := []string{"WARBIRD199"}
	bot.AddMiddleware(tbot.NewAuth(whitelist))

	bot.Handle("/ping", "pong!")

	bot.HandleFunc("/sonarr_status", sonarrStatus)

	bot.HandleFunc("/sonarr_version", sonarrVersion)

	// Start Listening
	err = bot.ListenAndServe()
	log.Fatal(err)

}
