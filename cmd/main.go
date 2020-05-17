package main

import (
	"log"

	"github.com/mattburchett/go_telegram/pkg/core/config"
	"github.com/mattburchett/go_telegram/pkg/service/telegram"
)

func main() {
	conf, err := config.GetConfig("config.json")
	if err != nil {
		log.Fatal("Failed to read JSON.")
	}

	tgBot := telegram.Bot{}
	tgBot.Config = conf
	tgBot.New(conf.Telegram.Token)
}
