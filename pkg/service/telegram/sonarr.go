package telegram

import (
	"fmt"
	"strings"

	"github.com/mattburchett/go_telegram/pkg/service/sonarr"

	"github.com/yanzay/tbot/v2"
)

// Sonarr Search
func (tb *Bot) sonarrSearch(m *tbot.Message) {
	if !tb.whitelistHandler(m) {
		return
	}

	text := strings.TrimPrefix(strings.TrimPrefix(m.Text, "/s"), " ")
	if len(text) == 0 {
		tb.Client.SendMessage(m.Chat.ID, "You must specify a show. Type /help for help.")
		return
	}

	request, err := sonarr.Search(m, tb.Config)
	if err != nil {
		tb.Client.SendMessage(m.Chat.ID, err.Error())
		return
	}

	inlineResponse := make([][]tbot.InlineKeyboardButton, 0)
	for _, i := range request {
		inlineResponse = append(inlineResponse, []tbot.InlineKeyboardButton{{
			Text:         i.Button,
			CallbackData: "tv_" + i.Callback,
		}})
	}

	if len(request) == 0 {
		tb.Client.SendMessage(m.Chat.ID, "No results found, try harder.")
		return
	}

	response, _ := tb.Client.SendMessage(m.Chat.ID, "Please select the show you would like to download.", tbot.OptInlineKeyboardMarkup(&tbot.InlineKeyboardMarkup{InlineKeyboard: inlineResponse}))
	tb.CallbackMessageID = response.MessageID
	tb.CallbackChatID = m.Chat.ID

}

// sonarrAdd will perform the add requests to Sonarr.
func (tb *Bot) sonarrAdd(cq *tbot.CallbackQuery) {
	if strings.Contains(cq.Data, "+") {
		if tb.adminCheck(cq.From.ID, true) {
			tb.Client.SendMessage(tb.CallbackChatID, sonarr.Add(cq.Data, tb.Config))
		} else {
			tb.Client.AnswerCallbackQuery(cq.ID, tbot.OptText("This request is over the season limit."))
		}
		return
	}

	tb.Client.SendMessage(tb.CallbackChatID, sonarr.Add(cq.Data, tb.Config))

}

// Admin Functions

// sonarrStatus queries Sonarr for it's system status information.
func (tb *Bot) sonarrStatus(m *tbot.Message) {
	if !tb.whitelistHandler(m) {
		return
	}

	if tb.adminCheck(m.From.ID, false) {
		request, err := sonarr.Status(m, tb.Config)

		if err != nil {
			tb.Client.SendMessage(m.Chat.ID, fmt.Sprintf("%v: \n %v", request, err))
			return
		}

		tb.Client.SendMessage(m.Chat.ID, "Sonarr Status:")
		tb.Client.SendMessage(m.Chat.ID, request)

	}

}
