package telegram

import (
	"strings"

	"github.com/yanzay/tbot/v2"
)

// Handler creates the active Telegram handlers.
func (tb *Bot) Handler() {

	// Bot Healthcheck
	tb.Bot.HandleMessage("/ping", func(m *tbot.Message) {
		tb.Client.SendMessage(m.Chat.ID, "pong")
	})

	// telegram/sonar.go
	tb.Bot.HandleMessage("/s", tb.sonarrSearch)
	tb.Bot.HandleMessage("/admin sonarrStatus", tb.sonarrStatus)

	// telegram/testhandler.go
	tb.Bot.HandleMessage("/test", tb.testHandler)

	// telegram/admin.go
	tb.Bot.HandleMessage("/admin myID", tb.myID)
	tb.Bot.HandleMessage("/admin chatID", tb.chatID)

	// Help
	tb.Bot.HandleMessage("/help$", tb.helpHandler)
	tb.Bot.HandleMessage("/h$", tb.helpHandler)

	// Callback Handler
	tb.Bot.HandleCallback(tb.callbackHandler)
}

// callbackHandler handles callbacks.
func (tb *Bot) callbackHandler(cq *tbot.CallbackQuery) {
	go func() {
		tb.Client.AnswerCallbackQuery(cq.ID, tbot.OptText("Request received."))
		tb.Client.DeleteMessage(tb.CallbackChatID, tb.CallbackMessageID)
	}()

	if strings.Contains(cq.Data, "tv_") {
		tb.sonarrAdd(cq)
		return
	}

	tb.Client.SendMessage(tb.CallbackChatID, cq.Data)

}

func (tb *Bot) helpHandler(m *tbot.Message) {
	if !tb.whitelistHandler(m) {
		return
	}

	tb.Client.SendMessage(m.Chat.ID, "USAGE:\n\n/movie <Movie Name> or /m <Movie Name>\n/show <TV Show Name> or /s <TV Show Name>\n\nEXAMPLES:\n\n/s The Walking Dead\n/m Avatar")
}

func (tb *Bot) whitelistHandler(m *tbot.Message) bool {
	for _, id := range tb.Config.Telegram.AuthorizedChats {
		if id == m.Chat.ID {
			return true
		}
	}

	tb.Client.SendMessage(m.Chat.ID, "This bot is not authorized for use in this chat.")
	return false

}
