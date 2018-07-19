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
	BotToken     string
	SonarrAPIURL string
	SonarrAPIKey string
	PlexPyAPIURL string
	PlexPyAPIKey string
	// RadarrAPIURL      string
	// RadarrAPIKey      string
	// CouchPotatoAPIURL string
	// CouchPotatoAPIKey string
	// PlexAPIURL        string
	// PlexAPIKey        string
}

func sonarrStatus(message *tbot.Message) {
	r, err := http.Get(Config.SonarrAPIURL + "system/status?apikey=" + Config.SonarrAPIKey)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	rd, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	message.Replyf("%s", rd)

}

func sonarrVersion(message *tbot.Message) {
	r, err := http.Get(Config.SonarrAPIURL + "system/status?apikey=" + Config.SonarrAPIKey)
	if err != nil {
		log.Fatal(err)
	}

	rd, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	type Version struct {
		Version string `json:"version"`
	}

	v := Version{}
	jv := json.Unmarshal(rd, &v)
	if jv != nil {
		log.Fatal(jv)
	}

	message.Replyf("%s", v.Version)

}

func activeSteamers(message *tbot.Message) {
	r, err := http.Get(Config.PlexPyAPIURL + "?apikey=" + Config.PlexPyAPIKey + "&cmd=get_activity")
	if err != nil {
		log.Fatal(err)
	}

	rd, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	type Activity struct {
		StreamCount             int `json:"stream_count"`
		StreamCountDirectPlay   int `json:"stream_count_direct_play"`
		StreamCountDirectStream int `json:"stream_count_direct_stream"`
		StreamCountTranscode    int `json:"stream_count_transcode"`
	}

	a := Activity{}
	ja := json.Unmarshal(rd, &a)
	if ja != nil {
		log.Fatal(ja)
	}

	message.Replyf("Stream Count: %v \nStream Count (Direct Play): %v \nStream Count (Direct Stream): %v \nStream Count (Transcode): %v", a.StreamCount, a.StreamCountDirectPlay, a.StreamCountDirectStream, a.StreamCountTranscode)
}

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

	bot.HandleFunc("/activity", activeSteamers)

	// Start Listening
	err = bot.ListenAndServe()
	log.Fatal(err)

}
