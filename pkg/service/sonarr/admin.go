package sonarr

import (
	"io/ioutil"
	"net/http"

	"github.com/mattburchett/go_telegram/pkg/core/config"
	"github.com/yanzay/tbot/v2"
)

// Status contains the Sonarr request for system status.
func Status(m *tbot.Message, config config.Config) (string, error) {
	r, err := http.Get(config.Sonarr.URL + "system/status?apikey=" + config.Sonarr.APIKey)
	if err != nil {
		return "Failed to contact Sonarr for data", err
	}

	rd, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "Failed to read Sonarr status data.", err
	}

	return string(rd), err
}
