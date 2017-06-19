package router

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/AnatolyRugalev/gusev-bot/src/commands"
)

type Router struct {
	Bot *telegram.BotAPI
}

func (r *Router) Route (update *telegram.Update) commands.Command {
	if update.Message.Text == "/cinema" {
		return commands.CinemaCommand{
			Bot: r.Bot,
		}
	} else {
		return nil
	}
}
