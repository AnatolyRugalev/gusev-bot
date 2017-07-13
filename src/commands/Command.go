package commands

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Command interface {
	Run(update *telegram.Update) (error)
}
