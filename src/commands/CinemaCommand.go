package commands

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"net/http"
	"regexp"
	"io/ioutil"
	"errors"
	"strconv"
	"gopkg.in/mgo.v2"
	"time"
)

type CinemaCommand struct {
	Bot   *telegram.BotAPI
	Mongo *mgo.Database
}

func (c CinemaCommand) Run(update *telegram.Update) error {
	callback := update.CallbackQuery
	var oldMessage *telegram.Message = &telegram.Message{}
	imageIndex := 0
	if callback != nil {
		oldMessage = callback.Message
		imageIndex, _ = strconv.Atoi(callback.Data)
	}
	urls, err := c.getPictureUrls()
	if err != nil {
		return err
	}
	resp, err := http.Get(urls[imageIndex])
	if err != nil {
		return err
	}
	if resp.Header.Get("Content-Type") != "image/jpeg" {
		msg := telegram.NewMessage(update.Message.Chat.ID, "Ошибка загрузки расписания")
		c.Bot.Send(msg)
	} else {
		var msg telegram.PhotoConfig
		var chatId int64
		file := telegram.FileReader{
			Name:   "cinema.jpg",
			Reader: resp.Body,
			Size:   resp.ContentLength,
		}
		if oldMessage.MessageID > 0 {
			chatId = oldMessage.Chat.ID
			c.Bot.DeleteMessage(telegram.DeleteMessageConfig{
				ChatID: chatId, MessageID: oldMessage.MessageID,
			})
		} else {
			chatId = update.Message.Chat.ID
		}
		msg = telegram.NewPhotoUpload(chatId, file)
		if len(urls) > 1 {
			msg.ReplyMarkup = c.makeMarkup(imageIndex, urls)
		}
		c.Bot.Send(msg)
	}

	return nil
}

func (c CinemaCommand) makeMarkup(index int, urls []string) telegram.InlineKeyboardMarkup {
	var buttons []telegram.InlineKeyboardButton
	if index != 0 {
		buttons = append(buttons, telegram.NewInlineKeyboardButtonData("<<", strconv.Itoa(index-1)))
	}
	if index != len(urls)-1 {
		buttons = append(buttons, telegram.NewInlineKeyboardButtonData(">>", strconv.Itoa(index+1)))
	}
	return telegram.NewInlineKeyboardMarkup(buttons)
}

const baseUrl string = "http://lumenfilm.com"

var urlCache = UrlCache{}

type UrlCache struct {
	Cached  bool
	Expires time.Time
	Urls    []string
}

func (c CinemaCommand) getPictureUrlsFromCache() ([]string, error) {
	if !urlCache.Cached {
		return []string{}, errors.New("no cached data")
	}
	if urlCache.Expires.Unix() < time.Now().Unix() {
		return []string{}, errors.New("data was expired")
	}
	return urlCache.Urls, nil
}

func (c CinemaCommand) putPictureUrlsToCache(urls []string) {
	urlCache = UrlCache{
		Cached: true,
		Expires: time.Now().Add(time.Duration(3600)),
		Urls: urls,
	}
}

func (c CinemaCommand) getPictureUrls() ([]string, error) {
	cached, err := c.getPictureUrlsFromCache()
	if err == nil {
		return cached, nil
	}
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
			result = append(result, baseUrl+match[1])
		}
	}
	if len(result) > 0 {
		c.putPictureUrlsToCache(result)
		return result, nil
	} else {
		return result, errors.New("Cannot match any image on page")
	}
}
