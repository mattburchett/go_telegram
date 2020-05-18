package telegram

import (
	"strconv"

	"github.com/yanzay/tbot/v2"
)

func (tb *Bot) myID(m *tbot.Message) {
	if tb.adminCheck(m.From.ID, false) {
		tb.Client.SendMessage(m.Chat.ID, strconv.Itoa(m.From.ID))
		return
	}

	tb.Client.SendMessage(m.Chat.ID, "You are not an authorized admin.")
}

func (tb *Bot) chatID(m *tbot.Message) {
	if tb.adminCheck(m.From.ID, false) {
		tb.Client.SendMessage(m.Chat.ID, m.Chat.ID)
		return
	}

	tb.Client.SendMessage(m.Chat.ID, "You are not an authorized admin.")
}

// adminCheck checks for valid bot admins.
func (tb *Bot) adminCheck(id int, callback bool) bool {
	for _, admin := range tb.Config.Telegram.Admins {
		if id == admin {
			return true
		}
	}

	return false
}
