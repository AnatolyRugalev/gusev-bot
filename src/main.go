package main

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"io/ioutil"
	"flag"
	"encoding/json"
	"github.com/AnatolyRugalev/gusev-bot/src/config"
	"github.com/earlcherry/gouter"
	"gopkg.in/mgo.v2"
	"fmt"
	"github.com/AnatolyRugalev/gusev-bot/src/commands"
)

func loadConfig() (*config.AppConfig, error) {
	configFile := flag.String("config", "config.json", "appConfig file path")
	configJson, err := ioutil.ReadFile(*configFile)
	var appConfig config.AppConfig
	err = json.Unmarshal(configJson, &appConfig)
	if err != nil {
		return nil, err
	}
	return &appConfig, nil
}

func makeBot(botConfig config.BotConfig) *telegram.BotAPI {
	bot, err := telegram.NewBotAPI(botConfig.Token)
	if err != nil {
		log.Panic(err)
	}
	return bot
}

func makeMongo(mongoConfig config.MongoConfig) *mgo.Database {
	session, err := mgo.Dial(mongoConfig.Host)

	if err != nil {
		log.Panic(err)
	}
	return session.DB(mongoConfig.Database)
}

type CommandArgs struct {
	gouter.SimpleRouteArgs
	Update *telegram.Update
}

func main() {
	appConfig, err := loadConfig()
	if err != nil {
		log.Panic(err)
	}

	bot := makeBot(appConfig.Bot)

	mongo := makeMongo(appConfig.Mongo)

	collections, err := mongo.CollectionNames()

	fmt.Printf("%v", collections)

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := telegram.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	g := gouter.New()
	g.Add("/cinema", commands.CinemaCommand{})

	for update := range updates {
		if update.Message == nil {
			continue
		}
		args := CommandArgs{
			Update: &update,
		}
		args.SetRoute(update.Message.Text)
		g.Run(args)
	}
}
