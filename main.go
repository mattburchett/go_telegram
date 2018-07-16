package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

// Configuration - Specify what to look for in Config file
type Configuration struct {
	Token string
}

// ReadConfig from file
func main() {
	c := flag.String("c", "./config.json", "Specify the configuration file.")
	flag.Parse()
	file, err := os.Open(*c)
	if err != nil {
		log.Fatal("can't open config file: ", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	Config := Configuration{}
	err = decoder.Decode(&Config)
	if err != nil {
		log.Fatal("can't decode config JSON: ", err)
	}

	b, err := tb.NewBot(tb.Settings{
		Token:  Config.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/ping", func(m *tb.Message) {
		b.Send(m.Sender, "pong")
	})

	log.Print("Starting bot...")

	b.Start()
}
