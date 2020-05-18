package telegram

import (
	"log"
	"strings"
	"time"

	"github.com/mattburchett/go_telegram/pkg/core/config"
	"github.com/yanzay/tbot/v2"
)

// Bot contains all the necessary bot and callback information.
type Bot struct {
	Client            *tbot.Client
	Config            config.Config
	Bot               *tbot.Server
	CallbackChatID    string
	CallbackMessageID int
}

// Stat middleware.
func stat(h tbot.UpdateHandler) tbot.UpdateHandler {
	return func(u *tbot.Update) {
		start := time.Now()
		h(u)
		log.Printf("Handle time: %v", time.Now().Sub(start))
	}
}

// New creates an active telegram bot and loads the handlers.
func (tb *Bot) New(token string) {
	tb.Bot = tbot.New(token)
	tb.Bot.Use(stat)
	tb.Client = tb.Bot.Client()
	tb.Handler()
	tb.Bot.Start()

}

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
	tb.Client.SendMessage(m.Chat.ID, "USAGE:\n\n/movie <Movie Name> or /m <Movie Name>\n/show <TV Show Name> or /s <TV Show Name>\n\nEXAMPLES:\n\n/s The Walking Dead\n/m Avatar")
}