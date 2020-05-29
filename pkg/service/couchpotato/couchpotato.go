package couchpotato

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/mattburchett/go_telegram/pkg/core/config"
	"github.com/yanzay/tbot/v2"
)

type response struct {
	Movies []struct {
		Titles    []string `json:"titles"`
		Imdb      string   `json:"imdb"`
		Year      int      `json:"year"`
		InLibrary struct {
			Status string `json:"status"`
		} `json:"in_library"`
		InWanted bool `json:"in_wanted"`
	} `json:"movies"`
	Success bool `json:"success"`
}

type request struct {
	ImdbID     string `json:"imdbid"`
	Title      string `json:"title"`
	Year       int    `json:"year"`
	Requested  bool   `json:"requested"`
	Downloaded struct {
		Status string `json:"status"`
	} `json:"downloaded"`
}

// Search performs the lookup actions within CouchPotato
func Search(m *tbot.Message, config config.Config) ([]response, error) {
	requestLookup, err := http.Get(config.CouchPotato.URL + "/api/" + config.CouchPotato.APIKey + "/movie.search?q=" + url.QueryEscape(strings.TrimPreFix(strings.TrimPrefix(m.Text, "/s"), " ")))
	if err != nil {
		return nil, err
	}

	request := []response{}

	requestData, err := ioutil.ReadAll(requestLookup.Body)
	if err != nil {
		return nil, err
	}

	requestJSON := json.Unmarshal(requestData, &request)
	if requestJSON != nil {
		return nil, err
	}

	for r := range request {

	}
}
