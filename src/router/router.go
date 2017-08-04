package router

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/AnatolyRugalev/gusev-bot/src/commands"
	"gopkg.in/mgo.v2"
)

type Router struct {
	Bot *telegram.BotAPI
	Mongo *mgo.Database
}

func (r *Router) Route(update *telegram.Update) {
	if update.Message != nil {
		if update.Message.Text == "/cinema" {
			r.cinema(update)
		}
	}
	if update.CallbackQuery != nil {
		r.cinema(update)
	}
}

func (r *Router) cinema(update *telegram.Update) {
	command := commands.CinemaCommand{
		Bot: r.Bot,
		Mongo: r.Mongo,
	}
	command.Run(update)
}
