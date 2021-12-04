package couchpotato

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/mattburchett/go_telegram/pkg/core/config"
	"github.com/yanzay/tbot/v2"
)

type couchpotatoSearch struct {
	Movies []struct {
		Title     string `json:"original_title"`
		Imdb      string `json:"imdb"`
		Year      int    `json:"year"`
		InLibrary bool   `json:"in_library"`
		InWanted  bool   `json:"in_wanted"`
	} `json:"movies"`
	Success bool `json:"success"`
}

type request struct {
	ImdbID     string `json:"imdbid"`
	Title      string `json:"title"`
	Year       int    `json:"year"`
	Requested  bool   `json:"requested"`
	Downloaded bool   `json:"downloaded"`
}

func (search couchpotatoSearch) Convert() []request {
	requests := []request{}
	for _, result := range search.Movies {
		requests = append(requests, request{
			ImdbID:     result.Imdb,
			Title:      result.Title,
			Year:       result.Year,
			Requested:  result.InWanted,
			Downloaded: result.InLibrary,
		})

	}

	return requests
}

type response struct {
	Button   string `json:"button"`
	Callback string `json:"callback"`
}

// Search performs the lookup actions within CouchPotato
func Search(m *tbot.Message, config config.Config) ([]response, error) {
	searchLookup, err := http.Get(config.CouchPotato.URL + config.CouchPotato.APIKey + "/movie.search?q=" + url.QueryEscape(strings.TrimPrefix(strings.TrimPrefix(m.Text, "/m"), " ")))
	if err != nil {
		return nil, err
	}

	search := couchpotatoSearch{}

	searchData, err := ioutil.ReadAll(searchLookup.Body)
	if err != nil {
		return nil, err
	}

	requestJSON := json.Unmarshal(searchData, &search)
	if requestJSON != nil {
		return nil, err
	}

	requests := search.Convert()

	responseData := []response{}
	for _, r := range requests {
		responseData = append(responseData,
			response{
				fmt.Sprintf("%v (%v)", r.Title, r.Year),
				fmt.Sprintf("%v", r.ImdbID),
			})

	}

	return responseData, err
}
