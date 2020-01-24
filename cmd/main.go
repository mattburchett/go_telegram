package main

import (
	"log"
	"strings"

	"github.com/mattburchett/go_telegram/pkg/config"
	"github.com/yanzay/tbot/v2"
)

type application struct {
	client            *tbot.Client
	callbackChatID    string
	callbackMessageID int
}

func main() {
	cfg, err := config.GetConfig("config.json")
	if err != nil {
		log.Fatal("Failed to read JSON.")
	}

	app := &application{}

	bot := tbot.New(cfg.TelegramToken)
	app.client = bot.Client()
	c := bot.Client()
	bot.HandleMessage("/ping", func(m *tbot.Message) {
		c.SendMessage(m.Chat.ID, "pong")
	})

	bot.HandleMessage("/test", app.testHandler)
	bot.HandleCallback(app.callbackHandler)

	err = bot.Start()
	if err != nil {
		log.Fatal(err)
	}

}

func (a *application) testHandler(m *tbot.Message) {
	buttons := make([]string, 0)
	buttons = append(buttons, "ping", "is", "stupid")

	inline2 := make([][]tbot.InlineKeyboardButton, 0)

	for _, i := range buttons {
		inline2 = append(inline2, []tbot.InlineKeyboardButton{{
			Text:         i,
			CallbackData: i,
		}})
	}

	msg, _ := a.client.SendMessage(m.Chat.ID, "Inline test. "+strings.TrimPrefix(m.Text, "/test "), tbot.OptInlineKeyboardMarkup(&tbot.InlineKeyboardMarkup{InlineKeyboard: inline2}))
	a.callbackMessageID = msg.MessageID
	a.callbackChatID = m.Chat.ID

}

func (a *application) callbackHandler(cq *tbot.CallbackQuery) {
	a.client.EditMessageText(a.callbackChatID, a.callbackMessageID, "Callback received.")
	a.client.SendMessage(a.callbackChatID, cq.Data)
}
