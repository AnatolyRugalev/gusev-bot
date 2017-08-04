package router

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/AnatolyRugalev/gusev-bot/src/commands"
)

func Route(update *telegram.Update){
	if update.Message.Text == "/cinema" {
		command:=commands.CinemaCommand{}
		command.Run(update)
	}
}