package commands

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"net/http"
	"regexp"
	"io/ioutil"
	"errors"
)

type CinemaCommand struct {
	Bot    *telegram.BotAPI
	Update *telegram.Update
	url    string
}

func (c CinemaCommand) Run(update *telegram.Update) error {
	urls, err := c.getPictureUrls()
	if err != nil {
		return err
	}
	resp, err := http.Get(urls[0])
	if err != nil {
		return err
	}
	var msg telegram.Chattable
	if resp.Header.Get("Content-Type") != "image/jpeg" {
		msg = telegram.NewMessage(c.Update.Message.Chat.ID, "Ошибка загрузки расписания")
	} else {
		msg = telegram.NewPhotoUpload(update.Message.Chat.ID, telegram.FileReader{
			Name:   "cinema.jpg",
			Reader: resp.Body,
			Size:   resp.ContentLength,
		})
	}

	c.Bot.Send(msg)
	return nil
}

const baseUrl string = "http://lumenfilm.com"

func (c CinemaCommand) getPictureUrls() ([]string, error) {
	resp, err := http.Get(baseUrl + "/gusev/affishe")
	var result []string
	if err != nil {
		return result, err
	}
	content, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	r := regexp.MustCompile("<img\\s.*?src=[\"']([^'\"]+)[\"']\\s?.*?>")
	matches := r.FindAllStringSubmatch(string(content), 7)
	for _, match := range matches {
		if len(match) > 1 {
			result = append(result, baseUrl + match[1])
		}
	}
	if len(result) > 0 {
		return result, nil
	} else {
		return result, errors.New("Cannot match any image on page")
	}
}
