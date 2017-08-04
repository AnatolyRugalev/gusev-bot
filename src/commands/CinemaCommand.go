package commands

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"net/http"
	"regexp"
	"io/ioutil"
)

type CinemaCommand struct {
	Bot    *telegram.BotAPI
	Update *telegram.Update
	url    string
}

func (c CinemaCommand) Run(update *telegram.Update) error {
	url := c.getPictureUrl()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	var msg telegram.Chattable
	//if resp.Header.Get("Content-Type") != "image/jpeg" {
	//	msg = telegram.NewMessage(c.Update.Message.Chat.ID, "Ошибка загрузки расписания")
	//} else {
		msg = telegram.NewPhotoUpload(update.Message.Chat.ID, telegram.FileReader{
			Name:   "cinema.jpg",
			Reader: resp.Body,
			Size:   resp.ContentLength,
		})
	//}

	c.Bot.Send(msg)
	return nil
}

const baseUrl string = "http://lumenfilm.com"

func (c CinemaCommand) getPictureUrl() string {
	resp, err := http.Get(baseUrl + "/gusev/affishe")
	if err != nil {
		panic(err)
	}
	content, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	r := regexp.MustCompile("<img\\s.*?src=[\"']([^'\"]+)[\"']\\s?.*?>")
	match := r.FindStringSubmatch(string(content))
	if len(match) > 1 {
		return baseUrl+match[1]
	}
	return ""
}
