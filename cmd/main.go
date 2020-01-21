package main

import (
	"fmt"
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

	test := make([][]tb.InlineButton, 0)
	test = append(test, []tb.InlineButton{{Unique: "1", Text: "test"}})

	fmt.Println(test)
	fmt.Println(len(test))
	fmt.Println(len(test[0]))

	b.Handle(test[1][1], func(c *tb.Callback) {
		b.Respond(c, &tb.CallbackResponse{Text: c.ID})
	})

	b.Handle("/ping", func(m *tb.Message) {
		b.Send(m.Sender, "pong")
	})

	b.Handle("/test", func(m *tb.Message) {
		b.Send(m.Sender, "Inline test.", &tb.ReplyMarkup{
			InlineKeyboard: test,
		})
	})

	b.Start()

}

func test(b *tb.Bot) {
	inlineBtns := []tb.InlineButton{tb.InlineButton{Unique: "1", Text: "Ping"}, tb.InlineButton{Unique: "2", Text: "Is"}, tb.InlineButton{Unique: "3", Text: "Stupid"}}
	inlineKeys := [][]tb.InlineButton{inlineBtns}

	for _, btn := range inlineBtns {
		b.Handle(&btn, func(c *tb.Callback) {
			b.Respond(c, &tb.CallbackResponse{Text: c.Message.Text})
		})
	}

	b.Handle("/test", func(m *tb.Message) {
		b.Send(m.Sender, "Inline test.", &tb.ReplyMarkup{
			InlineKeyboard: inlineKeys,
		})
	})
}

func inlineButton(txt string) [][]tb.InlineButton {
	return [][]tb.InlineButton{
		{
			tb.InlineButton{
				Text: txt,
			},
		},
	}
}
