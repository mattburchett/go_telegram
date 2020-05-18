package sonarr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/mattburchett/go_telegram/pkg/core/config"
	"github.com/yanzay/tbot/v2"
)

// Response holds the information needed for Telegram callback.
type response struct {
	Button   string `json:"button"`
	Callback string `json:"callback"`
}

type sonarrSearch []struct {
	Title            string `json:"title"`
	SeasonCount      int    `json:"seasonCount"`
	Year             int    `json:"year"`
	TvdbID           int    `json:"tvdbId"`
	Downloaded       bool   `json:"downloaded"`
	QualityProfileID int    `json:"qualityProfileId"`
	TitleSlug        string `json:"titleSlug"`
	Images           []struct {
		CoverType string `json:"coverType"`
		URL       string `json:"url"`
	} `json:"images"`
	Seasons []struct {
		SeasonNumber int  `json:"seasonNumber"`
		Monitored    bool `json:"monitored"`
	} `json:"seasons"`
	ProfileID int `json:"profileId"`
}

type sonarrAdd struct {
	Title            string `json:"title"`
	TvdbID           int    `json:"tvdbId"`
	QualityProfileID int    `json:"qualityProfileId"`
	TitleSlug        string `json:"titleSlug"`
	Images           []struct {
		CoverType string `json:"coverType"`
		URL       string `json:"url"`
	} `json:"images"`
	Seasons []struct {
		SeasonNumber int  `json:"seasonNumber"`
		Monitored    bool `json:"monitored"`
	} `json:"seasons"`
	ProfileID      int    `json:"profileId"`
	RootFolderPath string `json:"rootFolderPath"`
	AddOptions     struct {
		IgnoreEpisodesWithFiles    bool `json:"ignoreEpisodesWithFiles"`
		IgnoreEpisodesWithoutFiles bool `json:"ignoreEpisodesWithoutFiles"`
		SearchForMIssingEpisodes   bool `json:"searchForMissingEpisodes"`
	} `json:"addOptions"`
}

type RootFolderLookup []struct {
	Path            string `json:"path"`
	FreeSpace       int64  `json:"freeSpace"`
	TotalSpace      int64  `json:"totalSpace"`
	UnmappedFolders []struct {
		Name string `json:"name"`
		Path string `json:"path"`
	} `json:"unmappedFolders"`
	ID int `json:"id"`
}

// Search performs the lookup actions within Sonarr.
func Search(m *tbot.Message, config config.Config) ([]response, error) {

	var remote sonarrSearch
	var local sonarrSearch

	// Perform series lookup
	remoteLookup, err := http.Get(config.Sonarr.URL +
		"series/lookup?apikey=" +
		config.Sonarr.APIKey +
		"&term=" +
		url.QueryEscape(strings.TrimPrefix(strings.TrimPrefix(m.Text, "/s"), " ")))

	if err != nil {
		return nil, err
	}

	// Perform series local database lookup.
	localLookup, err := http.Get(config.Sonarr.URL + "series?apikey=" + config.Sonarr.APIKey)
	if err != nil {
		return nil, err
	}

	// Read remote and local Data
	remoteData, err := ioutil.ReadAll(remoteLookup.Body)
	if err != nil {
		return nil, err
	}

	localData, err := ioutil.ReadAll(localLookup.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal remote and local JSON
	remoteJSON := json.Unmarshal(remoteData, &remote)
	if remoteJSON != nil {
		return nil, err
	}

	localJSON := json.Unmarshal(localData, &local)
	if localJSON != nil {
		return nil, err
	}

	// Check for downloaded items.
	for r := range remote {
		for l := range local {
			if remote[r].TvdbID == local[l].TvdbID {
				remote[r].Downloaded = true
			}
		}
	}

	// Form URLs and return to Telegram
	responseData := []response{}
	for k := range remote {
		responseData = append(responseData,
			response{
				fmt.Sprintf("%v%v%v (%v) - %v Seasons",
					seasonHandler(remote[k].SeasonCount, config.Sonarr.SeasonLimit),
					downloadedHandler(remote[k].Downloaded),
					remote[k].Title,
					remote[k].Year,
					remote[k].SeasonCount),
				fmt.Sprintf("%d%s%s", remote[k].TvdbID,
					strings.TrimSuffix(downloadedHandler(remote[k].Downloaded), " "),
					strings.TrimSuffix(seasonHandler(remote[k].SeasonCount, config.Sonarr.SeasonLimit), " ")),
			})
	}

	return responseData, err

}

// Add will take the callback data and add the show to Sonarr.
func Add(callback string, config config.Config) string {
	// Separate the TVDB ID
	tvdbid := strings.TrimPrefix(strings.TrimSuffix(strings.TrimSuffix(callback, "+"), "*"), "tv_")

	// Look it up, to gather the information needed and the title.
	seriesLookup, err := http.Get(config.Sonarr.URL + "/series/lookup?apikey=" + config.Sonarr.APIKey + "&term=tvdb:" + tvdbid)
	if err != nil {
		return err.Error()
	}

	seriesData, err := ioutil.ReadAll(seriesLookup.Body)
	if err != nil {
		return err.Error()
	}

	series := sonarrSearch{}
	seriesJSON := json.Unmarshal(seriesData, &series)
	if seriesJSON != nil {
		return err.Error()
	}

	// If the "downloaded" asterisk is already in place, go ahead and return.
	if strings.Contains(callback, "*") {
		return fmt.Sprintf("%v has already been requested for download.", series[0].Title)
	}

	// Gather the root folder location.
	rootFolderLookup, err := http.Get(config.Sonarr.URL + "/rootfolder?apikey=" + config.Sonarr.APIKey)
	if err != nil {
		return err.Error()
	}
	rootFolderData, err := ioutil.ReadAll(rootFolderLookup.Body)
	if err != nil {
		return err.Error()
	}

	rootFolder := RootFolderLookup{}
	rootFolderJSON := json.Unmarshal(rootFolderData, &rootFolder)
	if rootFolderJSON != nil {
		return err.Error()
	}

	// Form the JSON needed for adding to Sonarr.
	seriesAdd := sonarrAdd{
		TvdbID: series[0].TvdbID,
		Title:  series[0].Title,
		// QualityProfileID: series[0].QualityProfileID,
		TitleSlug:      series[0].TitleSlug,
		Images:         series[0].Images,
		Seasons:        series[0].Seasons,
		RootFolderPath: rootFolder[0].Path,
		ProfileID:      config.Sonarr.ProfileID,
	}
	seriesAdd.AddOptions.IgnoreEpisodesWithFiles = false
	seriesAdd.AddOptions.IgnoreEpisodesWithoutFiles = false
	seriesAdd.AddOptions.SearchForMIssingEpisodes = true

	// Post it to Sonarr to be added.
	seriesAddJSON, err := json.Marshal(seriesAdd)
	if err != nil {
		return err.Error()
	}
	seriesAddReq, err := http.Post(config.Sonarr.URL+"/series?apikey="+config.Sonarr.APIKey, "application/json", bytes.NewBuffer(seriesAddJSON))
	if err != nil {
		return err.Error()
	}

	if seriesAddReq.StatusCode != 201 {
		return "There was an error processing this request."
	}
	return fmt.Sprintf("%s has been queued for download.", series[0].Title)

}

// downloadHandler returns the proper string that should be shown in the Telegram response.
func downloadedHandler(downloaded bool) string {
	if downloaded {
		return "* "
	}

	return ""
}

// seasonHandler returns the proper string that should be shown in the Telegram response.
// it is also used to prevent non-admins from downloading shows with a high number of seasons.
func seasonHandler(seasonCount int, seasonLimit int) string {
	if seasonLimit == 0 {
		return ""
	}

	if seasonCount >= seasonLimit {
		return "+ "
	}

	return ""
}
