package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/yanzay/tbot"
)

// Config - Set layout
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

func sonarrSearch(message *tbot.Message) {
	message.Vars["text"] = strings.Replace(message.Vars["text"], " ", "+", -1)
	r, err := http.Get(Config.SonarrAPIURL + "series/lookup?apikey=" + Config.SonarrAPIKey + "&term=" + message.Vars["text"])

	if err != nil {
		log.Fatalf("There was an error communicating with Sonarr: %v", err)
	}

	rd, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("There was an error parsing the JSON response: %v", err)
	}

	type seriesLookup struct {
		Title       string `json:"title"`
		SortTitle   string `json:"sortTitle"`
		SeasonCount int    `json:"seasonCount"`
		Status      string `json:"status"`
		Year        int    `json:"year"`
		TvdbID      int    `json:"tvdbId"`
		TvRageID    int    `json:"tvRageId"`
		TvMazeID    int    `json:"tvMazeId"`
		CleanTitle  string `json:"cleanTitle"`
		ImdbID      string `json:"imdbId"`
	}

	var sl []seriesLookup
	jsl := json.Unmarshal(rd, &sl)
	if jsl != nil {
		log.Fatalf("There was an error parsing the JSON response (unmarshal): %v", jsl)
	}

	buttons := make([][]string, 0)
	for k := range sl {
		if len(rd) == 2 {
			message.Reply("You must specify a show. Type /usage for usage.")
		} else {
			// message.Replyf("%v (%v) - %v Seasons", sl[k].Title, sl[k].Year, sl[k].SeasonCount)
			results := make([]string, 0)
			output := fmt.Sprintf("%v (%v) - %v Seasons", sl[k].Title, sl[k].Year, sl[k].SeasonCount)
			results = append(results, output)
			buttons = append(buttons, results)

		}
	}
	message.ReplyKeyboard("Please choose a show.", buttons, tbot.WithDataInlineButtons)
}

func usageHelp(message *tbot.Message) {
	message.Reply("USAGE:\n\n/movie <Movie Name> or /m <Movie Name>\n/show <TV Show Name> or /s <TV Show Name>\n\nEXAMPLES:\n\n/s The Walking Dead\n/m Avatar")
}

func badSyntax(message *tbot.Message) {
	message.Reply("You have to specify a name. Type /help for help.")
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
		Response struct {
			Data struct {
				StreamCount string `json:"stream_count"`
			} `json:"data"`
		} `json:"response"`
	}

	a := Activity{}
	ja := json.Unmarshal(rd, &a)
	if ja != nil {
		log.Fatal(ja)
	}

	message.Replyf("Stream Count: %v", a.Response.Data.StreamCount)
}

func movieSearch(message *tbot.Message) {

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

	bot.HandleFunc("/s", badSyntax)
	bot.HandleFunc("/show", badSyntax)
	bot.HandleFunc("/s {text}", sonarrSearch)
	bot.HandleFunc("/show {text}", sonarrSearch)

	bot.HandleFunc("/m", badSyntax)
	bot.HandleFunc("/movie", badSyntax)
	bot.HandleFunc("/m {text}", movieSearch)
	bot.HandleFunc("/movie {text}", sonarrSearch)

	bot.HandleFunc("/sonarr_status", sonarrStatus)

	bot.HandleFunc("/sonarr_version", sonarrVersion)

	bot.HandleFunc("/activity", activeSteamers)

	bot.HandleFunc("/help", usageHelp)
	bot.HandleFunc("/usage", usageHelp)

	// Start Listening
	err = bot.ListenAndServe()
	log.Fatal(err)

}
