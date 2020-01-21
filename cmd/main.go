package main

import (
	"log"
	"time"

	"github.com/mattburchett/go_telegram/pkg/config"
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {

	cfg, err := config.GetConfig("config.json")
	if err != nil {
		log.Fatal("Failed to load config.")
	}

	b, err := tb.NewBot(tb.Settings{
		Token:  cfg.TelegramToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/ping", func(m *tb.Message) {
		b.Send(m.Sender, "pong")
	})

	test(b)

	b.Start()

}

func test(b *tb.Bot) {
	var stored *tb.Message
	inlineBtns := []tb.InlineButton{tb.InlineButton{Unique: "1", Text: "Ping"}, tb.InlineButton{Unique: "2", Text: "Is"}, tb.InlineButton{Unique: "3", Text: "Stupid"}}
	inlineKeys := [][]tb.InlineButton{inlineBtns}

	b.Handle("/test", func(m *tb.Message) {
		msg, _ := b.Send(m.Sender, "Inline test.", &tb.ReplyMarkup{
			InlineKeyboard: inlineKeys,
		})
		stored = msg

	})

	for _, btn := range inlineBtns {
		type button struct {
			ID    int
			Stuff tb.InlineButton
		}

		b.Handle(&btn, func(c *tb.Callback) {
			b.Respond(c, &tb.CallbackResponse{Text: "Callback received."})
			func(m *tb.Message) {
				b.EditReplyMarkup(m, &tb.ReplyMarkup{ReplyKeyboardRemove: true})
				b.Edit(m, "Callback received. Processing.")
			}(stored)
			b.Send(c.Sender, btn.Text)
		})
	}
}
