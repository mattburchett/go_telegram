package telegram

import (
	"log"
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
