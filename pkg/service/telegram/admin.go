package telegram

import (
	"fmt"
	"strconv"

	"github.com/yanzay/tbot/v2"
)

func (tb *Bot) myID(m *tbot.Message) {
	if tb.AdminCheck(m) {
		tb.Client.SendMessage(m.Chat.ID, strconv.Itoa(m.From.ID))
		fmt.Println(m.From.ID)
	}
}

func (tb *Bot) chatID(m *tbot.Message) {
	tb.Client.SendMessage(m.Chat.ID, m.Chat.ID)
}

// AdminCheck checks for valid bot admins.
func (tb *Bot) AdminCheck(m *tbot.Message) bool {
	for _, admin := range tb.Config.Telegram.Admins {
		if m.From.ID == admin {
			return true
		}

	}
	tb.Client.SendMessage(m.Chat.ID, "You are not an authorized admin.")
	return false
}
