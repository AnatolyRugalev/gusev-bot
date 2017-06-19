package main

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"io/ioutil"
	"github.com/AnatolyRugalev/gusev-bot/src/router"
)

func loadToken() string {
	token, err := ioutil.ReadFile("bot.token")
	if err != nil {
		panic(err)
	}
	return string(token)
}

func main() {
	bot, err := telegram.NewBotAPI(loadToken())
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := telegram.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	r := router.Router{
		Bot: bot,
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		cmd := r.Route(&update)
		if cmd == nil {
			continue
		}
		cmd.Run(&update)
	}
}
