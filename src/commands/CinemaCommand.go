package commands

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"net/http"
	"time"
	"fmt"
	"github.com/earlcherry/gouter"
)

type CinemaCommand struct {
	Bot *telegram.BotAPI
	Update *telegram.Update
	url string
}

func (c CinemaCommand) Run(args gouter.Args) error {

	url := c.getPictureUrl()
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	var msg telegram.Chattable
	if resp.Header.Get("Content-Type") != "image/jpeg" {
		msg = telegram.NewMessage(c.Update.Message.Chat.ID, "Ошибка загрузки расписания")
	} else {
		msg = telegram.NewPhotoUpload(c.Update.Message.Chat.ID, telegram.FileReader{
			Name:   "cinema.jpg",
			Reader: resp.Body,
			Size:   resp.ContentLength,
		})
	}

	c.Bot.Send(msg)
	return nil
}

func (c CinemaCommand) getPictureUrl() string {
	t := time.Now()
	return fmt.Sprintf("http://lumenfilm.com/uploads/%02d.jpg", t.Day())
}
