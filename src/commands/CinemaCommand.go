package commands

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"net/http"
	"time"
	"fmt"
)

type CinemaCommand struct {
	Command
	Bot *telegram.BotAPI
	url string
}

func (c CinemaCommand) Run(update *telegram.Update) error {

	url := c.getPictureUrl()
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	msg := telegram.NewPhotoUpload(update.Message.Chat.ID, telegram.FileReader{
		Name: "KINO.JPG",
		Reader: resp.Body,
		Size: resp.ContentLength,
	})
	msg.Caption = "KINOSHKA"

	c.Bot.Send(msg)
	return nil
}

func (c CinemaCommand) getPictureUrl() string {
	t := time.Now()
	return fmt.Sprintf("http://lumenfilm.com/uploads/%02d_%02d.jpg", t.Day(), t.Month())
}
