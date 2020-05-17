package telegram

import (
	"fmt"

	"github.com/mattburchett/go_telegram/pkg/service/sonarr"

	"github.com/yanzay/tbot/v2"
)

func (tb *Bot) sonarrStatus(m *tbot.Message) {
	if tb.AdminCheck(m) {
		request, err := sonarr.SonarrStatus(m, tb.Config)

		if err != nil {
			tb.Client.SendMessage(m.Chat.ID, fmt.Sprintf("%v: \n %v", request, err))
		} else {
			tb.Client.SendMessage(m.Chat.ID, "Sonarr Status:")
			tb.Client.SendMessage(m.Chat.ID, request)
		}
	}

}
