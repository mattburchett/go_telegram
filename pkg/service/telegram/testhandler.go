package telegram

import (
	"strings"

	"github.com/yanzay/tbot/v2"
)

func (tb *Bot) testHandler(m *tbot.Message) {
	buttons := make([]string, 0)
	buttons = append(buttons, "ping", "is", "stupid")

	inline2 := make([][]tbot.InlineKeyboardButton, 0)

	for _, i := range buttons {
		inline2 = append(inline2, []tbot.InlineKeyboardButton{{
			Text:         i,
			CallbackData: i,
		}})
	}

	msg, _ := tb.Client.SendMessage(m.Chat.ID, "Inline test. "+strings.TrimPrefix(m.Text, "/test "), tbot.OptInlineKeyboardMarkup(&tbot.InlineKeyboardMarkup{InlineKeyboard: inline2}))
	tb.CallbackMessageID = msg.MessageID
	tb.CallbackChatID = m.Chat.ID
}
